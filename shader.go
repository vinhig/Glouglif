package main

import (
	"fmt"
	"github.com/go-gl/gl/v3.3-core/gl"
	"github.com/go-gl/mathgl/mgl32"
	"strings"
)

type GLProgram uint32
type GLShader uint32

type IShader interface {
	Bind()
	UseTexture(texture GLTexture, index int)
	UseMatrix4(mat4 mgl32.Mat4, input string)
	UseVector3(vec3 mgl32.Vec3, input string)
	UseVector4(vec3 mgl32.Vec3, input string)
	UseInt(vec int, input string)
	GetProgram() GLProgram
}

type ShadowShader struct {
	program GLProgram
}

func NewShadowShader() *ShadowShader {
	vertexSource :=
		`#version 330 core //shadow vertex
#extension GL_ARB_shading_language_420pack : enable

layout(location = 0) in vec3 position;

uniform mat4 i_projection;
uniform mat4 i_view;
uniform mat4 i_model;

void main() { gl_Position = i_projection * i_view * i_model * vec4(position, 1.0); }`

	fragmentSource :=
		`#version 330 core //shadow fragment
#extension GL_ARB_shading_language_420pack : enable

// Output data
// layout(location = 0) out float fragmentdepth;

void main() {
  // Not really needed, OpenGL does it anyway
  // fragmentdepth = 10000000000.0;
  // fragmentdepth = 0.0;
}`
	vertexShader, err := compileShader(vertexSource, gl.VERTEX_SHADER)
	if err != nil {
		panic(err)
	}

	fragmentShader, err := compileShader(fragmentSource, gl.FRAGMENT_SHADER)
	if err != nil {
		panic(err)
	}

	program := gl.CreateProgram()
	gl.AttachShader(program, uint32(vertexShader))
	gl.AttachShader(program, uint32(fragmentShader))
	gl.LinkProgram(program)

	gl.DeleteShader(uint32(vertexShader))
	gl.DeleteShader(uint32(fragmentShader))

	shader := ShadowShader{program: GLProgram(program)}

	return &shader
}

func (shader *ShadowShader) Bind() {
	gl.UseProgram(uint32(shader.program))
}

func (shader *ShadowShader) UseTexture(texture GLTexture, index int) {
	panic(index)
}

func (shader *ShadowShader) UseMatrix4(mat4 mgl32.Mat4, input string) {
	inputID := gl.GetUniformLocation(uint32(shader.program), gl.Str(input+"\x00"))
	gl.UniformMatrix4fv(inputID, 1, false, &mat4[0])
}

func (shader *ShadowShader) UseVector3(vec3 mgl32.Vec3, input string) {

}

func (shader *ShadowShader) UseVector4(vec3 mgl32.Vec3, input string) {

}

func (shader *ShadowShader) UseInt(vec int, input string) {
	inputID := gl.GetUniformLocation(uint32(shader.program), gl.Str(input+"\x00"))
	gl.Uniform1i(inputID, int32(vec))
}

func (shader *ShadowShader) GetProgram() GLProgram {
	return shader.program
}

type SpriteShader struct {
	program GLProgram
}

func NewSpriteShader() *SpriteShader {
	vertexSource :=
		`#version 330

layout(location = 0) in vec2 position;
layout(location = 1) in vec2 uv;

out vec2 o_uv;
out vec2 o_shadow_coords;

void main() {
	gl_Position = vec4(position - vec2(1.0, -1.0), 0.0, 1.0);
	o_uv = uv;
}`

	fragmentSource :=
		`#version 330
#extension GL_ARB_shading_language_420pack : enable

out vec4 frag_color;
layout(binding = 0) uniform sampler2D diffuse;

in vec2 o_uv;

void main() {
  frag_color = texture(diffuse, o_uv).rgba;
}`

	vertexShader, err := compileShader(vertexSource, gl.VERTEX_SHADER)
	if err != nil {
		panic(err)
	}

	fragmentShader, err := compileShader(fragmentSource, gl.FRAGMENT_SHADER)
	if err != nil {
		panic(err)
	}

	program := gl.CreateProgram()
	gl.AttachShader(program, uint32(vertexShader))
	gl.AttachShader(program, uint32(fragmentShader))
	gl.LinkProgram(program)

	gl.DeleteShader(uint32(vertexShader))
	gl.DeleteShader(uint32(fragmentShader))

	shader := SpriteShader{program: GLProgram(program)}

	return &shader
}

type BasicShader struct {
	program GLProgram
}

func (shader *SpriteShader) Bind() {
	gl.UseProgram(uint32(shader.program))
}

func (shader *SpriteShader) UseTexture(texture GLTexture, index int) {
	gl.ActiveTexture(gl.TEXTURE0)
	gl.BindTexture(gl.TEXTURE_2D, uint32(texture))
}

func (shader *SpriteShader) UseMatrix4(mat4 mgl32.Mat4, input string) {
	inputID := gl.GetUniformLocation(uint32(shader.program), gl.Str(input+"\x00"))
	gl.UniformMatrix4fv(inputID, 1, false, &mat4[0])
}

func (shader *SpriteShader) UseVector3(vec3 mgl32.Vec3, input string) {

}

func (shader *SpriteShader) UseVector4(vec3 mgl32.Vec3, input string) {

}

func (shader *SpriteShader) UseInt(vec int, input string) {
	inputID := gl.GetUniformLocation(uint32(shader.program), gl.Str(input+"\x00"))
	gl.Uniform1i(inputID, int32(vec))
}

func (shader *SpriteShader) GetProgram() GLProgram {
	return shader.program
}

func NewBasicShader() *BasicShader {
	vertexSource :=
		`#version 330
#extension GL_ARB_shading_language_420pack : enable

layout(location = 0) in vec3 position;
layout(location = 1) in vec2 uv;

uniform mat4 i_projection;
uniform mat4 i_view;
uniform mat4 i_model;

uniform mat4 i_light_projection;
uniform mat4 i_light_view;
uniform mat4 i_light_model;

uniform mat4 i_bias;

out vec4 o_color;
out vec2 o_uv;
out vec4 o_shadow_coords;

void main() {
	o_color = i_projection * i_view * i_model * vec4(position, 1.0);
	gl_Position = i_projection *  i_view * i_model * vec4(position, 1.0);

	o_uv = uv;
	o_shadow_coords = i_bias * i_light_projection * i_light_view * i_model * vec4(position, 1.0);
}`

	fragmentSource :=
		`#version 330
#extension GL_ARB_shading_language_420pack : enable

//vec2 poisson[4] = vec2[](
//	vec2(-0.94201624, -0.39906216 ),
//	vec2(0.94558609, -0.76890725 ),
//	vec2(-0.094184101, -0.92938870 ),
//	vec2(0.34495938, 0.29387760 )
//);

vec2 triton[8] = vec2[](
	vec2(0.1, 0),
	vec2(0.0866, 0.05),
	vec2(0.0707, 0.0707),
	vec2(0.05, 0.05),
	vec2(0, 0.1),
	vec2(-0.05,0.0866),
	vec2(-0.0707,0.05),
	vec2(-0.1,0)
);

layout(binding = 0) uniform sampler2D diffuse;
layout(binding = 1) uniform sampler2D shadow_map;

out vec4 frag_color;

uniform vec4 i_color;
uniform int shadow_receiver;

in vec4 o_color;
in vec2 o_uv;
in vec4 o_shadow_coords;

void main() {
	float shadow = 1.0;
	/*if (texture(shadow_map, o_shadow_coords.xy).r < o_shadow_coords.z - 0.010) {
		shadow = 0.1;
	}*/
	if (shadow_receiver == 1) {
		for (int i = 0; i < 8; i++) {
			if (texture(shadow_map, o_shadow_coords.xy + triton[i]/200.0).r < o_shadow_coords.z - 0.001) {
				shadow -= 0.1;
			}
		}
	}
	// frag_color = vec4(o_shadow_coords.xy, 1.0, 1.0);
	frag_color = texture(diffuse, o_uv).rgba;
	frag_color = vec4(frag_color.xyz * shadow, frag_color.w);
}`
	vertexShader, err := compileShader(vertexSource, gl.VERTEX_SHADER)
	if err != nil {
		panic(err)
	}

	fragmentShader, err := compileShader(fragmentSource, gl.FRAGMENT_SHADER)
	if err != nil {
		panic(err)
	}

	program := gl.CreateProgram()
	gl.AttachShader(program, uint32(vertexShader))
	gl.AttachShader(program, uint32(fragmentShader))
	gl.LinkProgram(program)

	gl.DeleteShader(uint32(vertexShader))
	gl.DeleteShader(uint32(fragmentShader))

	shader := BasicShader{program: GLProgram(program)}

	return &shader
}

func (shader *BasicShader) Bind() {
	gl.UseProgram(uint32(shader.program))
}

func (shader *BasicShader) UseTexture(texture GLTexture, index int) {
	if index == 0 {
		gl.ActiveTexture(gl.TEXTURE0)
	} else if index == 1 {
		gl.ActiveTexture(gl.TEXTURE1)
	} else if index == 2 {
		gl.ActiveTexture(gl.TEXTURE2)
	}
	gl.BindTexture(gl.TEXTURE_2D, uint32(texture))
}

func (shader *BasicShader) UseMatrix4(mat4 mgl32.Mat4, input string) {
	inputID := gl.GetUniformLocation(uint32(shader.program), gl.Str(input+"\x00"))
	gl.UniformMatrix4fv(inputID, 1, false, &mat4[0])
}

func (shader *BasicShader) UseVector3(vec3 mgl32.Vec3, input string) {

}

func (shader *BasicShader) UseVector4(vec3 mgl32.Vec3, input string) {

}

func (shader *BasicShader) UseInt(vec int, input string) {
	inputID := gl.GetUniformLocation(uint32(shader.program), gl.Str(input+"\x00"))
	gl.Uniform1i(inputID, int32(vec))
}

func (shader *BasicShader) GetProgram() GLProgram {
	return shader.program
}

func compileShader(source string, shaderType uint32) (GLShader, error) {
	shader := gl.CreateShader(shaderType)

	csources, free := gl.Strs(source + "\x00")
	gl.ShaderSource(shader, 1, csources, nil)
	free()
	gl.CompileShader(shader)

	var status int32
	gl.GetShaderiv(shader, gl.COMPILE_STATUS, &status)
	if status == gl.FALSE {
		var logLength int32
		gl.GetShaderiv(shader, gl.INFO_LOG_LENGTH, &logLength)

		log := strings.Repeat("\x00", int(logLength+1))
		gl.GetShaderInfoLog(shader, logLength, nil, gl.Str(log))

		return 0, fmt.Errorf("Failed to compile.\n%v", source, log)
	}
	return GLShader(shader), nil
}
