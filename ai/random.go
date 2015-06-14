// Copyright 2015 Muir Manders.  All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package ai

import (
	"combo/game"
	"math/rand"
)

type randomPlayer game.Color

func NewRandomPlayer(c game.Color) game.Player {
	return randomPlayer(c)
}

func (p randomPlayer) Color() game.Color {
	return game.Color(p)
}

func (p randomPlayer) Name() string {
	return "random player"
}

func (p randomPlayer) Move(b game.Board) game.Move {
	moves := game.AvailableMoves(b, p.Color())
	return moves[rand.Intn(len(moves))]
}
