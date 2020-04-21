package main

type Mesh struct {
	Node
	primitive      IPrimitive
	texture        GLTexture
	shadowReceiver bool
	modifier       IModifier
}

func NewMesh(primitive string, texture string) *Mesh {
	mesh := new(Mesh)
	mesh.visible = true
	if primitive == "cube" {
		mesh.primitive = NewCube()
	} else if primitive == "triangle" {
		mesh.primitive = NewTriangle()
	} else {
		mesh.primitive = NewImported(primitive)
	}
	mesh.inode = mesh

	glTexture := NewTexture(texture)
	mesh.texture = glTexture

	return mesh
}

func (mesh *Mesh) SetModifier(modifier IModifier) {
	mesh.modifier = modifier
}

func (mesh *Mesh) GetModifier(modifier IModifier) IModifier{
	return mesh.modifier
}

func (mesh *Mesh) Render(pipeline IPipeline) {
	if mesh.visible {
		if mesh.modifier != nil {
			mesh.modifier.Render(mesh, pipeline)
		} else {
			pipeline.LinkNode(mesh)
			if pipeline.GetType() != SHADOW_PIPELINE {
				pipeline.LinkTexture(mesh.texture, 0)
				if mesh.shadowReceiver {
					pipeline.GetShader().UseInt(1, "shadow_receiver")
				} else {
					pipeline.GetShader().UseInt(0, "shadow_receiver")
				}
			}
			pipeline.DrawPrimitive(mesh.primitive)
		}
		for i := 0; i < len(mesh.children); i++ {
			mesh.children[i].Render(pipeline)
		}
	}
}
