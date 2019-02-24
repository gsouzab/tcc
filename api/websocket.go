package main

import (
	"encoding/json"
	"log"
	"net/http"
	"reflect"
	"time"

	"github.com/gorilla/mux"
	"github.com/gorilla/websocket"
)

const (
	// Time allowed to write a message to the peer.
	writeWait = 10 * time.Second

	// Time allowed to read the next pong message from the peer.
	pongWait = 60 * time.Second

	// Send pings to peer with this period. Must be less than pongWait.
	pingPeriod = (pongWait * 9) / 10

	// Maximum message size allowed from peer.
	maxMessageSize = 2048
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  2048,
	WriteBufferSize: 2048,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

var typeRegistry = make(map[string]reflect.Type)

// Hub maintains the set of active clients and broadcasts messages to the
// clients.
type Hub struct {
	// Registered clients.
	clients map[*Client]bool

	// Inbound messages from the clients.
	broadcast chan interface{}

	// Register requests from the clients.
	register chan *Client

	// Unregister requests from clients.
	unregister chan *Client

	// Unregister requests from clients.
	messageType reflect.Type
}

// WebsocketInit initializes websocket features
func WebsocketInit(router *mux.Router) {
	telemetryHub := newHub(reflect.TypeOf(TelemetryMessage{}))
	go telemetryHub.run()

	probeHub := newHub(reflect.TypeOf(ProbeMessage{}))
	go probeHub.run()

	// Configure websocket route
	router.HandleFunc("/ws/telemetry", func(w http.ResponseWriter, r *http.Request) {
		serveWs(telemetryHub, w, r)
	})

	// Configure websocket route
	router.HandleFunc("/ws/probe", func(w http.ResponseWriter, r *http.Request) {
		serveWs(probeHub, w, r)
	})

	// Configure websocket route
	router.HandleFunc("/broadcast/telemetry", func(w http.ResponseWriter, r *http.Request) {
		broadcast(telemetryHub, w, r)
	})

	// Configure websocket route
	router.HandleFunc("/broadcast/probe", func(w http.ResponseWriter, r *http.Request) {
		broadcast(probeHub, w, r)
	})
}

func newHub(messageType reflect.Type) *Hub {
	return &Hub{
		broadcast:   make(chan interface{}),
		register:    make(chan *Client),
		unregister:  make(chan *Client),
		clients:     make(map[*Client]bool),
		messageType: messageType,
	}
}

func (h *Hub) run() {
	for {
		select {
		case client := <-h.register:
			h.clients[client] = true
		case client := <-h.unregister:
			if _, ok := h.clients[client]; ok {
				delete(h.clients, client)
				close(client.send)
			}
		case message := <-h.broadcast:
			for client := range h.clients {
				select {
				case client.send <- message:
				default:
					close(client.send)
					delete(h.clients, client)
				}
			}
		}
	}
}

// Client is a middleman between the websocket connection and the hub.
type Client struct {
	hub *Hub

	// The websocket connection.
	conn *websocket.Conn

	// Buffered channel of outbound messages.
	send chan interface{}
}

// readPump pumps messages from the websocket connection to the hub.
//
// The application runs readPump in a per-connection goroutine. The application
// ensures that there is at most one reader on a connection by executing all
// reads from this goroutine.
func (c *Client) readPump() {
	defer func() {
		c.hub.unregister <- c
		c.conn.Close()
	}()
	c.conn.SetReadLimit(maxMessageSize)
	c.conn.SetReadDeadline(time.Now().Add(pongWait))
	c.conn.SetPongHandler(func(string) error { c.conn.SetReadDeadline(time.Now().Add(pongWait)); return nil })
	for {

		message := reflect.New(c.hub.messageType)
		err := c.conn.ReadJSON(message.Interface())

		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Printf("error: %v", err)
			}
			break
		}

		c.hub.broadcast <- message
	}
}

func broadcast(hub *Hub, w http.ResponseWriter, r *http.Request) {
	messages := reflect.MakeSlice(reflect.SliceOf(hub.messageType), 0, 10)

	x := reflect.New(messages.Type())
	x.Elem().Set(messages)
	err := json.NewDecoder(r.Body).Decode(x.Interface())

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte{})
		return
	}

	n := reflect.Indirect(x).Len()
	for i := 0; i < n; i++ {
		hub.broadcast <- reflect.Indirect(x).Index(i)
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte{})
}

// writePump pumps messages from the hub to the websocket connection.
//
// A goroutine running writePump is started for each connection. The
// application ensures that there is at most one writer to a connection by
// executing all writes from this goroutine.
func (c *Client) writePump() {
	ticker := time.NewTicker(pingPeriod)
	defer func() {
		ticker.Stop()
		c.conn.Close()
	}()
	for {
		select {
		case message, ok := <-c.send:
			c.conn.SetWriteDeadline(time.Now().Add(writeWait))
			if !ok {
				// The hub closed the channel.
				c.conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}

			w, err := c.conn.NextWriter(websocket.TextMessage)

			if err != nil {
				return
			}

			if err := json.NewEncoder(w).Encode(message.(reflect.Value).Interface()); err != nil {
				return
			}

			// Add queued chat messages to the current websocket message.
			n := len(c.send)
			for i := 0; i < n; i++ {
				m := <-c.send
				if err := json.NewEncoder(w).Encode(m.(reflect.Value).Interface()); err != nil {
					return
				}
			}

			if err := w.Close(); err != nil {
				return
			}
		case <-ticker.C:
			c.conn.SetWriteDeadline(time.Now().Add(writeWait))
			if err := c.conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				return
			}
		}
	}
}

// serveWs handles websocket requests from the peer.
func serveWs(hub *Hub, w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
		return
	}
	client := &Client{hub: hub, conn: conn, send: make(chan interface{}, 256)}
	client.hub.register <- client

	// Allow collection of memory referenced by the caller by doing all work in
	// new goroutines.
	go client.writePump()
	go client.readPump()
}
