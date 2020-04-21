package main

import (
	"fmt"
	"github.com/icrowley/fake"
	"net"
	"strconv"
	"strings"
)

type Server struct {
	ln         net.Listener
	port       string
	player     string   // Local player name
	id         int      // Local player id
	players    []string // Name of distant player
	syncs      []string // All sync since start
	syncsDelay [][]int  // Syncs not sent for each player
	started    bool     // Is the game started ?
}

// Connect launches the server according to given port.
func (server *Server) Connect() {
	// Listen to incoming request
	// But doesn't handle them
	var err error
	server.ln, err = net.Listen("tcp", server.port)
	if err != nil {
		panic(err)
	}
}

// Start starts the server on another thread.
func (server *Server) Start() {
	// Here we create the local player
	server.syncsDelay = append(server.syncsDelay, []int{})
	server.player = fake.FullName()
	server.players = append(server.players, server.player)
	server.id = 0
	server.syncs = append(server.syncs, fmt.Sprintf("new %s %d", server.player, server.id))

	// For ever loop on another thread
	go func() {
		for {
			// Handle incoming request
			conn, err := server.ln.Accept()
			if err != nil {
				panic(err)
			}
			data := make([]byte, 256)
			_, err = conn.Read(data)
			if err != nil {
				print(err.Error())
			}
			body := string(data)

			if strings.Contains(body, "new") {
				// New player can only join if the game isn't started
				if !server.started {
					// Clean given command
					// Something from "new Sarah Montgomery\x00" to "Sarah Montgomery"
					name := strings.Replace(body, "new ", "", 1)
					name = strings.Replace(name, "\x00", "", 1)

					// Send id for new player
					_, _ = fmt.Fprintf(conn, strconv.Itoa(len(server.players)))

					sync := fmt.Sprintf("new %s %d", name, len(server.players))
					println(sync)

					// Register this new player
					server.players = append(server.players, name) // new name
					firstSyncs := make([]int, len(server.syncs))  // all modification before creation has to be sync
					for i, _ := range server.syncs {
						firstSyncs[i] = i
					}
					server.syncsDelay = append(server.syncsDelay, firstSyncs) // new array of not sent sync
					server.NewSync(sync, len(server.players)-1)               // other should be warned about this creation
				} else {
					_, _ = fmt.Fprintf(conn, "kick them all")
				}
			} else if strings.Contains(body, "fetch") {
				// Get player id first
				id := getPlayerID(body, 6)

				// Return syncs that hasn't be sent yet
				data := ""
				for _, syncID := range server.syncsDelay[id] {
					data += server.syncs[syncID] + "\n"
				}
				if data == "" {
					_, _ = fmt.Fprintf(conn, "<empty>")
				} else {
					_, _ = fmt.Fprintf(conn, data)

					// Player was warned
					// So he won't be warned again
					server.syncsDelay[id] = []int{}
				}
			} else if strings.Contains(body, "sync") {
				// Don't bother me with modification if you can't do modification
				// No useless job
				if !strings.Contains(body, "<empty>") {
					// Get player id first
					id := getPlayerID(body, 5)

					// Remove "sync %d" from received sync
					body = body[7 : len(body)-1]

					// Register sync
					server.NewSync(body, id)
					// println(id)
				}
				_, _ = fmt.Fprintf(conn, "good")
			} else {
				_, _ = fmt.Fprintf(conn, "what?")
			}
			_ = conn.Close()
		}
	}()
}

// Act applies modifications that occurred since last sync.
func (server *Server) Act(scene INode) {

}

// GetPlayerName returns local Player name.
func (server *Server) GetPlayerName() string {
	return server.player
}

// SetPlayerName sets local Player name.
func (server *Server) SetPlayerName(name string) {
	server.player = name
}

// GetPlayerID returns local Player name.
func (server *Server) GetPlayerID() int {
	return server.id
}

// SetPlayerID sets local Player name.
func (server *Server) SetPlayerID(id int) {
	server.id = id
}

// Register registers a new NetworkAction to sync over the players.
func (server *Server) Register(action INetworkAction) {

}

// Register registers a new NetworkAction to sync over the players.
func (server *Server) Sync() {

}

func (server *Server) NewSync(sync string, author int) {
	syncID := len(server.syncs)
	server.syncs = append(server.syncs, sync)

	// All player that aren't author should be warned
	for i, _ := range server.syncsDelay {
		if i != author {
			server.syncsDelay[i] = append(server.syncsDelay[i], syncID)
		}
	}
}
