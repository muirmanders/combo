// Copyright 2015 Muir Manders.  All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package game

import (
	"fmt"
	"io"
)

type Color string

const (
	White Color = "white"
	Black Color = "black"

	emptySquare Color = ""
)

func (c Color) String() string {
	return string(c)
}

func (c Color) Other() Color {
	if c == Black {
		return White
	} else {
		return Black
	}
}

type Position struct {
	X int `json:"x"`
	Y int `json:"y"`
}

type Square struct {
	Position

	PieceColor Color `json:"piece_color,omitempty"`
	PieceCount int   `json:"piece_count"`
}

type Board interface {
	Width() int
	Height() int
	Get(Position) (Square, error)
	AvailableMoves(Color) []Move
	IfMove(Move) Board
}

type Move struct {
	From       Position `json:"from"`
	To         Position `json:"to"`
	PieceCount int      `json:"piece_count"`
}

type Player interface {
	Name() string
	Color() Color
	Move(Board) Move
}

type Game interface {
	Play() Player
}

type Config struct {
	Black  Player
	White  Player
	Width  int
	Height int
	Logger io.Writer
}

func NewGame(config Config) (Game, error) {
	if config.Black == nil || config.White == nil {
		return nil, fmt.Errorf("must specify black and white players")
	}

	if config.Width <= 0 || config.Height < 4 {
		return nil, fmt.Errorf("must specify valid width and height")
	}

	return &game{
		board:  NewBoard(config.Width, config.Height).(*board),
		black:  config.Black,
		white:  config.White,
		turn:   config.Black,
		logger: config.Logger,
	}, nil
}

type game struct {
	board  *board
	black  Player
	white  Player
	turn   Player
	logger io.Writer
}

func (g *game) log(f string, args ...interface{}) {
	if g.logger == nil {
		return
	}
	fmt.Fprintf(g.logger, f, args...)
}

// Play the game until completion, returning the winning Player.
func (g *game) Play() Player {
	var otherPlayer Player

	for {
		if g.turn == g.black {
			otherPlayer = g.white
		} else {
			otherPlayer = g.black
		}

		availableMoves := g.board.AvailableMoves(g.turn.Color())
		if len(availableMoves) == 0 {
			g.log("Player %s (%s) loses for having no moves.\n", g.turn.Name(), g.turn.Color())
			return otherPlayer
		}

		move := g.turn.Move(g.board)

		var moveOK bool
		for _, availableMove := range availableMoves {
			if move == availableMove {
				moveOK = true
				break
			}
		}

		if !moveOK {
			g.log("Player %s (%s) loses for illegal move %+v.\n", g.turn.Name(), g.turn.Color(), move)
			return otherPlayer
		}

		g.board.applyMove(move)

		g.turn = otherPlayer
	}
}
