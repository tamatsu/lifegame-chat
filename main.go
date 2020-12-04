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

const roomName = "chatRoom"
const chatEventName = "chat"
const boardEventName = "board"

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

	app := Init()

	server.OnConnect("/", func(s socketio.Conn) error {
		server.JoinRoom("/", roomName, s)

		newID, _ := uuid.NewRandom()
		userDict[s.ID()] = newID

		s.Emit(boardEventName, jsonEncode(app.board))

		return nil // no error
	})

	server.OnEvent("/", chatEventName, func(s socketio.Conn, msg string) {
		user := userDict[s.ID()]
		m := ToChatMsg(user, msg)

		server.BroadcastToRoom("/", roomName, chatEventName, jsonEncode(m))
	})

	server.OnEvent("/", "tick", func(s socketio.Conn) {
		var b bool
		b, app.lastTime = Clock(app.lastTime)

		if b {
			app.board = Transition(app.board)

			server.BroadcastToRoom("/", roomName, boardEventName, jsonEncode(app.board))
		}
	})

	server.OnEvent("/", "toggl", func(s socketio.Conn, msg string) {
		user := userDict[s.ID()]

		var cmd TogglCmd
		json.Unmarshal([]byte(msg), &cmd)

		app.board = Toggl(user, app.board, cmd)

		server.BroadcastToRoom("/", roomName, boardEventName, jsonEncode(app.board))
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
