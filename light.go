package main

import "github.com/go-gl/mathgl/mgl32"

type DirectionalLight struct {
	Node
	view       mgl32.Mat4
	projection mgl32.Mat4
}

func NewDirectionalLight() *DirectionalLight {
	light := new(DirectionalLight)
	light.Init()
	light.inode = light

	light.projection = mgl32.Ortho(-10, 10, -10, 10, -10, 20)
	light.view = mgl32.LookAtV(light.position, mgl32.Vec3{0.0, 0.0, 1.0}, mgl32.Vec3{0.0, 1.0, 0.0})

	return light
}

func (light *DirectionalLight) GetView() mgl32.Mat4 {
	if light.modelNeedUpdate {
		light.view = light.GetWorldModel().Mul4(mgl32.LookAtV(light.position, mgl32.Vec3{0.0, 0.0, 1.0}, mgl32.Vec3{0.0, 1.0, 0.0}))
	}
	return light.view
}

func (light *DirectionalLight) GetProjection() mgl32.Mat4 {
	return light.projection
}