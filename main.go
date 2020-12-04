package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/google/uuid"
	socketio "github.com/googollee/go-socket.io"
)

const sizeX = 9
const sizeY = 9
const roomName = "chatRoom"
const chatEventName = "chat"
const boardEventName = "board"
const clockIntervalMilliseconds = 3000

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

type ChatMsg struct {
	Msg   string
	Color int
}

func chat(user User, msg string) ChatMsg {
	return ChatMsg{
		Msg:   msg,
		Color: hashToInt(user),
	}
}

type TogglCmd struct {
	X int
	Y int
}

func toggl(user User, board Board, cmd TogglCmd) Board {
	if cmd.X >= 0 && cmd.X < sizeX && cmd.Y >= 0 && cmd.Y < sizeY {
		color := hashToInt(user)

		if board[cmd.Y][cmd.X] == -1 { // is Dead
			board[cmd.Y][cmd.X] = color
		} else {
			board[cmd.Y][cmd.X] = -1 // Die
		}

	}

	return board
}

func clock(a time.Time) (bool, time.Time) {
	b := time.Now()
	if b.Sub(a).Milliseconds() >= clockIntervalMilliseconds {
		return true, b
	}

	return false, a
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
	lastTime := time.Now()

	server.OnConnect("/", func(s socketio.Conn) error {
		server.JoinRoom("/", roomName, s)

		newID, _ := uuid.NewRandom()
		userDict[s.ID()] = newID

		s.Emit(boardEventName, jsonEncode(board))

		return nil // no error
	})

	server.OnEvent("/", chatEventName, func(s socketio.Conn, msg string) {
		user := userDict[s.ID()]
		m := chat(user, msg)

		server.BroadcastToRoom("/", roomName, chatEventName, jsonEncode(m))
	})

	server.OnEvent("/", "tick", func(s socketio.Conn) {
		var b bool
		b, lastTime = clock(lastTime)

		if b {
			board = Transition(board)

			server.BroadcastToRoom("/", roomName, boardEventName, jsonEncode(board))
		}
	})

	server.OnEvent("/", "toggl", func(s socketio.Conn, msg string) {
		user := userDict[s.ID()]

		var cmd TogglCmd
		json.Unmarshal([]byte(msg), &cmd)

		board = toggl(user, board, cmd)

		server.BroadcastToRoom("/", roomName, boardEventName, jsonEncode(board))
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
