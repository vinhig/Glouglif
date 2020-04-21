package main

import (
	"github.com/go-gl/mathgl/mgl32"
	"github.com/veandco/go-sdl2/sdl"
	"math"
)

type ICamera interface {
	GetView() mgl32.Mat4
	GetProjection() mgl32.Mat4
	GetWorldModel() mgl32.Mat4
}

type PerspectiveCamera struct {
	Node
	view       mgl32.Mat4
	projection mgl32.Mat4
}

func NewPerspectiveCamera() *PerspectiveCamera {
	camera := new(PerspectiveCamera)
	camera.Init()
	camera.visible = true
	camera.inode = camera
	camera.projection = mgl32.Perspective(mgl32.DegToRad(60.0), float32(WIDTH)/float32(HEIGHT), 0.1, 1000.0)
	camera.view = mgl32.LookAtV(camera.position, mgl32.Vec3{0.0, 0.0, 1.0}, mgl32.Vec3{0.0, 1.0, 0.0})
	camera.view = camera.GetWorldModel().Mul4(camera.view)

	return camera
}

func (camera *PerspectiveCamera) GetView() mgl32.Mat4 {
	// if camera.modelNeedUpdate {
	center := mgl32.Vec3{
		float32(math.Cos(float64(camera.rotation.Y())) * math.Cos(float64(camera.rotation.X()))),
		float32(math.Sin(float64(camera.rotation.X()))),
		float32(math.Sin(float64(camera.rotation.Y())) * math.Cos(float64(camera.rotation.X())))}.Add(camera.position)
	up := mgl32.Vec3{0.0, 1.0, 0.0}

	camera.view = mgl32.LookAtV(camera.position, center, up)
	// }

	return camera.view
}

func (camera *PerspectiveCamera) GetProjection() mgl32.Mat4 {
	camera.projection = mgl32.Perspective(mgl32.DegToRad(60.0), float32(WIDTH)/float32(HEIGHT), 0.1, 1000.0)

	return camera.projection
}

type FirstPersonController struct {
	yaw   float32
	pitch float32
}

func (camera *FirstPersonController) Act(node INode, delta float32) {
	x, y, _ := sdl.GetRelativeMouseState()
	camera.yaw += float32(x) / 10.0

	pitch := camera.pitch - float32(y)/10.0
	if pitch > -60 && pitch < 60 {
		camera.pitch = pitch
	} else {
		camera.pitch = mgl32.Clamp(camera.pitch, -60, 60)
	}

	front := mgl32.Vec3{0.0, 0.0, 0.0}
	left := mgl32.Vec3{0.0, 0.0, 0.0}

	if input.IsKeyDown(Z) {
		front = mgl32.Vec3{1.0, 0.0, 0.0}
		left = mgl32.Vec3{0.0, 0.0, 1.0}
	} else if input.IsKeyDown(S) {
		front = mgl32.Vec3{-1.0, 0.0, 0.0}
		left = mgl32.Vec3{0.0, 0.0, -1.0}
	}

	if input.IsKeyDown(Q) {
		front = front.Add(mgl32.Vec3{0.0, 0.0, -1.0})
		left = left.Add(mgl32.Vec3{1.0, 0.0, 0.0})
	} else if input.IsKeyDown(D) {
		front = front.Add(mgl32.Vec3{0.0, 0.0, 1.0})
		left = left.Add(mgl32.Vec3{-1.0, 0.0, 0.0})
	}

	front = front.Mul(float32(math.Sin(float64(mgl32.DegToRad(camera.yaw+90)))) * 15 * delta)
	left = left.Mul(float32(math.Cos(float64(mgl32.DegToRad(camera.yaw+90)))) * -15 * delta)

	node.(*Node).position = node.(*Node).position.Add(front.Add(left))
	node.Rotate(camera.pitch, camera.yaw, 0.0)

	network.Register(&CameraMovementAction{
		playerID: network.GetPlayerID(),
		position: node.(*Node).position,
		rotation: node.(*Node).rotation,
	})
}
