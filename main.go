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


func main() {
  // fileserver
  fs := http.Fileserver(http.Dir("../public"))
  http.Handle("/", fs)
  // websocket
  http.HandleFunc("/ws", handleConnections)
  // listen for incoming messages
  go HandleMessages()
  // start server on localhost:8000 with logging
  log.Println("http server started on :8000")
  err := http.ListenAndServe(":8000", nil)
  if err != nil {
    log.Fatall("ListenAndServe: ", err)
  }
}
