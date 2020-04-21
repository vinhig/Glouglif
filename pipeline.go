package main

import (
	"github.com/go-gl/gl/v3.3-core/gl"
	"github.com/go-gl/mathgl/mgl32"
)

type PipelineType int
type GLTexture uint32
type GLVertexArrayObject uint32

const (
	DIFFUSE_PIPELINE = iota
	SHADOW_PIPELINE
)

var bias = mgl32.Mat4{
0.5, 0.0, 0.0, 0.0,
0.0, 0.5, 0.0, 0.0,
0.0, 0.0, 0.5, 0.0,
0.5, 0.5, 0.5, 1.0}

type GLFrameBuffer uint32

type IPipeline interface {
	GetShader() IShader
	Begin()
	End()
	GetType() PipelineType
	LinkNode(node INode)
	LinkTexture(texture GLTexture, binding int)
	DrawPrimitive(prim IPrimitive)
	GetFrameBuffer() GLFrameBuffer
	GetTexture() GLTexture
}

type DiffusePipeline struct {
	basic       IShader
	quad        IShader
	vao         GLVertexArrayObject
	frameBuffer GLFrameBuffer
	texture     GLTexture
}

func NewDiffusePipeline() *DiffusePipeline {
	var frameBuffer uint32
	gl.GenFramebuffers(1, &frameBuffer)
	gl.BindFramebuffer(gl.FRAMEBUFFER, frameBuffer)

	var texture uint32
	gl.GenTextures(1, &texture)
	gl.GenTextures(1, &texture)
	gl.BindTexture(gl.TEXTURE_2D, texture)
	gl.TexImage2D(gl.TEXTURE_2D, 0, gl.RGB, int32(WIDTH), int32(HEIGHT), 0, gl.RGB, gl.UNSIGNED_BYTE, nil)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MAG_FILTER, gl.NEAREST)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MIN_FILTER, gl.NEAREST)

	var depthBuffer uint32
	gl.GenRenderbuffers(1, &depthBuffer)
	gl.BindRenderbuffer(gl.RENDERBUFFER, depthBuffer)
	gl.RenderbufferStorage(gl.RENDERBUFFER, gl.DEPTH_COMPONENT, int32(WIDTH), int32(HEIGHT))
	gl.FramebufferRenderbuffer(gl.FRAMEBUFFER, gl.DEPTH_ATTACHMENT, gl.RENDERBUFFER, depthBuffer)

	gl.FramebufferTexture(gl.FRAMEBUFFER, gl.COLOR_ATTACHMENT0, texture, 0)
	gl.DrawBuffer(gl.COLOR_ATTACHMENT0)

	vertices := []float32{
		0.0, 0.0,
		0.0, -2.0,
		2.0, 0.0,

		0.0, -2.0,
		2.0, -2.0,
		2.0, 0.0,
	}

	uvs := []float32{
		0.0, 0.0,
		0.0, -1.0,
		1.0, 0.0,

		0.0, -1.0,
		1.0, -1.0,
		1.0, 0.0,
	}

	vao := GenVertexArray()

	GenBuffer(vao, len(vertices), 2, vertices, 0)
	GenBuffer(vao, len(uvs), 2, uvs, 1)

	program := NewBasicShader()
	quad := NewSpriteShader()
	pipeline := DiffusePipeline{
		basic:       program,
		quad:        quad,
		vao:         GLVertexArrayObject(vao),
		frameBuffer: GLFrameBuffer(frameBuffer),
		texture:     GLTexture(texture),
	}

	return &pipeline
}

func (pipeline *DiffusePipeline) Begin() {
	gl.BindFramebuffer(gl.FRAMEBUFFER, uint32(pipeline.frameBuffer))
	gl.Viewport(0, 0, int32(WIDTH), int32(HEIGHT))

	gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT)
	gl.CullFace(gl.BACK)

	pipeline.basic.Bind()
}

func (pipeline *DiffusePipeline) End() {
	gl.BindFramebuffer(gl.FRAMEBUFFER, 0)
	gl.Viewport(0, 0, int32(WIDTH), int32(HEIGHT))
	gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT)

	pipeline.quad.Bind()
	pipeline.LinkTexture(pipeline.texture, 0)
	gl.BindVertexArray(uint32(pipeline.vao))
	gl.DrawArrays(gl.TRIANGLES, 0, 9)
}

func (pipeline *DiffusePipeline) GetType() PipelineType {
	return DIFFUSE_PIPELINE
}

func (pipeline *DiffusePipeline) DrawPrimitive(prim IPrimitive) {
	DrawArrays(prim.GetVao(), prim.GetLength())
}

func (pipeline *DiffusePipeline) LinkNode(node INode) {
	pipeline.basic.UseMatrix4(node.GetWorldModel(), "i_model")
}

func (pipeline *DiffusePipeline) LinkTexture(texture GLTexture, binding int) {
	pipeline.basic.UseTexture(texture, binding)
}

func (pipeline *DiffusePipeline) LinkCamera(camera ICamera) {
	pipeline.basic.UseMatrix4(camera.GetProjection(), "i_projection")
	pipeline.basic.UseMatrix4(camera.GetView(), "i_view")
}

func (pipeline *DiffusePipeline) GetFrameBuffer() GLFrameBuffer {
	return GLFrameBuffer(pipeline.frameBuffer)
}

func (pipeline *DiffusePipeline) GetTexture() GLTexture {
	return GLTexture(pipeline.texture)
}

func (pipeline *DiffusePipeline) GetShader() IShader {
	return pipeline.basic
}


type ShadowPipeline struct {
	shadow      IShader
	texture     GLTexture
	frameBuffer GLFrameBuffer
}

func NewShadowPipeline() *ShadowPipeline {
	shadow := NewShadowShader()

	var frameBuffer uint32
	gl.GenFramebuffers(1, &frameBuffer)
	gl.BindFramebuffer(gl.FRAMEBUFFER, frameBuffer)

	var texture uint32
	gl.GenTextures(1, &texture)
	gl.BindTexture(gl.TEXTURE_2D, texture)
	gl.TexImage2D(gl.TEXTURE_2D, 0, gl.DEPTH_COMPONENT16, 2048, 2048, 0, gl.DEPTH_COMPONENT, gl.FLOAT, nil)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MAG_FILTER, gl.NEAREST)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MIN_FILTER, gl.NEAREST)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_WRAP_S, gl.CLAMP_TO_EDGE)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_WRAP_T, gl.CLAMP_TO_EDGE)

	gl.FramebufferTexture(gl.FRAMEBUFFER, gl.DEPTH_ATTACHMENT, texture, 0)
	gl.DrawBuffer(gl.NONE)

	pipeline := ShadowPipeline{
		frameBuffer: GLFrameBuffer(frameBuffer),
		texture:     GLTexture(texture),
		shadow:      shadow,
	}

	return &pipeline
}

func (pipeline *ShadowPipeline) Begin() {
	gl.BindFramebuffer(gl.FRAMEBUFFER, uint32(pipeline.frameBuffer))
	gl.Viewport(0, 0, 2048, 2048)

	gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT)
	gl.CullFace(gl.FRONT)

	pipeline.shadow.Bind()
}

func (pipeline *ShadowPipeline) End() {

}

func (pipeline *ShadowPipeline) GetType() PipelineType {
	return SHADOW_PIPELINE
}

func (pipeline *ShadowPipeline) DrawPrimitive(prim IPrimitive) {
	DrawArrays(prim.GetVao(), prim.GetLength())
}

func (pipeline *ShadowPipeline) LinkNode(node INode) {
	pipeline.shadow.UseMatrix4(node.GetWorldModel(), "i_model")
}

func (pipeline *ShadowPipeline) LinkTexture(texture GLTexture, binding int) {
	// It's a depth rendering
	// No need to link texture
	panic(texture)
}

func (pipeline *ShadowPipeline) LinkCamera(camera ICamera) {
	pipeline.shadow.UseMatrix4(camera.GetProjection(), "i_projection")
	pipeline.shadow.UseMatrix4(camera.GetView(), "i_view")
	// pipeline.shadow.UseMatrix4(camera.GetWorldModel(), "i_camera_model")
	// pipeline.shadow.UseMatrix4(camera.GetWorldModel(), "i_model")
}

func (pipeline *ShadowPipeline) GetFrameBuffer() GLFrameBuffer {
	return GLFrameBuffer(pipeline.frameBuffer)
}

func (pipeline *ShadowPipeline) GetTexture() GLTexture {
	return GLTexture(pipeline.texture)
}

func (pipeline *ShadowPipeline) GetShader() IShader {
	return pipeline.shadow
}
