package main

type IModifier interface {
	Render(mesh *Mesh, pipeline IPipeline)
}