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


func handleConnections(w http.ResponseWriter, r *http.Request) {
  // use upgrader on GET to make websocket
  ws, err := upgrader.Upgrade(w, r, nil)
  if err != nil {
    log.Fatal(err)
  }
  // close connection
  defer ws.Close()
  // register client
  clients[ws] = true
  // infinite loop that waits for new message on websocket
  for {
    var msg Message
    // read new message as JSON and map to Message object
    err := ws.ReadJSON(&msg)
    if err != nil {
      log.Printf("errorL %v", err)
      delete(clients, ws)
      break
    }
    // send new message to broadcast channel
    broadcast <- msg
  }

}


func handleMessages() {
  for {
    msg := <-broadcast
    // send to connected clients
    for client := range clients {
      err := client.WriteJSON(msg)
      if err != nil {
        log.Printf("error: %v", err)
        client.Close()
        delete(clients, client)
      }
    }
  }
}
