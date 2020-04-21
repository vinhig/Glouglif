package main

import (
	"github.com/go-gl/gl/v3.3-core/gl"
)

func GenVertexArray() GLVertexArrayObject {
	var vao uint32
	vao = 0
	gl.GenVertexArrays(1, &vao)
	return GLVertexArrayObject(vao)
}

func GenBuffer(vao GLVertexArrayObject, length int, size int, data []float32, attribIndex uint32) {
	vertexArray := uint32(vao)
	// Bind future buffer to this vertex array
	gl.BindVertexArray(vertexArray)

	// Generate the actual buffer
	var buffer uint32
	buffer = 0
	gl.GenBuffers(1, &buffer)
	gl.BindBuffer(gl.ARRAY_BUFFER, buffer)
	gl.BufferData(
		gl.ARRAY_BUFFER,
		length*4,
		gl.Ptr(data),
		gl.STATIC_DRAW)
	gl.EnableVertexAttribArray(attribIndex)
	gl.VertexAttribPointer(uint32(attribIndex), int32(size), gl.FLOAT, false, 0, nil)
}

func DrawArrays(vao GLVertexArrayObject, length int) {
	gl.BindVertexArray(uint32(vao))
	gl.DrawArrays(gl.TRIANGLES, 0, int32(length))
}
