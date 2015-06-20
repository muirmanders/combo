// Copyright 2015 Muir Manders.  All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package ai

import (
	"combo/game"
	"math"
	"math/rand"
	"time"
)

type mediumPlayer game.Color

func init() {
	rand.Seed(time.Now().UnixNano())
}

func NewMediumPlayer(c game.Color) game.Player {
	return mediumPlayer(c)
}

func (p mediumPlayer) Color() game.Color {
	return game.Color(p)
}

func (p mediumPlayer) Name() string {
	return "medium player"
}

func (p mediumPlayer) Move(b game.Board) game.Move {
	_, move := negamax(b, 4, math.MinInt32, math.MaxInt32, p.Color())
	return move
}
