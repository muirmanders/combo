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

type weakPlayer game.Color

func init() {
	rand.Seed(time.Now().UnixNano())
}

func NewWeakPlayer(c game.Color) game.Player {
	return weakPlayer(c)
}

func (p weakPlayer) Color() game.Color {
	return game.Color(p)
}

func (p weakPlayer) Name() string {
	return "weak player"
}

func (p weakPlayer) Move(b game.Board) game.Move {
	_, move := negamax(b, 1, math.MinInt32, math.MaxInt32, p.Color())
	return move
}
