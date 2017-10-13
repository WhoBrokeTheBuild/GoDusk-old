package main

import (
	"log"
	"runtime"

	"github.com/WhoBrokeTheBuild/GoDusk/dusk"
	"github.com/go-gl/gl/v4.1-core/gl"
	"github.com/go-gl/mathgl/mgl32"
)

var camera *dusk.Camera
var shader *dusk.Shader
var model *dusk.Model

var rotation = float32(0)

var update = func(data interface{}) {
	ctx := data.(*dusk.UpdateContext)

	rotation += 1.0 * ctx.DeltaTime
}

var render = func(data interface{}) {
	var glerr uint32
	log.SetFlags(log.LstdFlags | log.Lshortfile)

	//ctx := data.(*dusk.RenderContext)

	model.Transform = model.Transform.Mul4(mgl32.HomogRotate3D(mgl32.DegToRad(rotation), mgl32.Vec3{0, 1, 0}))
	rotation = 0.0

	mvp := camera.Proj.Mul4(camera.View.Mul4(model.Transform))
	gl.UniformMatrix4fv(shader.GetUniformLocation("_Model"), 1, false, &model.Transform[0])
	if glerr = gl.GetError(); glerr > 0 {
		log.Printf("gl.GetError returned %v", glerr)
	}
	gl.UniformMatrix4fv(shader.GetUniformLocation("_View"), 1, false, &camera.View[0])
	if glerr = gl.GetError(); glerr > 0 {
		log.Printf("gl.GetError returned %v", glerr)
	}
	gl.UniformMatrix4fv(shader.GetUniformLocation("_Proj"), 1, false, &camera.Proj[0])
	if glerr = gl.GetError(); glerr > 0 {
		log.Printf("gl.GetError returned %v", glerr)
	}
	gl.UniformMatrix4fv(shader.GetUniformLocation("_MVP"), 1, false, &mvp[0])
	if glerr = gl.GetError(); glerr > 0 {
		log.Printf("gl.GetError returned %v", glerr)
	}

	eye := mgl32.Vec3{2, 2, 2}

	gl.Uniform3fv(shader.GetUniformLocation("_LightPos"), 1, &eye[0])
	if glerr = gl.GetError(); glerr > 0 {
		log.Printf("gl.GetError returned %v", glerr)
	}
	gl.Uniform3fv(shader.GetUniformLocation("_ViewPos"), 1, &eye[0])
	if glerr = gl.GetError(); glerr > 0 {
		log.Printf("gl.GetError returned %v", glerr)
	}

	model.Render(shader)
}

func main() {
	runtime.LockOSThread()

	app, err := dusk.NewApp()
	if err != nil {
		dusk.LogError("%v", err)
		return
	}
	defer app.Cleanup()

	dusk.LogInfo("GL_NO_ERROR %v", gl.NO_ERROR)
	dusk.LogInfo("GL_INVALID_ENUM %v", gl.INVALID_ENUM)
	dusk.LogInfo("GL_INVALID_VALUE %v", gl.INVALID_VALUE)
	dusk.LogInfo("GL_INVALID_OPERATION %v", gl.INVALID_OPERATION)
	dusk.LogInfo("GL_INVALID_FRAMEBUFFER_OPERATION %v", gl.INVALID_FRAMEBUFFER_OPERATION)
	dusk.LogInfo("GL_OUT_OF_MEMORY %v", gl.OUT_OF_MEMORY)

	app.AssetFunction = Asset
	app.EvtUpdate.Subscribe(&update)
	app.EvtRender.Subscribe(&render)

	camera = dusk.NewCamera(&app, 45.0, 0.1, 100.0)
	defer camera.Cleanup(&app)

	camera.SetPosition(mgl32.Vec3{2, 2, 2})
	camera.SetDirection(mgl32.Vec3{-1, -1, -1})

	shader, err = dusk.NewShader(&app, "assets/default.vs.glsl", "assets/default.fs.glsl")
	if err != nil {
		dusk.LogError("%v", err)
		return
	}
	shader.Use()

	model, err = dusk.NewModelFromFile(&app, "assets/globe/globe.obj")
	if err != nil {
		dusk.LogError("%v", err)
		return
	}
	defer model.Cleanup()

	app.Start()
}
