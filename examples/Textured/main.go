package main

import (
	"log"

	"github.com/WhoBrokeTheBuild/GoDusk/dusk"
)

var model *dusk.Model

func render(ctx *dusk.RenderContext) {
	model.Render()
}

func main() {
	app, err := dusk.NewApp()
	if err != nil {
		log.Fatal(err)
	}
	defer app.Cleanup()

	app.EvtRender = append(app.EvtRender, render)

	model, err = dusk.NewModelFromFile("cube.obj")
	if err != nil {
		log.Fatal(err)
	}

	app.Start()
}
