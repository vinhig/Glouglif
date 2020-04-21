package main

import "github.com/veandco/go-sdl2/sdl"

// KeyCode is a keyboard key represented by an integer
type KeyCode int

// KeyCode is a keyboard key represented by an integer
const (
	A KeyCode = iota
	Z
	E
	R
	T
	Y
	Q
	S
	D
	W
	UP
	DOWN
	LEFT
	RIGHT
)

type Input struct {
	keys []bool
	last []bool
}

func NewInput() *Input {
	input := new(Input)
	input.keys = make([]bool, 14)
	input.last = make([]bool, 14)

	return input
}

// IsKeyPressed indicates if the keys has been just down
func (input *Input) IsKeyPressed(key KeyCode) bool {
	return input.keys[key] && !input.last[key]
}

// IsKeyDown indicates if the keys is currently down
func (input *Input) IsKeyDown(key KeyCode) bool {
	return input.keys[key]
}

// IsKeyDown indicates if the keys is currently up
func (input *Input) IsKeyUp(key KeyCode) bool {
	return !input.keys[key]
}


func (input *Input) Clean() {
	for key, state := range input.keys {
		input.last[key] = state
		// input.keys[key] = false
	}
}

func (input *Input) ProcessKeys(event *sdl.KeyboardEvent) {
	if event.State == sdl.PRESSED {
		switch event.Keysym.Scancode {
		case sdl.SCANCODE_Q:
			input.keys[A] = true
			break
		case sdl.SCANCODE_W:
			input.keys[Z] = true
			break
		case sdl.SCANCODE_E:
			input.keys[E] = true
			break
		case sdl.SCANCODE_R:
			input.keys[R] = true
			break
		case sdl.SCANCODE_T:
			input.keys[T] = true
			break
		case sdl.SCANCODE_Y:
			input.keys[Y] = true
			break
		case sdl.SCANCODE_A:
			input.keys[Q] = true
			break
		case sdl.SCANCODE_S:
			input.keys[S] = true
			break
		case sdl.SCANCODE_D:
			input.keys[D] = true
			break
		case sdl.SCANCODE_Z:
			input.keys[W] = true
			break
		case sdl.SCANCODE_UP:
			input.keys[UP] = true
			break
		case sdl.SCANCODE_DOWN:
			input.keys[DOWN] = true
			break
		case sdl.SCANCODE_LEFT:
			input.keys[LEFT] = true
			break
		case sdl.SCANCODE_RIGHT:
			input.keys[RIGHT] = true
			break
		default:
			println("KeyCode event not implemented.")
		}
	} else if event.State == sdl.RELEASED {
		switch event.Keysym.Scancode {
		case sdl.SCANCODE_Q:
			input.keys[A] = false
			break
		case sdl.SCANCODE_W:
			input.keys[Z] = false
			break
		case sdl.SCANCODE_E:
			input.keys[E] = false
			break
		case sdl.SCANCODE_R:
			input.keys[R] = false
			break
		case sdl.SCANCODE_T:
			input.keys[T] = false
			break
		case sdl.SCANCODE_Y:
			input.keys[Y] = false
			break
		case sdl.SCANCODE_A:
			input.keys[Q] = false
			break
		case sdl.SCANCODE_S:
			input.keys[S] = false
			break
		case sdl.SCANCODE_D:
			input.keys[D] = false
			break
		case sdl.SCANCODE_Z:
			input.keys[W] = false
			break
		case sdl.SCANCODE_UP:
			input.keys[UP] = false
			break
		case sdl.SCANCODE_DOWN:
			input.keys[DOWN] = false
			break
		case sdl.SCANCODE_LEFT:
			input.keys[LEFT] = false
			break
		case sdl.SCANCODE_RIGHT:
			input.keys[RIGHT] = false
			break
		default:
			println("KeyCode event not implemented.")
		}
	}
}
