package main

import (
	"fmt"
	"github.com/go-gl/mathgl/mgl32"
	"os"
	"strconv"
)

type INetworkAction interface {
	Author() string
}

type NewPlayerAction struct {
	player string
}

func (action *NewPlayerAction) Author() string {
	return action.player
}

type CameraMovementAction struct {
	player   string
	playerID int
	position mgl32.Vec3
	rotation mgl32.Vec3
}

func (action *CameraMovementAction) Author() string {
	return action.player
}

type INetwork interface {
	Connect()
	Start()
	Act(scene INode)
	Sync()
	GetPlayerName() string
	SetPlayerName(name string)
	GetPlayerID() int
	SetPlayerID(id int)
	Register(action INetworkAction)
}

// NewNetwork initializes a INetwork object according to os.Args.
func NewNetwork() INetwork {
	if len(os.Args) != 3 {
		panic("Not enough argument.")
	}
	mode := os.Args[1]
	address := os.Args[2]

	if mode == "CLIENT" {
		fmt.Printf("Run as client connected to '%s'\n", address)

		network := new(Client)
		network.address = address
		return network
	} else if mode == "SERVER" {
		port, err := strconv.Atoi(address)
		if err != nil {
			println("While converting port to integer:")
			panic(err)
		}
		fmt.Printf("Run as server listening on localhost:%d\n", port)

		network := new(Server)
		network.port = ":" + address
		return network
	}
	panic(fmt.Sprintf("'%s' isn't a know game mode\n", mode))
}

func getPlayerID(command string, place int) int {
	id, err := strconv.Atoi(command[place:place+1])
	if err != nil {
		panic(err)
	}
	return id
}
