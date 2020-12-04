package main

import (
	"time"

	"github.com/google/uuid"
)

const sizeX = 9
const sizeY = 9
const clockIntervalMilliseconds = 3000

type User = uuid.UUID

type App = struct {
	board    Board
	lastTime time.Time
}

func Init() App {
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

	return App{
		board:    board,
		lastTime: lastTime,
	}
}

func hashToInt(user User) int {
	const a = 57
	const mod = 997

	current := 57
	for _, byte := range user {
		current = (current*int(byte) + a) % mod
	}

	return current % 360
}

type ChatMsg struct {
	Msg   string
	Color int
}

func ToChatMsg(user User, msg string) ChatMsg {
	return ChatMsg{
		Msg:   msg,
		Color: hashToInt(user),
	}
}

type TogglCmd struct {
	X int
	Y int
}

func Toggl(user User, board Board, cmd TogglCmd) Board {
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

func Clock(a time.Time) (bool, time.Time) {
	b := time.Now()
	if b.Sub(a).Milliseconds() >= clockIntervalMilliseconds {
		return true, b
	}

	return false, a
}
