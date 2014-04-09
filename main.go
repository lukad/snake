package main

import (
	"code.google.com/p/go.net/websocket"
	"github.com/go-martini/martini"
	"github.com/lukad/snake/game"
	"io"
	"log"
	"net/http"
	"os"
)

var games = make(map[string]*game.Game)

func connect(ws *websocket.Conn) {
	// Message receive from client should either be 'new' or a game id
	log.Println("New websocket connection:", ws.Request().RemoteAddr)
	var msg string
	websocket.Message.Receive(ws, &msg)

	var g *game.Game
	var ok bool

	// Check if the player creates a new game or is connecting to an exiting game
	if msg == "new" {
		g = game.New()
		games[g.Id] = g
		log.Printf("Player %s created game %s", ws.Request().RemoteAddr, g.Id)
	} else {
		log.Printf("Player %s is requesting game %s", ws.Request().RemoteAddr, msg)
		g, ok = games[msg]
		if !ok {
			websocket.Message.Send(ws, "notfound")
			return
		}
	}

	log.Printf("Adding Player %s to game %s", ws.Request().RemoteAddr, g.Id)
	// Send back game id
	g.AddPlayer(ws)
	websocket.Message.Send(ws, g.Id)
	if len(g.Players) == 1 {
		go g.Run()
	}

	defer func() {
		g.RemovePlayer(ws)
		if len(g.Players) == 0 {
			log.Println("Removing game", g.Id)
			delete(games, g.Id)
		}
	}()

	for {
		err := websocket.Message.Receive(ws, &msg)
		if err != nil {
			log.Printf("Player %s disconnected and left game %s", ws.Request().RemoteAddr, g.Id)
			return
		}
		switch msg {
		case "left":
			g.Input(ws, game.Left)
		case "right":
			g.Input(ws, game.Right)
		case "up":
			g.Input(ws, game.Up)
		case "down":
			g.Input(ws, game.Down)
		}
	}
}

func joinGame(res http.ResponseWriter) {
	file, err := os.Open("public/index.html")
	if err != nil {
		res.Write([]byte(err.Error()))
		return
	}
	io.Copy(res, file)
}

func main() {
	m := martini.New()

	m.Use(martini.Recovery())
	m.Use(martini.Logger())
	m.Use(martini.Static("public"))

	r := martini.NewRouter()

	r.Get("/connect", websocket.Handler(connect).ServeHTTP)
	r.Get("/:id", joinGame)

	m.Action(r.Handle)

	m.Run()
}
