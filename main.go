package main

import (
	"fmt"
	"log"
	"net/http"
  "encoding/json"
  "os"

	"github.com/googollee/go-socket.io"
  "github.com/google/uuid"
)


type User uuid.UUID

type Message struct {
  SocketID string
  Msg string
  Color int
}

const chatroom = "chatroom" // Room name

func hash(user uuid.UUID) int {
  current := 57

  for _, b := range user {
    n := int(b) + 1
    current = (current * n) % 575757 + 57
  }

  return current % 360
}

func _log(v interface{}) interface{} {
  fmt.Println(v)
  return v
}

func emitBoard(board Board, s socketio.Conn) {
  v, err := json.Marshal(board)
  if err != nil {
    fmt.Println(err)
  }

  s.Emit("board", string(v))
}

func broadCastBoard(board Board, server* socketio.Server) {
  v, err := json.Marshal(board)
  if err != nil {
    fmt.Println(err)
  }
  // _log(string(v))
  server.BroadcastToRoom("/", chatroom, "board", string(v))
}

func encode(obj interface{}) string {
  v, err := json.Marshal(obj)
  if err != nil {
    fmt.Println(err)
  }

  return string(v)
}

func main() {

  userDict := make(map[string]uuid.UUID)

	server, _ := socketio.NewServer(nil)

  board := Board{
    { -1, -1, -1, -1, -1, -1, -1, -1, -1 },
    { -1, -1, -1, -1, -1, -1, -1, -1, -1 },
    { -1, -1, -1, -1, -1, -1, -1, -1, -1 },
    { -1, -1, -1, -1, -1, -1, -1, -1, -1 },
    { -1, -1, -1, -1, -1, -1, -1, -1, -1 },
    { -1, -1, -1, -1, -1, -1, -1, -1, -1 },
    { -1, -1, -1, -1, -1, -1, -1, -1, -1 },
    { -1, -1, -1, -1, -1, -1, -1, -1, -1 },
    { -1, -1, -1, -1, -1, -1, -1, -1, -1 },
  }
	
	server.OnConnect("/", func(s socketio.Conn) error {
		fmt.Println("connected:", s.ID())

    u, _ := uuid.NewRandom()
    userDict[s.ID()] = u
    fmt.Println("Users:", userDict)

    server.JoinRoom("/", chatroom, s)

    emitBoard(board, s)

		return nil // no error
	})

  server.OnEvent("/", "chat", func(s socketio.Conn, msg string) {
    fmt.Println("chat:", msg)

    fmt.Println("rooms", server.Rooms("/"))
    fmt.Println("roomlen", server.RoomLen("/", chatroom))

    user := userDict[s.ID()]
    _log(user)

    m := Message{
      SocketID: s.ID(),
      Msg: msg,
      Color: hash(user),
    }

    _log(m)
    server.BroadcastToRoom("/", chatroom, "chat", encode(m))
  })

  server.OnEvent("/", "tick", func(s socketio.Conn) {
    fmt.Println("tick: ")

    board = Transition(board)

    // for _, l := range board {
    //   fmt.Println(l)
    // }

    broadCastBoard(board, server)
  })

  server.OnEvent("/", "toggl", func(s socketio.Conn, msg string) {
    fmt.Println("toggl:")

    user := userDict[s.ID()]

    type Command struct {
      X int
      Y int
    }

    var cmd Command
    _log(msg)
    json.Unmarshal([]byte(msg), &cmd)
    fmt.Println(cmd)

    cell := board[cmd.Y][cmd.X]
    if (cell == -1) { // Dead
      board[cmd.Y][cmd.X] = hash(user)
    } else {
      board[cmd.Y][cmd.X] = -1
    }

    broadCastBoard(board, server)

  })

	server.OnError("/", func(s socketio.Conn, e error) {
		fmt.Println("error:", e)
	})

	server.OnDisconnect("/", func(s socketio.Conn, reason string) {
		fmt.Println("disconnect: ", reason)
	})

	go server.Serve()



	defer server.Close()

  port := os.Getenv("PORT")
  if (port == "") {
    port = "8000"
  }

	http.Handle("/socket.io/", server)
	http.Handle("/", http.FileServer(http.Dir("./public")))
	log.Println("Serving at localhost:" + port + "...")
	log.Fatal(http.ListenAndServe(":" + port , nil))



}

