package main

import "github.com/go-gl/mathgl/mgl32"

type IAnimation interface {
	Act(node INode, delta float32)
	Play()
	Pause()
}

type AnimationStep struct {
	position mgl32.Vec3
	rotation mgl32.Vec3
	scale    mgl32.Vec3
}

type Transformanimation struct {
	steps       []AnimationStep
	currentStep int
	intervals   []float32
	elapsed     float32

	running bool
}

func (anim *Transformanimation) Act(node INode, delta float32) {
	anim.elapsed += delta

	// Find the currentStep
	var total float32
	total = 0.0
	found := false
	for i, interval := range anim.intervals {
		if anim.elapsed < total+interval {
			anim.currentStep = i
			found = true
			break
		}
		total += interval
	}

	if !found {
		anim.currentStep = 0
		anim.elapsed = 0.0
	}

	step := anim.steps[anim.currentStep]
	interval := anim.intervals[anim.currentStep]

	var next AnimationStep
	if len(anim.steps)-1 == anim.currentStep {
		next = anim.steps[0]
	} else {
		next = anim.steps[anim.currentStep+1]
	}

	diffPosition := next.position.Sub(step.position).Mul((1/interval) * delta)
	diffRotation := next.rotation.Sub(step.rotation).Mul((1/interval) * delta)
	diffScale := next.scale.Sub(step.scale).Mul((1/interval) * delta)

	newPosition := node.(*Node).position.Add(diffPosition)
	newRotation := node.GetLocalRotation().Add(diffRotation)
	newScale := node.(*Node).scale.Add(diffScale)

	node.Translate(newPosition.X(), newPosition.Y(), newPosition.Z())
	node.Rotate(newRotation.X(), newRotation.Y(), newRotation.Z())
	node.Scale(newScale.X(), newScale.Y(), newScale.Z())
}

func (anim *Transformanimation) Play() {
	anim.running = true
}

func (anim *Transformanimation) Pause() {
	anim.running = false
}
