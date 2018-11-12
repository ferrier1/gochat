package main

import (
  "log"
  "net/http"

  "github.com/gorilla/websocket"
)


// connected clients
var clients = make(map[*websocket.Conn]bool)

// broadcast channel
var broadcast = make(chan Message)


var upgrader = websocket.Upgrader{}

// message object

type Message struct {
  Email string `json:"email"`
  Username string `json:"username"`
  Message string `json:"message"`
}
