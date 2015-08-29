// Copyright 2015 Muir Manders.  All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package game

import (
	"fmt"
	"io"
	"time"
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
	Play() GameResult
}

type GameResult struct {
	Winner                Player
	WinnerClock           time.Duration
	WinnerPiecesRemaining int

	Loser                Player
	LoserClock           time.Duration
	LoserPiecesRemaining int
}

type Config struct {
	Black     Player
	White     Player
	Width     int
	Height    int
	Logger    io.Writer
	GameClock time.Duration
}

func NewGame(config Config) (Game, error) {
	if config.Black == nil || config.White == nil {
		return nil, fmt.Errorf("must specify black and white players")
	}

	if config.Width <= 0 || config.Height < 4 {
		return nil, fmt.Errorf("must specify valid width and height")
	}

	return &game{
		board:      NewBoard(config.Width, config.Height).(*board),
		black:      config.Black,
		white:      config.White,
		turn:       config.Black,
		logger:     config.Logger,
		useClock:   config.GameClock > time.Duration(0),
		blackClock: config.GameClock,
		whiteClock: config.GameClock,
	}, nil
}

type game struct {
	board      *board
	black      Player
	white      Player
	turn       Player
	logger     io.Writer
	useClock   bool
	blackClock time.Duration
	whiteClock time.Duration
}

func (g *game) log(f string, args ...interface{}) {
	if g.logger == nil {
		return
	}
	fmt.Fprintf(g.logger, f, args...)
}

func piecesRemaining(b Board, c Color) int {
	remaining := 0

	for x := 0; x < b.Width(); x++ {
		for y := 0; y < b.Height(); y++ {

			sq, _ := b.Get(Position{X: x, Y: y})
			if sq.PieceColor == c {
				remaining += sq.PieceCount
			}
		}
	}

	return remaining
}

func (g *game) gameResult(winner Player) GameResult {
	var (
		loser       Player
		loserClock  time.Duration
		winnerClock time.Duration
	)
	if winner == g.black {
		loser = g.white
		winnerClock = g.blackClock
		loserClock = g.whiteClock
	} else {
		loser = g.black
		winnerClock = g.whiteClock
		loserClock = g.blackClock
	}

	return GameResult{
		Winner:                winner,
		WinnerClock:           winnerClock,
		WinnerPiecesRemaining: piecesRemaining(g.board, winner.Color()),

		Loser:                loser,
		LoserClock:           loserClock,
		LoserPiecesRemaining: piecesRemaining(g.board, loser.Color()),
	}
}

// Play the game until completion, returning the winning Player.
func (g *game) Play() GameResult {
	var (
		otherPlayer Player
		clock       *time.Duration
	)

	for {
		if g.turn == g.black {
			otherPlayer = g.white
			clock = &g.blackClock
		} else {
			otherPlayer = g.black
			clock = &g.whiteClock
		}

		availableMoves := g.board.AvailableMoves(g.turn.Color())
		if len(availableMoves) == 0 {
			g.log("Player %s (%s) loses for having no moves.\n", g.turn.Name(), g.turn.Color())
			return g.gameResult(otherPlayer)
		}

		var move Move

		if g.useClock {
			moveCh := make(chan Move)

			go func() {
				moveCh <- g.turn.Move(g.board)
			}()

			moveStart := time.Now()

			select {
			case move = <-moveCh:
				*clock -= time.Since(moveStart)
			case <-time.After(*clock):
				*clock = time.Duration(0)
				g.log("Player %s (%s) loses for running out of time.\n", g.turn.Name(), g.turn.Color())
				return g.gameResult(otherPlayer)
			}
		} else {
			move = g.turn.Move(g.board)
		}

		var moveOK bool
		for _, availableMove := range availableMoves {
			if move == availableMove {
				moveOK = true
				break
			}
		}

		if !moveOK {
			g.log("Player %s (%s) loses for illegal move %+v.\n", g.turn.Name(), g.turn.Color(), move)
			return g.gameResult(otherPlayer)
		}

		g.board.applyMove(move)

		g.turn = otherPlayer
	}
}
