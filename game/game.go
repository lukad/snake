package game

import (
	"code.google.com/p/go.net/websocket"
	"math/rand"
	"time"
)

var rnd *rand.Rand

type Point struct {
	X, Y int
}

type Player struct {
	Points    []Point
	Direction Direction
	Score     int
	Lifes     int
}

type Direction int

const (
	Left Direction = iota
	Right
	Up
	Down
)

type Game struct {
	Players map[*websocket.Conn]*Player
	Id      string
}

type State struct {
	Others []*Player
	You    *Player
	Time   int64
}

func New() *Game {
	game := &Game{
		Players: make(map[*websocket.Conn]*Player),
		Id:      randomString(7),
	}

	return game
}

func (g *Game) Update() {
	for _, player := range g.Players {
		switch player.Direction {
		case Left:
			player.Points[0].X -= 1
		case Right:
			player.Points[0].X += 1
		case Up:
			player.Points[0].Y -= 1
		case Down:
			player.Points[0].Y += 1
		}
		if player.Points[0].X > 15 { player.Points[0].X = 0 }
		if player.Points[0].X < 0 { player.Points[0].X = 15 }
		if player.Points[0].Y > 15 { player.Points[0].Y = 0 }
		if player.Points[0].Y < 0 { player.Points[0].Y = 15 }
	}
}

func (g *Game) Run() {
	for len(g.Players) > 0 {

		g.Update()

		for ws, player := range g.Players {
			state := State{
				You:  player,
				Time: time.Now().UnixNano(),
			}
			for otherWs, otherPlayer := range g.Players {
				if ws != otherWs {
					state.Others = append(state.Others, otherPlayer)
				}
			}
			websocket.JSON.Send(ws, &state)
		}

		time.Sleep(50 * time.Millisecond)
	}
}

func (g *Game) AddPlayer(ws *websocket.Conn) {
	g.Players[ws] = &Player{
		Points: []Point{{rnd.Intn(16), rnd.Intn(16)}},
		Direction: Left,
	}
}

func (g *Game) RemovePlayer(ws *websocket.Conn) {
	delete(g.Players, ws)
}

func (g *Game) Input(ws *websocket.Conn, direction Direction) {
	if player, ok := g.Players[ws]; ok {
		player.Direction = direction
	}
}

func randomString(length int) string {
	b := make([]byte, length)
	chars := []byte("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789")
	for i := 0; i < length; i++ {
		b[i] = chars[rnd.Intn(len(chars))]
	}
	return string(b)
}

func init() {
	rnd = rand.New(rand.NewSource(time.Now().UnixNano()))
}
