package main

import (
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/gorilla/websocket"
)

var producers = make(map[*websocket.Conn]bool) // connected producer clients
var consumers = make(map[*websocket.Conn]bool) // connected consumer clients
var broadcast = make(chan TelemetryMessage)    // broadcast channel

// Configure the upgrader
var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true	
	},
}

//WebsocketInit ...
func WebsocketInit(router *mux.Router) {
	// Configure websocket route
	router.HandleFunc("/producer", handleProducerConnections)
	router.HandleFunc("/consumer", handleConsumerConnections)

	// Start listening for incoming chat messages
	go handleMessages()
}

func handleProducerConnections(w http.ResponseWriter, r *http.Request) {
	// Upgrade initial GET request to a websocket
	ws, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Fatal(err)
	}
	// Make sure we close the connection when the function returns
	defer ws.Close()

	// Register our new client
	producers[ws] = true

	for {
		var msg TelemetryMessage
		// Read in a new message as JSON and map it to a Message object
		err := ws.ReadJSON(&msg)
		if err != nil {
			// log.Printf("error: %v", err)
			delete(producers, ws)
			break
		}
		// Send the newly received message to the broadcast channel
		broadcast <- msg
	}
}

func handleConsumerConnections(w http.ResponseWriter, r *http.Request) {
	// Upgrade initial GET request to a websocket
	ws, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Fatal(err)
	}

	// Register our new client
	consumers[ws] = true
}

func handleMessages() {
	for {
		// Grab the next message from the broadcast channel
		msg := <-broadcast
		// Send it out to every client that is currently connected
		for client := range consumers {
			err := client.WriteJSON(msg)
			if err != nil {
				// log.Printf("error: %v", err)
				client.Close()
				delete(consumers, client)
			}
			log.Printf("message: %v", msg)
		}
	}
}
