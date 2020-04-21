package main

import (
	"fmt"
	"github.com/go-gl/mathgl/mgl32"
	"github.com/icrowley/fake"
	"net"
	"strconv"
	"strings"
	"time"
)

type Client struct {
	// conn    net.Conn
	address string
	player  string
	id      int
	actions []INetworkAction
	syncs   string
}

// Connect connects to server and save connection.
func (client *Client) Connect() {
	/*	go func() {
		for {*/
	start := time.Now()
	conn, err := net.Dial("tcp", client.address)
	if err != nil {
		panic(err)
	}

	// Ask server for an id
	name := fake.FullName()
	_, _ = fmt.Fprintf(conn, "new "+name)
	id := make([]byte, 2)
	_, err = conn.Read(id)
	if err != nil {
		panic(err)
	}
	client.player = name
	client.id, _ = strconv.Atoi(string(string(id)[0]))
	fmt.Printf("We are in da place as '%s' '%s'\n", name, client.id)
	_ = conn.Close()
	fmt.Printf("ping %d", time.Now().Sub(start).Milliseconds())
	/*		}
	}()*/
}

// Start does nothing as a client
func (client *Client) Start() {

}

// Act applies modifications that occurred since last sync.
func (client *Client) Act(scene INode) {
	syncs := strings.Split(client.syncs, "\n")
	for _, sync := range syncs {
		if strings.HasPrefix(sync, "skin") {
			var id int
			position := mgl32.Vec3{}
			rotation := mgl32.Vec3{}
			_, err := fmt.Sscanf(
				sync,
				"skin %d %f %f %f %f %f %f",
				&id,
				&position[0],
				&position[1],
				&position[2],
				&rotation[0],
				&rotation[1],
				&rotation[2])

			if err != nil {
				print(err.Error())
			}

			skin := scene.GetChild(fmt.Sprintf("skin %d", id))
			if skin == nil {
				panic("kiki m'a mordu")
			}

			skin.Translate(position.X(), position.Y(), position.Z())
			skin.RotateRad(-rotation.X(), rotation.Y(), rotation.Z())
		}
	}
}

// Sync fetches distant modification and sends local modification.
func (client *Client) Sync() {
	/*go func() {*/
	start := time.Now()
	// Fetch
	conn, err := net.Dial("tcp", client.address)
	if err != nil {
		panic(err)
	}
	_, _ = fmt.Fprintf(conn, "fetch "+strconv.Itoa(client.id))
	response := make([]byte, 1024)
	_, err = conn.Read(response)
	if err != nil {
		panic(err)
	}
	body := string(response)
	if !strings.Contains(body, "<empty>") {
		client.syncs = body
	}
	_ = conn.Close()

	// Then sync
	conn, err = net.Dial("tcp", client.address)
	if err != nil {
		panic(err)
	}

	// Build sync object
	data := ""
	if len(client.actions) != 0 {
		for _, action := range client.actions {
			switch action.(type) {
			case *CameraMovementAction:
				sync := action.(*CameraMovementAction)
				data += fmt.Sprintf("skin %d %f %f %f %f %f %f",
					client.id,
					sync.position.X(),
					sync.position.Y(),
					sync.position.Z(),
					sync.rotation.X(),
					sync.rotation.Y(),
					sync.rotation.Z())
			default:
				panic("unknown INetworkAction type")
			}
		}
	} else {
		data = "<empty>"
	}

	_, _ = fmt.Fprintf(conn, fmt.Sprintf("sync %d\n", client.id)+data)
	response = make([]byte, 256)
	_, err = conn.Read(response)
	if err != nil {
		panic(err)
	}
	if !strings.Contains(string(response), "good") {
		panic("The server is a fucking liar!")
	}
	_ = conn.Close()

	// All modification send
	// Clean modifications
	client.actions = []INetworkAction{}
	/*}()*/
	fmt.Printf("ping :%d\r", time.Now().Sub(start).Milliseconds())
}

// GetPlayerName returns local Player name.
func (client *Client) GetPlayerName() string {
	return client.player
}

// SetPlayerName sets local Player name.
func (client *Client) SetPlayerName(name string) {
	client.player = name
}

// GetPlayerID returns local Player name.
func (client *Client) GetPlayerID() int {
	return client.id
}

// SetPlayerID sets local Player name.
func (client *Client) SetPlayerID(id int) {
	client.id = id
}

// Register registers a new NetworkAction to sync over the players.
func (client *Client) Register(action INetworkAction) {
	client.actions = append(client.actions, action)
}
