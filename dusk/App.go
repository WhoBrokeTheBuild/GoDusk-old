package dusk

import (
	"fmt"
	"strings"
	"time"

	"github.com/go-gl/gl/v4.1-core/gl"
	"github.com/go-gl/mathgl/mgl32"
	"github.com/veandco/go-sdl2/sdl"
)

type App struct {
	WindowWidth  int
	WindowHeight int
	WindowTitle  string
	TargetFps    float32

	EvtUpdate []func(ctx *UpdateContext)
	EvtRender []func(ctx *RenderContext)

	updateCtx  UpdateContext
	renderCtx  RenderContext
	sdlWindow  *sdl.Window
	sdlContext sdl.GLContext
	running    bool
}

func NewApp() (App, error) {
	var err error

	app := App{
		WindowTitle:  "Dusk",
		WindowWidth:  640,
		WindowHeight: 480,
		TargetFps:    60,
		updateCtx: UpdateContext{
			Frame: 0,
		},
		renderCtx: RenderContext{},
	}

	sdl.Init(sdl.INIT_EVERYTHING)

	sdl.GL_SetAttribute(sdl.GL_CONTEXT_FLAGS, sdl.GL_CONTEXT_FORWARD_COMPATIBLE_FLAG)
	sdl.GL_SetAttribute(sdl.GL_CONTEXT_MAJOR_VERSION, 4)
	sdl.GL_SetAttribute(sdl.GL_CONTEXT_MINOR_VERSION, 1)
	sdl.GL_SetAttribute(sdl.GL_CONTEXT_PROFILE_MASK, sdl.GL_CONTEXT_PROFILE_CORE)
	sdl.GL_SetAttribute(sdl.GL_DOUBLEBUFFER, 1)
	sdl.GL_SetAttribute(sdl.GL_DEPTH_SIZE, 24)
	sdl.GL_SetAttribute(sdl.GL_MULTISAMPLEBUFFERS, 1)

	app.sdlWindow, err = sdl.CreateWindow(
		app.WindowTitle,
		sdl.WINDOWPOS_CENTERED, sdl.WINDOWPOS_CENTERED,
		app.WindowWidth, app.WindowHeight,
		sdl.WINDOW_OPENGL|sdl.WINDOW_SHOWN,
	)

	if err != nil {
		LogError("Failed to create Window, %v", err)
		return app, err
	}

	app.sdlContext, err = sdl.GL_CreateContext(app.sdlWindow)
	if err != nil {
		LogError("Failed to create GL Context, %v", err)
		return app, err
	}

	err = gl.Init()
	if err != nil {
		LogError("Failed to initialize GLOW, %v", err)
		return app, err
	}

	LogInfo("OpenGL Version %v", gl.GoStr(gl.GetString(gl.VERSION)))
	LogInfo("GLSL Version %v", gl.GoStr(gl.GetString(gl.SHADING_LANGUAGE_VERSION)))
	LogInfo("OpenGL Vendor %v", gl.GoStr(gl.GetString(gl.VENDOR)))
	LogInfo("OpenGL Renderer %v", gl.GoStr(gl.GetString(gl.RENDERER)))

	var samples int32
	gl.GetIntegerv(gl.SAMPLES, &samples)
	LogInfo("Anti-Aliasing %vx", samples)

	var formats int32
	gl.GetIntegerv(gl.NUM_PROGRAM_BINARY_FORMATS, &formats)
	LogInfo("Binary Shader Formats %v", formats)

	sdl.GL_SetSwapInterval(1)

	gl.Enable(gl.MULTISAMPLE)

	gl.Enable(gl.DEPTH_TEST)
	gl.DepthFunc(gl.LEQUAL)

	gl.Enable(gl.BLEND)
	gl.BlendFunc(gl.SRC_ALPHA, gl.ONE_MINUS_SRC_ALPHA)

	gl.ClearColor(0.3, 0.3, 0.3, 1.0)

	// --
	// --
	// --

	progId := gl.CreateProgram()

	vertId := gl.CreateShader(gl.VERTEX_SHADER)

	srcs, free := gl.Strs(vertShader)
	gl.ShaderSource(vertId, 1, srcs, nil)
	free()

	gl.CompileShader(vertId)

	fragId := gl.CreateShader(gl.FRAGMENT_SHADER)

	srcs, free = gl.Strs(fragShader)
	gl.ShaderSource(fragId, 1, srcs, nil)
	free()

	gl.CompileShader(fragId)

	gl.AttachShader(progId, vertId)
	gl.AttachShader(progId, fragId)
	gl.LinkProgram(progId)

	var status int32
	gl.GetProgramiv(progId, gl.LINK_STATUS, &status)
	if status == gl.FALSE {
		var logLength int32
		gl.GetProgramiv(progId, gl.INFO_LOG_LENGTH, &logLength)

		log := strings.Repeat("\x00", int(logLength+1))
		gl.GetProgramInfoLog(progId, logLength, nil, gl.Str(log))

		fmt.Printf("failed to link program: %v", log)
	}

	gl.DeleteShader(vertId)
	gl.DeleteShader(fragId)

	gl.UseProgram(progId)

	model := mgl32.Ident4()
	view := mgl32.LookAt(
		3, 3, 3,
		0, 0, 0,
		0, 1, 0,
	)
	projection := mgl32.Perspective(
		mgl32.DegToRad(45.0),
		float32(app.WindowWidth)/float32(app.WindowHeight),
		0.001, 1024.0,
	)
	mvp := projection.Mul4(view.Mul4(model))

	gl.UniformMatrix4fv(gl.GetUniformLocation(progId, gl.Str("_MVP\x00")), 1, false, &mvp[0])

	// --
	// --
	// --

	return app, err
}

func (app App) Cleanup() {
	app.sdlWindow.Destroy()
	sdl.GL_DeleteContext(app.sdlContext)
}

func (app App) Start() error {
	var evt sdl.Event

	frameDelay := float64(1000.0 / app.TargetFps)
	frameElap := float64(0.0)

	fpsUpdateFrames := 0
	fpsUpdateDelay := float64(250.0)
	fpsUpdateElap := float64(0.0)

	now := func() float64 {
		return float64(time.Now().UnixNano()) / float64(time.Millisecond)
	}

	timeOffset := now()

	app.running = true
	for app.running {
		elapsedTime := now() - timeOffset
		timeOffset = now()

		for evt = sdl.PollEvent(); evt != nil; evt = sdl.PollEvent() {
			switch evt.(type) {
			case *sdl.QuitEvent:
				app.running = false
			}
		}

		app.updateCtx.DeltaTime = float32(elapsedTime / frameDelay)
		app.updateCtx.ElapsedTime = elapsedTime
		app.updateCtx.TotalTime += elapsedTime

		for _, f := range app.EvtUpdate {
			f(&app.updateCtx)
		}

		frameElap += elapsedTime
		if frameDelay <= frameElap {
			gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT)

			for _, f := range app.EvtRender {
				f(&app.renderCtx)
			}

			sdl.GL_SwapWindow(app.sdlWindow)

			frameElap = 0.0
			fpsUpdateFrames += 1
			app.updateCtx.Frame += 1
		}

		fpsUpdateElap += elapsedTime
		if fpsUpdateDelay <= fpsUpdateElap {
			app.updateCtx.CurrentFps = float32(float64(fpsUpdateFrames)/fpsUpdateElap) * 1000.0

			title := fmt.Sprintf("%s - %0.2f", app.WindowTitle, app.updateCtx.CurrentFps)
			app.sdlWindow.SetTitle(title)

			fpsUpdateElap = 0.0
			fpsUpdateFrames = 0
		}

		//if frameDelay - frameElap < 1 {
		//    sdl.Delay(1)
		//}
	}

	sdl.Quit()
	return nil
}

var vertShader = `
#version 330 core

uniform mat4 _MVP;

in layout(location = 0) vec3 _Vertex;
in layout(location = 1) vec3 _Normal;
in layout(location = 2) vec2 _TexCoord;

out vec4 p_Vertex;
out vec4 p_Normal;
out vec2 p_TexCoord;

void main() {
	p_Vertex = _MVP * vec4(_Vertex, 1.0);
	p_Normal = _MVP * vec4(_Normal, 1.0);
	p_TexCoord = _TexCoord;

	gl_Position = _MVP * vec4(_Vertex, 1.0);
}
` + "\x00"

var fragShader = `
#version 330 core

in vec4 p_Vertex;
in vec4 p_Normal;
in vec2 p_TexCoord;

out vec4 o_Color;

void main() {
	o_Color = vec4(1, 0, 0, 1);
}
` + "\x00"
