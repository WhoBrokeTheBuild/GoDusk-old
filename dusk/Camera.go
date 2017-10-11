package dusk

import (
	"github.com/go-gl/mathgl/mgl32"
)

type Camera struct {
	View mgl32.Mat4
	Proj mgl32.Mat4

	pos mgl32.Vec3
	dir mgl32.Vec3
	up  mgl32.Vec3

	fov          float32
	near         float32
	far          float32
	aspectWidth  float32
	aspectHeight float32

	resize func(data interface{})
}

func NewCamera(app *App, fov float32, near float32, far float32) *Camera {
	camera := &Camera{
		pos:          mgl32.Vec3{0, 0, 0},
		dir:          mgl32.Vec3{0, 0, 0},
		up:           mgl32.Vec3{0, 1, 0},
		fov:          fov,
		near:         near,
		far:          far,
		aspectWidth:  float32(app.WindowWidth),
		aspectHeight: float32(app.WindowHeight),
	}
	camera.resize = func(data interface{}) {
		size := data.(mgl32.Vec2)

		camera.aspectWidth = size.X()
		camera.aspectHeight = size.Y()
		camera.calculateProj()
	}

	app.EvtResize.Subscribe(&camera.resize)

	camera.calculateView()
	camera.calculateProj()

	return camera
}

func (camera *Camera) Cleanup(app *App) {
	app.EvtResize.Unsubscribe(&camera.resize)
}

func (camera *Camera) calculateView() {
	center := camera.pos.Add(camera.dir)
	camera.View = mgl32.LookAt(
		camera.pos.X(), camera.pos.Y(), camera.pos.Z(),
		center.X(), center.Y(), center.Z(),
		camera.up.X(), camera.up.Y(), camera.up.Z(),
	)
}

func (camera *Camera) calculateProj() {
	camera.Proj = mgl32.Perspective(
		mgl32.DegToRad(camera.fov),
		camera.aspectWidth/camera.aspectHeight,
		camera.near, camera.far,
	)
}

func (camera *Camera) SetPosition(pos mgl32.Vec3) {
	camera.pos = pos
	camera.calculateView()
}

func (camera *Camera) SetDirection(dir mgl32.Vec3) {
	camera.dir = dir
	camera.calculateView()
}

func (camera *Camera) SetUp(up mgl32.Vec3) {
	camera.up = up
	camera.calculateView()
}
