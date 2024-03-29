package dusk

import (
	"fmt"
	"io/ioutil"
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

	EvtUpdate *Event // ctx *UpdateContext
	EvtRender *Event // ctx *RenderContext
	EvtResize *Event // size mgl32.Vec3

	AssetFunction func(string) ([]byte, error)

	updateCtx  UpdateContext
	renderCtx  RenderContext
	sdlWindow  *sdl.Window
	sdlContext sdl.GLContext
	running    bool
}

func AssetFromFile(filename string) ([]byte, error) {
	return ioutil.ReadFile(filename)
}

func NewApp() (App, error) {
	var err error

	app := App{
		WindowTitle:   "Dusk",
		WindowWidth:   640,
		WindowHeight:  480,
		TargetFps:     60,
		EvtUpdate:     NewEvent(),
		EvtRender:     NewEvent(),
		EvtResize:     NewEvent(),
		AssetFunction: AssetFromFile,
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
			case *sdl.WindowEvent:
				windowEvt := evt.(*sdl.WindowEvent)
				switch windowEvt.Event {
				case sdl.WINDOWEVENT_RESIZED:

					gl.Viewport(0, 0, windowEvt.Data1, windowEvt.Data2)
					app.EvtResize.Call(mgl32.Vec2{
						float32(windowEvt.Data1),
						float32(windowEvt.Data2),
					})
				}

			}
		}

		app.updateCtx.DeltaTime = float32(elapsedTime / frameDelay)
		app.updateCtx.ElapsedTime = elapsedTime
		app.updateCtx.TotalTime += elapsedTime

		app.EvtUpdate.Call(&app.updateCtx)

		frameElap += elapsedTime
		if frameDelay <= frameElap {
			gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT)

			app.EvtRender.Call(&app.renderCtx)

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
