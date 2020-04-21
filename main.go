package main

import (
	"github.com/go-gl/mathgl/mgl32"
	"github.com/veandco/go-sdl2/sdl"
)

const WIDTH = 600
const HEIGHT = 400

var input *Input
var network INetwork

func main() {
	network = NewNetwork()
	network.Connect()

	network.Start()

	display := NewDisplay(WIDTH, HEIGHT, "Glouglif")

	input = NewInput()

	sdl.SetRelativeMouseMode(true)

	_ = NewTriangle()

	light := NewDirectionalLight()
	light.Translate(0.0, 0.0, 5.0)
	light.Scale(1.0, 1.0, 1.0)
	light.Rotate(90.0, 1.0, 1.0)
	_ = light

	camera := NewPerspectiveCamera()
	camera.SetName("Camera")
	camera.Translate(0.0, 0.0, 0.0)
	camera.Scale(1.0, 1.0, 1.0)
	camera.Rotate(0.0, 0.0, 0.0)
	camera.SetController(&FirstPersonController{})

	scene := NewNode()
	scene.SetName("MainScene")
	scene.Translate(0.0, 0.0, 0.0)
	scene.Scale(1.0, 1.0, 1.0)
	scene.Rotate(0.0, 0.0, 0.0)

	coucou := NewMesh("cube", "square.png")
	coucou.SetName("Glouglou")
	coucou.Translate(1.0, 0.0, 0.0)
	coucou.Scale(.2, .2, .2)
	coucou.Rotate(0.0, 0.0, 0.0)
	coucou.SetAnimation(&Transformanimation{
		steps: []AnimationStep{
			{
				position: mgl32.Vec3{-1.0, 0.0, 0.0},
				scale:    mgl32.Vec3{.2, .2, .2},
				rotation: mgl32.Vec3{0.0, 0.0, 0.0},
			},
			{
				position: mgl32.Vec3{3.0, 0.0, 0.0},
				scale:    mgl32.Vec3{1.0, 1.0, 1.0},
				rotation: mgl32.Vec3{180.0, 0.0, 0.0},
			},
			// {position: mgl32.Vec3{1.0,0.0,0.0}},
		},
		intervals: []float32{
			1.0, 1.0,
		},
	})

	player1 := NewMesh("monkey.obj", "square.png")
	player1.SetName("skin 0")
	player1.Translate(0.0, 0.0, 0.0)
	player1.Scale(.2, .2, .2)
	player1.Rotate(0.0, 0.0, 0.0)

	player2 := NewMesh("monkey.obj", "belgium.png")
	player2.SetName("skin 1")
	player2.Translate(1.0, 1.0, 0.0)
	player2.Scale(.2, .2, .2)
	player2.Rotate(0.0, 0.0, 0.0)

	player3 := NewMesh("monkey.obj", "hexagon.png")
	player3.SetName("skin 2")
	player3.Translate(2.0, 2.0, 0.0)
	player3.Scale(.2, .2, .2)
	player3.Rotate(0.0, 0.0, 0.0)

	tritri := NewMesh("hexagon.obj", "hexagon.png")
	tritri.SetName("Triangularity")
	tritri.Translate(0.0, -1, 5)
	tritri.Scale(1, 1, 1)
	tritri.Rotate(180.0, 0.0, 0.0)
	tritri.shadowReceiver = true
	tritri.SetModifier(&HexaModifier{sizeX: 10, sizeY: 10, width: 1.727 / 2, height: 1.495, position: mgl32.Vec3{0.0, -1, 5}})

	scene.AddChild(tritri)
	scene.AddChild(coucou)
	scene.AddChild(camera)
	scene.AddChild(player1)
	scene.AddChild(player2)
	scene.AddChild(player3)

	diffusePipeline := NewDiffusePipeline()
	shadowPipeline := NewShadowPipeline()

	// tritri.texture = shadowPipeline.texture

	/*	var frame float32
		frame = 0.0*/

	/*var start time.Time
	start = time.Now()*/

	// coucou.GetAnimation().Play()

	for display.Running() {
		// Update SDL Window and input
		// delta := float32(time.Now().Sub(start).Seconds())
		network.Act(scene)
		display.Update(input)
		scene.Update(0.0016)

		shadowPipeline.Begin()
		shadowPipeline.LinkCamera(light)

		scene.Render(shadowPipeline)

		shadowPipeline.End()

		// Link diffuse stuff
		// for color rendering
		diffusePipeline.Begin()

		// Main camera
		diffusePipeline.LinkCamera(camera)

		// Shadow map
		diffusePipeline.LinkTexture(shadowPipeline.texture, 1)

		// Light stuff
		diffusePipeline.GetShader().UseMatrix4(light.view, "i_light_view")
		diffusePipeline.GetShader().UseMatrix4(light.projection, "i_light_projection")
		diffusePipeline.GetShader().UseMatrix4(light.GetWorldModel(), "i_light_model")
		diffusePipeline.GetShader().UseMatrix4(bias, "i_bias")

		scene.Render(diffusePipeline)
		// Show to screen what was rendered
		diffusePipeline.End()

		// camera.Rotate(0.0, float32(frame), 0.0)

		/*		if input.IsKeyDown(E) {
				coucou.Rotate(frame, frame, frame)
				frame += 50 * delta
			}*/

		/*start = time.Now()*/
		// fmt.Printf("\rFrame computed in: %f", delta)

		// Network sync
		network.Sync()

		// GL swap buffer
		display.Present()
	}
}
