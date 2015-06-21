// Copyright 2015 Muir Manders.  All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package ai

import (
	"combo/game"
	"math"
	"math/rand"
)

func abs(i int) int {
	if i < 0 {
		return -i
	}
	return i
}

// simple heuristic to maximize # of pieces vs opponent
func scoreBoard(b game.Board, c game.Color) int {
	return b.PieceCount(c) - b.PieceCount(c.Other())
}

var nullMove game.Move

func negamax(b game.Board, depth, alpha, beta int, color game.Color) (int, game.Move) {
	if depth == 0 {
		return scoreBoard(b, color), nullMove
	}

	moves := b.AvailableMoves(color)

	if len(moves) == 0 {
		return math.MinInt32 + 1, nullMove
	}

	var (
		bestScore = math.MinInt32
		bestMove  game.Move
	)

	for _, i := range rand.Perm(len(moves)) {
		val, _ := negamax(b.IfMove(moves[i]), depth-1, -beta, -alpha, color.Other())
		val = -val
		if val > bestScore {
			bestScore = val
			bestMove = moves[i]
			if val > alpha {
				alpha = val
			}
			if alpha >= beta {
				break
			}
		}
	}

	return bestScore, bestMove
}
