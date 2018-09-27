package main

import (
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/gorilla/websocket"
)

var clients = make(map[*websocket.Conn]bool) // connected clients
var broadcast = make(chan TelemetryMessage)  // broadcast channel

// Configure the upgrader
var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

// TelemetryMessage object
type TelemetryMessage struct {
	Sensor    string `json:"sensor"`
	CO2       int    `json:"co2"`
	Temp      int    `json:"temp"`
	Hum       int    `json:"hum"`
	Latitude  string `json:"latitude"`
	Longitude string `json:"longitude"`
	CreatedAt int64  `json:"createdAt"`
}

func WebsocketInit(router *mux.Router) {
	// Configure websocket route
	router.HandleFunc("/ws", handleConnections)

	// Start listening for incoming chat messages
	go handleMessages()
}

func handleConnections(w http.ResponseWriter, r *http.Request) {
	// Upgrade initial GET request to a websocket
	ws, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Fatal(err)
	}
	// Make sure we close the connection when the function returns
	defer ws.Close()

	// Register our new client
	clients[ws] = true

	for {
		var msg TelemetryMessage
		// Read in a new message as JSON and map it to a Message object
		err := ws.ReadJSON(&msg)
		if err != nil {
			// log.Printf("error: %v", err)
			delete(clients, ws)
			break
		}
		// Send the newly received message to the broadcast channel
		broadcast <- msg
	}
}

func handleMessages() {
	for {
		// Grab the next message from the broadcast channel
		msg := <-broadcast
		// Send it out to every client that is currently connected
		for client := range clients {
			err := client.WriteJSON(msg)
			if err != nil {
				// log.Printf("error: %v", err)
				client.Close()
				delete(clients, client)
			}

		}
	}
}
