package main

import (
	"log"

	"github.com/WhoBrokeTheBuild/GoDusk/dusk"
)

func update(ctx *dusk.UpdateContext) {
	//dusk.LogInfo("Update FPS: %.2f, Delta: %.3f, Elapsed: %.2f, Total: %.2f", ctx.CurrentFps, ctx.DeltaTime, ctx.ElapsedTime, ctx.TotalTime)
}

func main() {
	app, err := dusk.NewApp()
	if err != nil {
		log.Fatal(err)
	}

	app.EvtUpdate = append(app.EvtUpdate, update)

	model, err := dusk.NewModelFromFile("cube.obj")
	if err != nil {
		log.Fatal(err)
	}

	_ = model

	app.Start()
}
