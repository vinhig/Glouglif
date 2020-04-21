package main

import (
	"bytes"

	"github.com/sheenobu/go-obj/obj"
)

type IPrimitive interface {
	GetVao() GLVertexArrayObject
	GetLength() int
}

func NewCube() *Cube {
	vertices := []float32{
		//  X, Y, Z, U, V
		// Bottom
		-1.0, -1.0, -1.0,
		1.0, -1.0, -1.0,
		-1.0, -1.0, 1.0,
		1.0, -1.0, -1.0,
		1.0, -1.0, 1.0,
		-1.0, -1.0, 1.0,

		// Top
		-1.0, 1.0, -1.0,
		-1.0, 1.0, 1.0,
		1.0, 1.0, -1.0,
		1.0, 1.0, -1.0,
		-1.0, 1.0, 1.0,
		1.0, 1.0, 1.0,

		// Front
		-1.0, -1.0, 1.0,
		1.0, -1.0, 1.0,
		-1.0, 1.0, 1.0,
		1.0, -1.0, 1.0,
		1.0, 1.0, 1.0,
		-1.0, 1.0, 1.0,

		// Back
		-1.0, -1.0, -1.0,
		-1.0, 1.0, -1.0,
		1.0, -1.0, -1.0,
		1.0, -1.0, -1.0,
		-1.0, 1.0, -1.0,
		1.0, 1.0, -1.0,

		// Left
		-1.0, -1.0, 1.0,
		-1.0, 1.0, -1.0,
		-1.0, -1.0, -1.0,
		-1.0, -1.0, 1.0,
		-1.0, 1.0, 1.0,
		-1.0, 1.0, -1.0,

		// Right
		1.0, -1.0, 1.0,
		1.0, -1.0, -1.0,
		1.0, 1.0, -1.0,
		1.0, -1.0, 1.0,
		1.0, 1.0, -1.0,
		1.0, 1.0, 1.0,
	}

	uvs := []float32{
		0.0, 0.0,
		1.0, 0.0,
		0.0, 1.0,
		1.0, 0.0,
		1.0, 1.0,
		0.0, 1.0,

		0.0, 0.0,
		0.0, 1.0,
		1.0, 0.0,
		1.0, 0.0,
		0.0, 1.0,
		1.0, 1.0,

		1.0, 0.0,
		0.0, 0.0,
		1.0, 1.0,
		0.0, 0.0,
		0.0, 1.0,
		1.0, 1.0,

		0.0, 0.0,
		0.0, 1.0,
		1.0, 0.0,
		1.0, 0.0,
		0.0, 1.0,
		1.0, 1.0,

		0.0, 1.0,
		1.0, 0.0,
		0.0, 0.0,
		0.0, 1.0,
		1.0, 1.0,
		1.0, 0.0,

		1.0, 1.0,
		1.0, 0.0,
		0.0, 0.0,
		1.0, 1.0,
		0.0, 0.0,
		0.0, 1.0,
	}

	vao := GenVertexArray()
	// panic(len(vertices))
	GenBuffer(vao, len(vertices), 3, vertices, 0)
	GenBuffer(vao, len(uvs), 2, uvs, 1)

	prim := Cube{
		vao: vao,
	}

	return &prim
}

type Cube struct {
	vao GLVertexArrayObject
}

func (cube *Cube) GetVao() GLVertexArrayObject {
	return cube.vao
}

func (cube *Cube) GetLength() int {
	return 108 / 3
}

type Triangle struct {
	vao GLVertexArrayObject
}

func NewTriangle() *Triangle {
	vertices := []float32{
		-0.5, -0.5, 0,
		0.5, -0.5, 0,
		0.0, 0.5, 0,
	}

	uvs := []float32{
		0, 0,
		1, 0,
		.5, 1,
	}

	vao := GenVertexArray()
	GenBuffer(vao, 9, 3, vertices, 0)
	GenBuffer(vao, 9, 2, uvs, 1)

	prim := Triangle{
		vao: vao,
	}

	return &prim
}

func (triangle *Triangle) GetVao() GLVertexArrayObject {
	return triangle.vao
}

func (triangle *Triangle) GetLength() int {
	return 9 / 3
}

type Imported struct {
	vao    GLVertexArrayObject
	length int
}

func NewImported(path string) *Imported {
	obj, err := obj.NewReader(bytes.NewReader(ReadFile(path))).Read()
	if err != nil {
		panic(err)
	}

	var vertices, uvs, normals []float32
	// Iterate through face and vertex and feed vertices/uvs/normals
	for _, f := range obj.Faces {
		for _, p := range f.Points {
			vx := p.Vertex

			nx := float32(0.0)
			ny := float32(0.0)
			nz := float32(0.0)
			if p.Normal != nil {
				nx = float32(p.Normal.X)
				ny = float32(p.Normal.Y)
				nz = float32(p.Normal.Z)
			}

			u := float32(0.0)
			v := float32(0.0)
			if p.Texture != nil {
				u = float32(p.Texture.U)
				v = float32(p.Texture.V)
			}

			// Feed vertices/uvs/normals
			vertices = append(vertices, float32(vx.X), float32(vx.Y), float32(vx.Z))
			uvs = append(uvs, u, v)
			normals = append(normals, nx, ny, nz)
		}
	}

	vao := GenVertexArray()

	GenBuffer(vao, len(vertices), 3, vertices, 0)
	GenBuffer(vao, len(uvs), 2, uvs, 1)
	GenBuffer(vao, len(normals), 3, normals, 2)

	prim := Imported{
		length: len(vertices) / 3,
		vao: vao,
	}

	return &prim
}

func (obj *Imported) GetVao() GLVertexArrayObject {
	return obj.vao
}

func (obj *Imported) GetLength() int {
	return obj.length
}
