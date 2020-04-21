package main

import (
	"github.com/go-gl/gl/v3.3-core/gl"
	"github.com/veandco/go-sdl2/sdl"
	"strings"
	"unsafe"
)

type Display struct {
	window  *sdl.Window
	context sdl.GLContext
	event   sdl.Event
	running bool
	frame   float64
}

func NewDisplay(width int, height int, title string) Display {
	err := sdl.Init(sdl.INIT_EVERYTHING)
	if err != nil {
		panic(err)
	}

	// OpenGL 3.3 Core
	_ = sdl.GLSetAttribute(sdl.GL_CONTEXT_MAJOR_VERSION, 3)
	_ = sdl.GLSetAttribute(sdl.GL_CONTEXT_MINOR_VERSION, 3)
	_ = sdl.GLSetAttribute(sdl.GL_CONTEXT_PROFILE_MASK, sdl.GL_CONTEXT_PROFILE_CORE)
	_ = sdl.GLSetAttribute(sdl.GL_CONTEXT_FLAGS, sdl.GL_CONTEXT_DEBUG_FLAG)

	// Create a little window at the center of the screen
	window, err := sdl.CreateWindow(title, sdl.WINDOWPOS_CENTERED, sdl.WINDOWPOS_CENTERED, int32(width), int32(height), sdl.WINDOW_OPENGL)
	if err != nil {
		panic(err)
	}

	// OpenGL initialization
	context, err := window.GLCreateContext()
	if err != nil {
		panic(err)
	}

	err = gl.Init()
	if err != nil {
		panic(err)
	}

	// In case of problem, debug to ogldebugcb
	gl.Enable(gl.DEBUG_OUTPUT_SYNCHRONOUS)
	gl.DebugMessageCallback(gl.DebugProc(ogldebugcb), gl.Ptr(nil))

	gl.Enable(gl.DEPTH_TEST)
	gl.Enable(gl.BLEND)
	gl.BlendFunc(gl.SRC_ALPHA, gl.ONE_MINUS_SRC_ALPHA)
	gl.Enable(gl.DEBUG_OUTPUT_SYNCHRONOUS)
	gl.Enable(gl.MULTISAMPLE)
	gl.Enable(gl.CULL_FACE)

/*	gl.DebugMessageInsert(
		gl.DEBUG_SOURCE_APPLICATION,
		gl.DEBUG_TYPE_ERROR,
		1, // Id
		gl.DEBUG_SEVERITY_NOTIFICATION,
		-1, // Length (negative => null-terminated)
		gl.Str("hello world\x00"))*/

	gl.ClearColor(0.0, 0.0, 0.0, 1.0)

	// vsync please
	_ = sdl.GLSetSwapInterval(1)

	// window.Hide()

	return Display{
		window:  window,
		context: context,
		event:   nil,
		running: true,
	}
}

// Running checks if the window is open
func (display *Display) Running() bool {
	return display.running
}

// Present just does window.GLSwap
func (display *Display) Present() {
	display.window.GLSwap()
}

// Update updates input stuff and running boolean
func (display *Display) Update(input *Input) {
	input.Clean()
	var event sdl.Event
	for event = sdl.PollEvent(); event != nil; event =
		sdl.PollEvent() {
		switch t := event.(type) {
		case *sdl.QuitEvent:
			_ = t
			display.running = false
		case *sdl.KeyboardEvent:
			if t.Keysym.Scancode == sdl.SCANCODE_ESCAPE {
				display.running = false
			}
			input.ProcessKeys(t)
		}
	}
	display.frame += 0.05

	// gl.ClearColor(float32(math.Cos(display.frame)), float32(math.Sin(display.frame)), 0.7, 1.0)
	gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT)
}

func ogldebugcb(source uint32, gltype uint32, id uint32, severity uint32, length int32, message string, userParam unsafe.Pointer) {
	// Program was mainly dev with a NVIDIA GPU
	// Kind of annoying false error
	if !strings.Contains(message, "GL_STATIC_DRAW") {
		println(message)
	}
}
