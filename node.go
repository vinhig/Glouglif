package main

// DISCLAIMER
// This part is HEAVILY inspired by https://github.com/g3n/engine/blob/master/core/node.go
// END OF DISCLAIMER

import (
	"github.com/go-gl/mathgl/mgl32"
)

type INode interface {
	Parent() INode
	SetParent(parent INode)
	Children() []INode
	AddChild(node INode)
	RemoveChild(node INode)
	Visible() bool
	SetVisible(visible bool)
	Name() string
	SetName(name string)
	GetNode() *Node
	GetINode() INode
	Render(pipeline IPipeline)
	Update(delta float32)
	Translate(x float32, y float32, z float32)
	Rotate(x float32, y float32, z float32)
	RotateRad(x float32, y float32, z float32)
	Scale(x float32, y float32, z float32)
	GetLocalPosition() mgl32.Vec3
	GetWorldPosition() mgl32.Vec3
	GetLocalRotation() mgl32.Vec3
	GetWorldRotation() mgl32.Vec3
	GetLocalScale() mgl32.Vec3
	GetWorldScale() mgl32.Vec3
	GetLocalModel() mgl32.Mat4
	GetWorldModel() mgl32.Mat4
	GetController() IController
	SetController(controller IController)
	GetAnimation() IAnimation
	SetAnimation(anim IAnimation)
	GetChild(name string) INode
}

type Node struct {
	parent          INode
	inode           INode
	visible         bool
	name            string
	position        mgl32.Vec3
	rotation        mgl32.Vec3
	scale           mgl32.Vec3
	modelNeedUpdate bool
	localModel      mgl32.Mat4
	worldModel      mgl32.Mat4
	children        []INode
	controller      IController
	animation       IAnimation
}

func NewNode() *Node {
	node := new(Node)
	node.scale = mgl32.Vec3{1.0, 1.0, 1.0}
	node.visible = true
	node.Init()

	return node
}

func (node *Node) Init() {
	// Position, rotation and scale
	node.scale = mgl32.Vec3{1.0, 1.0, 1.0}

	node.inode = node
	node.visible = true
	node.children = make([]INode, 0)
}

func (node *Node) Parent() INode {
	return node.parent
}

func (node *Node) SetParent(parent INode) {
	node.parent = parent
	node.modelNeedUpdate = true
}

func (node *Node) Children() []INode {
	return node.children
}

func (node *Node) AddChild(child INode) {
	if child.Parent() != nil {
		child.Parent().RemoveChild(child)
	}
	child.SetParent(node.GetINode())
	node.children = append(node.children, child.GetINode())
}

func (node *Node) RemoveChild(child INode) {
	for pos, current := range node.children {
		if current == child {
			copy(node.children[pos:], node.children[pos+1:])
			node.children[len(node.children)-1] = nil
			node.children = node.children[:len(node.children)-1]
			child.GetNode().parent = nil
		}
	}
}

func (node *Node) Visible() bool {
	return node.visible
}

func (node *Node) SetVisible(visible bool) {
	node.visible = visible
}

func (node *Node) Name() string {
	return node.name
}

func (node *Node) SetName(name string) {
	node.name = name
}

func (node *Node) GetNode() *Node {
	return node
}

func (node *Node) GetINode() INode {
	return node.inode
}

func (node *Node) Render(pipeline IPipeline) {
	if node.visible {
		// fmt.Printf("Rendering from %s\n", node.Name())
		for i := 0; i < len(node.children); i++ {
			node.children[i].Render(pipeline)
		}
	}
}

func (node *Node) Update(delta float32) {
	if node.visible {
		if node.controller != nil {
			node.controller.Act(node, delta)
		}
		if node.animation != nil {
			node.animation.Act(node, delta)
		}
		// fmt.Printf("Updating from %s\n", node.Name())
		for i := 0; i < len(node.children); i++ {
			node.children[i].Update(delta)
		}
	}
}

func (node *Node) Translate(x float32, y float32, z float32) {
	node.position = mgl32.Vec3{x, y, z}
	node.modelNeedUpdate = true
}

func (node *Node) Rotate(x float32, y float32, z float32) {
	node.rotation = mgl32.Vec3{mgl32.DegToRad(x), mgl32.DegToRad(y), mgl32.DegToRad(z)}
	node.modelNeedUpdate = true
}

func (node *Node) RotateRad(x float32, y float32, z float32) {
	node.rotation = mgl32.Vec3{x, y, z}
	node.modelNeedUpdate = true
}

func (node *Node) Scale(x float32, y float32, z float32) {
	node.scale = mgl32.Vec3{x, y, z}
	node.modelNeedUpdate = true
}

func (node *Node) GetLocalPosition() mgl32.Vec3 {
	return node.position
}

func (node *Node) GetWorldPosition() mgl32.Vec3 {
	panic("nope")
}

func (node *Node) GetLocalRotation() mgl32.Vec3 {
	rotation := mgl32.Vec3{
		mgl32.RadToDeg(node.rotation.X()),
		mgl32.RadToDeg(node.rotation.Y()),
		mgl32.RadToDeg(node.rotation.Z()),
	}
	return rotation
}

func (node *Node) GetWorldRotation() mgl32.Vec3 {
	panic("nope")
}

func (node *Node) GetLocalScale() mgl32.Vec3 {
	return node.scale
}

func (node *Node) GetWorldScale() mgl32.Vec3 {
	panic("nope")
}

func (node *Node) GetLocalModel() mgl32.Mat4 {
	if node.modelNeedUpdate {
		// fmt.Printf("Re-computing transform matrix from '%s'\n", node.name)
		var posMat, scaleMat mgl32.Mat4
		var rotMatX, rotMatY, rotMatZ, rotMat mgl32.Mat3
		posMat = mgl32.Translate3D(node.position.X(), node.position.Y(), node.position.Z())
		scaleMat = mgl32.Scale3D(node.scale.X(), node.scale.Y(), node.scale.Z())
		rotMatX = mgl32.Rotate3DX(node.rotation.X())
		rotMatY = mgl32.Rotate3DY(node.rotation.Y())
		rotMatZ = mgl32.Rotate3DZ(node.rotation.Z())
		rotMat = rotMatX.Mul3(rotMatY).Mul3(rotMatZ)

		node.localModel = posMat.Mul4(rotMat.Mat4()).Mul4(scaleMat)
		node.modelNeedUpdate = false
	}

	return node.localModel

}

func (node *Node) GetWorldModel() mgl32.Mat4 {

	if node.parent != nil {
		// return node.parent.GetWorldModel().Mul4(node.GetLocalModel())
		return node.GetLocalModel().Mul4(node.parent.GetWorldModel())
	}

	return node.GetLocalModel()
}

func (node *Node) GetController() IController {
	return node.controller
}

func (node *Node) SetController(controller IController) {
	node.controller = controller
}

func (node *Node) GetAnimation() IAnimation {
	return node.animation
}

func (node *Node) SetAnimation(anim IAnimation) {
	node.animation = anim
}

func (node *Node) GetChild(name string) INode {
	if node.name == name {
		return node
	}

	for _, child := range node.children {
		d := child.GetChild(name)
		if d != nil {
			return d
		}
	}

	return nil
}
