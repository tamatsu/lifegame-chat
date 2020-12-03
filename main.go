package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/google/uuid"
	socketio "github.com/googollee/go-socket.io"
)

const sizeX = 9
const sizeY = 9
const roomName = "chatRoom"
const chatEventName = "chat"
const boardEventName = "board"

type User = uuid.UUID

func hashToInt(user User) int {
	const a = 57
	const mod = 997

	current := 57
	for _, byte := range user {
		current = (current*int(byte) + a) % mod
	}

	return current % 360
}

func jsonEncode(obj interface{}) string {
	v, err := json.Marshal(obj)
	if err != nil {
		fmt.Println(err)
	}

	return string(v)
}

func main() {
	server, _ := socketio.NewServer(nil)
	userDict := make(map[string]User)
	board := Board{
		{-1, -1, -1, -1, -1, -1, -1, -1, -1},
		{-1, -1, -1, -1, -1, -1, -1, -1, -1},
		{-1, -1, -1, -1, -1, -1, -1, -1, -1},
		{-1, -1, -1, -1, -1, -1, -1, -1, -1},
		{-1, -1, -1, -1, -1, -1, -1, -1, -1},
		{-1, -1, -1, -1, -1, -1, -1, -1, -1},
		{-1, -1, -1, -1, -1, -1, -1, -1, -1},
		{-1, -1, -1, -1, -1, -1, -1, -1, -1},
		{-1, -1, -1, -1, -1, -1, -1, -1, -1},
	}

	server.OnConnect("/", func(s socketio.Conn) error {
		server.JoinRoom("/", roomName, s)

		newID, _ := uuid.NewRandom()
		userDict[s.ID()] = newID

		s.Emit(boardEventName, jsonEncode(board))

		return nil // no error
	})

	server.OnEvent("/", chatEventName, func(s socketio.Conn, msg string) {
		user := userDict[s.ID()]

		type Message struct {
			SocketID string
			Msg      string
			Color    int
		}

		server.BroadcastToRoom("/", roomName, chatEventName, jsonEncode(Message{
			SocketID: s.ID(),
			Msg:      msg,
			Color:    hashToInt(user),
		}))
	})

	server.OnEvent("/", "tick", func(s socketio.Conn) {
		board = Transition(board)

		server.BroadcastToRoom("/", roomName, boardEventName, jsonEncode(board))
	})

	server.OnEvent("/", "toggl", func(s socketio.Conn, msg string) {
		user := userDict[s.ID()]

		type UserCommand struct {
			X int
			Y int
		}

		var cmd UserCommand
		json.Unmarshal([]byte(msg), &cmd)

		if cmd.X >= 0 && cmd.X < sizeX && cmd.Y >= 0 && cmd.Y < sizeY {
			cell := board[cmd.Y][cmd.X]
			color := hashToInt(user)

			if cell == -1 { // is Dead
				board[cmd.Y][cmd.X] = color
			} else {
				board[cmd.Y][cmd.X] = -1 // Die
			}

			server.BroadcastToRoom("/", roomName, boardEventName, jsonEncode(board))

		} else {
			log.Println("Invalid toggl command", cmd)
		}
	})

	server.OnError("/", func(s socketio.Conn, e error) {
		log.Println("error:", e)
	})

	server.OnDisconnect("/", func(s socketio.Conn, msg string) {
		log.Println("disconnect: ", msg)
	})

	go server.Serve()

	defer server.Close()

	port := os.Getenv("PORT")
	if port == "" {
		port = "8000"
	}

	http.Handle("/socket.io/", server)
	http.Handle("/", http.FileServer(http.Dir("./public")))
	log.Println("Your app is listening port:" + port)
	log.Fatal(http.ListenAndServe(":"+port, nil))

}
