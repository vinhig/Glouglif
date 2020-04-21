package main

type IController interface {
	Act(node INode, delta float32)
}
