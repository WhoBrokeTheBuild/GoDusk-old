package main

import (
	"github.com/WhoBrokeTheBuild/GoDusk/dusk"
	"github.com/go-gl/gl/v4.1-core/gl"
	"github.com/go-gl/mathgl/mgl32"
)

var model *dusk.Model

func render(ctx *dusk.RenderContext) {
	model.Render()
}

func main() {
	app, err := dusk.NewApp()
	if err != nil {
		dusk.LogError("%v", err)
        return
	}
	defer app.Cleanup()

	app.EvtRender = append(app.EvtRender, render)

	shader, err := dusk.NewShader("assets/default.vs.glsl", "assets/default.fs.glsl")
	if err != nil {
		dusk.LogError("%v", err)
        return
	}
	shader.Use()

	modelMat := mgl32.Ident4()
	viewMat := mgl32.LookAt(
		3, 3, 3,
		0, 0, 0,
		0, 1, 0,
	)
	projMat := mgl32.Perspective(
		mgl32.DegToRad(45.0),
		float32(app.WindowWidth)/float32(app.WindowHeight),
		0.001, 1024.0,
	)
	mvp := projMat.Mul4(viewMat.Mul4(modelMat))

	gl.UniformMatrix4fv(shader.GetUniformLocation("_MVP"), 1, false, &mvp[0])

	model, err = dusk.NewModelFromFile("assets/cube.obj")
	if err != nil {
		dusk.LogError("%v", err)
        return
	}

	app.Start()
}
