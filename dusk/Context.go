package dusk

type UpdateContext struct {
	DeltaTime   float32
	ElapsedTime float64
	TotalTime   float64
	CurrentFps  float32
	Frame       uint64
}

type RenderContext struct {
}
