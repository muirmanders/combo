// Copyright 2015 Muir Manders.  All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package http

import (
	"combo/game"
	"encoding/json"
	"fmt"
	"log"
	gohttp "net/http"
	"os"

	"github.com/gorilla/websocket"
)

var (
	singleGame game.Game
	cpuPlayer  game.Player
)

func Go(listenAddr string, cpu game.Player) {
	mux := gohttp.NewServeMux()

	cpuPlayer = cpu

	mux.HandleFunc("/", func(w gohttp.ResponseWriter, r *gohttp.Request) {
		w.Header().Set("Content-Type", "text/html")
		w.WriteHeader(gohttp.StatusOK)
		w.Write(MustAsset("combo.html"))
	})

	mux.HandleFunc("/combo.js", func(w gohttp.ResponseWriter, r *gohttp.Request) {
		w.Header().Set("Content-Type", "application/javascript")
		w.WriteHeader(gohttp.StatusOK)
		w.Write(MustAsset("combo.js"))
	})

	mux.HandleFunc("/combo.css", func(w gohttp.ResponseWriter, r *gohttp.Request) {
		w.Header().Set("Content-Type", "text/css")
		w.WriteHeader(gohttp.StatusOK)
		w.Write(MustAsset("combo.css"))
	})

	var upgrader = websocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
	}

	mux.HandleFunc("/connect", func(w gohttp.ResponseWriter, r *gohttp.Request) {
		conn, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			log.Printf("error making websocket: %s", err)
			return
		}
		handleConnect(conn)
	})

	fmt.Fprintf(os.Stdout, "Listening on %s...\n", listenAddr)
	panic(gohttp.ListenAndServe(listenAddr, mux))
}

type httpPlayer struct {
	color game.Color
	conn  *websocket.Conn
}

func (p httpPlayer) Color() game.Color {
	return p.color
}

func (p httpPlayer) Name() string {
	return "http"
}

type commandToClient struct {
	Command string      `json:"command"`
	Args    interface{} `json:"args"`
}

func (p httpPlayer) Move(b game.Board) game.Move {

	var move game.Move

	moveCommand := commandToClient{
		Command: "move",
		Args: map[string]interface{}{
			"board": b,
			"moves": b.AvailableMoves(p.Color()),
		},
	}

	if err := p.conn.WriteJSON(moveCommand); err != nil {
		log.Printf("error requesting move: %s", err)
		return move
	}

	var resp commandFromClient
	if err := p.conn.ReadJSON(&resp); err != nil {
		log.Printf("error receiving move: %s", err)
		return move
	}

	if resp.Command != "move" {
		log.Printf("expected move command, got: %s", resp.Command)
		return move
	}

	if err := json.Unmarshal(resp.Args, &move); err != nil {
		log.Printf("failed unmarshaling move: %s (%s)", resp.Command, string(resp.Args))
		return move
	}

	return move
}

type commandFromClient struct {
	Command string          `json:"command"`
	Args    json.RawMessage `json:"args"`
}

type newGameArgs struct {
	Width  int `json:"width"`
	Height int `json:"height"`
}

func handleConnect(conn *websocket.Conn) {
	player := httpPlayer{game.Black, conn}

	for {
		var c commandFromClient
		err := conn.ReadJSON(&c)
		if err != nil {
			log.Println("websocket error:", err)
			return
		}

		switch c.Command {
		case "new_game":
			var args newGameArgs
			if err := json.Unmarshal(c.Args, &args); err != nil {
				log.Printf("bad new_game args: %s (%s)", err, string(c.Args))
				return
			}

			singleGame, err = game.NewGame(game.Config{
				Black:  player,
				White:  cpuPlayer,
				Width:  args.Width,
				Height: args.Height,
				Logger: os.Stderr,
			})

			if err != nil {
				log.Printf("error creating game: %s", err)
				return
			}

			winner := singleGame.Play()

			gameOver := commandToClient{
				Command: "game_over",
				Args: map[string]string{
					"message": fmt.Sprintf("%s (%s) is the winner!", winner.Color(), winner.Name()),
				},
			}
			if err := conn.WriteJSON(gameOver); err != nil {
				log.Printf("error sending game_over: %s", err)
				return
			}
		}
	}
}
