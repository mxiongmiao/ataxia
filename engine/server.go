package engine

/*
   Server structures and functions
*/

import (
	"io"
	"log"
	"net"
	//	"bufio"
	"sync"
	"time"
	//	"container/list"
	golua "github.com/aarzilli/golua/lua"
	"github.com/xenith-studios/ataxia/game"
	"github.com/xenith-studios/ataxia/handler"
	"github.com/xenith-studios/ataxia/lua"
)

type PlayerList struct {
	players map[string]*Account
	mu      sync.RWMutex
}

type Server struct {
	socket     *net.TCPListener
	luaState   *golua.State
	World      *game.World
	PlayerList *PlayerList
	In         chan string
	shutdown   chan bool
}

func NewPlayerList() (list *PlayerList) {
	return &PlayerList{players: make(map[string]*Account)}
}

func (list *PlayerList) Add(name string, player *Account) {
	list.mu.Lock()
	defer list.mu.Unlock()
	list.players[name] = player
}

func (list *PlayerList) Delete(name string) {
	list.mu.Lock()
	defer list.mu.Unlock()
	list.players[name] = nil
}

func (list *PlayerList) Get(name string) (player *Account) {
	list.mu.RLock()
	defer list.mu.RUnlock()
	player = list.players[name]
	return
}

// NewServer creates a new server and returns a pointer to it
func NewServer(port int, shutdown chan bool) *Server {
	listener, err := net.ListenTCP("tcp", &net.TCPAddr{IP: net.ParseIP(""), Port: port, Zone: ""})
	if err != nil {
		log.Fatalln("Failed to create server:", err)
		return nil
	}

	server := &Server{
		luaState:   lua.MainState,
		World:      game.NewWorld(lua.MainState),
		PlayerList: NewPlayerList(),
		socket:     listener,
		shutdown:   shutdown,
		In:         make(chan string, 1024),
	}

	server.PublishAccessors(server.luaState)
	server.World.PublishAccessors(server.luaState)

	return server
}

func (server *Server) InitializeWorld() {
	server.World.LoadAreas()
	server.World.Initialize()

}

func (server *Server) Shutdown() {
	if server.socket != nil {
		server.SendToPlayers("Server is shutting down!")
		for _, player := range server.PlayerList.players {
			if player != nil {
				player.Close()
			}
		}
		server.socket.Close()
	}
}

func (server *Server) Listen() {
	for {
		if server.socket == nil {
			log.Println("Server socket closed")
			server.shutdown <- true
			return
		}
		conn, err := server.socket.Accept()
		if err != nil {
			log.Println("Failed to accept new connection:", err)
			continue
		} else {
			c := new(connection)
			c.remoteAddr = conn.RemoteAddr().String()
			c.socket = conn
			c.handler = handler.NewTelnetHandler(conn)
			log.Println("Accepted a new connection:", c.remoteAddr)
			player := NewAccount(server, c)
			go player.Run()
		}
	}
}

func (server *Server) Run() {
	// Main loop
	// Handle network messages (push user events)
	// Handle game updates
	// Game tick
	// Time update
	// Weather update
	// Entity updates (push events)
	// Handle pending events
	// Handle pending messages (network and player)
	// Sleep
	for {
		// Sleep for 1 ms
		time.Sleep(1000000)
	}
}

func (server *Server) AddPlayer(player *Account) {
	server.PlayerList.Add(player.Name, player)
}

func (server *Server) RemovePlayer(player *Account) {
	server.PlayerList.Delete(player.Name)
}

func (server *Server) Write(buf []byte) (n int, err error) {
	return 0, io.EOF
}
