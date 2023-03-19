package coordination

import (
	"sync"

	"github.com/gofiber/websocket/v2"
)

type socketMutex struct {
	socket  *websocket.Conn
	channel chan interface{}
}

type SocketMapChannelParams struct {
	SocketId string
	Socket   *websocket.Conn
}

var socketMapChannel chan SocketMapChannelParams

// Mutex must be used when accessing socketIdsToConnection for thread safety
var socketIdsToConnection = make(map[string]*socketMutex)
var socketIdsToConnectionMutex = sync.Mutex{}

/*
Initializes values for thread-safe
*/
func InitializeThreadSafeSocketHandling() {
	if socketMapChannel == nil {
		socketMapChannel = make(chan SocketMapChannelParams)
		// Handle connection map in a single go routine
		go handleConnections()
	}
}

/*
Adds connection to internal map (thread-safe)
*/
func AddConnection(socketId string, socket *websocket.Conn) {
	socketMapChannel <- SocketMapChannelParams{SocketId: socketId, Socket: socket}
}

/*
Removes connection to internal map (thread-safe)
*/
func RemoveConnection(socketId string) {
	socketMapChannel <- SocketMapChannelParams{SocketId: socketId, Socket: nil}
}

/*
Messages all current users in a separate go routines (thread-safe)
*/
func MessageEveryone(message interface{}) {
	socketIdsToConnectionMutex.Lock()
	defer socketIdsToConnectionMutex.Unlock()
	for socketId, socketInfo := range socketIdsToConnection {
		if socketInfo != nil {
			// Write message in different routine, as we don't have to wait
			go WriteMessage(socketId, message)
		}
	}
}

/*
Closes all sockets and channels.
*/
func Cleanup() {
	socketIdsToConnectionMutex.Lock()
	defer socketIdsToConnectionMutex.Unlock()
	for _, socketInfo := range socketIdsToConnection {
		if socketInfo != nil {
			socketInfo.socket.Close()
			close(socketInfo.channel)
		}
	}
}

/*
Writes message to socket (thread-safe)
*/
func WriteMessage(socketId string, message interface{}) {
	socketIdsToConnectionMutex.Lock()
	socketInfo := socketIdsToConnection[socketId]
	socketIdsToConnectionMutex.Unlock()
	if socketInfo != nil {
		socketInfo.channel <- message
	}
}

/*
Should be called as a go routine. Handles storing/deleting values from
the map of socket ids to sockets.
*/
func handleConnections() {
	// Read channel until closed
	for socketParams := range socketMapChannel {
		if socketParams.Socket != nil {
			// Create channel for the socket
			socketChannel := make(chan interface{})
			socketIdsToConnectionMutex.Lock()
			socketIdsToConnection[socketParams.SocketId] = &socketMutex{socket: socketParams.Socket, channel: socketChannel}
			socketIdsToConnectionMutex.Unlock()
			go handleMessaging(socketChannel, socketParams.Socket)
		} else {
			// Close socket messaging channel
			socketIdsToConnectionMutex.Lock()
			close(socketIdsToConnection[socketParams.SocketId].channel)
			// Delete value
			delete(socketIdsToConnection, socketParams.SocketId)
			socketIdsToConnectionMutex.Unlock()
		}
	}
}

/*
Should be called as a go routine and only used by a single socket.
Allows for writing messages to the socket in a thread-safe way.
*/
func handleMessaging(socketChannel chan interface{}, socket *websocket.Conn) {
	// Read channel until closed
	for message := range socketChannel {
		socket.WriteJSON(message)
	}
}
