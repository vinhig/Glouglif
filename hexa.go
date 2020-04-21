package main

import "github.com/go-gl/mathgl/mgl32"

type HexaModifier struct {
	sizeX    int
	sizeY    int
	width    float32
	height   float32
	position mgl32.Vec3
}

func (modifier *HexaModifier) Render(mesh *Mesh, pipeline IPipeline) {
	if pipeline.GetType() != SHADOW_PIPELINE {
		for x := 0; x < modifier.sizeX; x++ {
			// fmt.Printf("coucou %d", x)
			for y := 0; y < modifier.sizeY; y++ {
				if y%2 == 0 {
					mesh.Translate(modifier.position.X()+modifier.width*float32(x)*2+modifier.width, modifier.position.Y()+0.0001*float32(x), modifier.position.Y()+modifier.height*float32(y))
				} else {
					mesh.Translate(modifier.position.X()+modifier.width*float32(x)*2, modifier.position.Y(), modifier.position.Y()+modifier.height*float32(y))
				}
				pipeline.LinkNode(mesh)
				pipeline.LinkTexture(mesh.texture, 0)
				pipeline.GetShader().UseInt(1, "shadow_receiver")
				pipeline.DrawPrimitive(mesh.primitive)
			}
		}
	}
}
