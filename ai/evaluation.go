// Copyright 2015 Muir Manders.  All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package ai

import (
	"combo/game"
	"math"
	"math/rand"
)

// simple heuristic for board evaluation
func scoreBoard(b game.Board, c game.Color) (score int) {

	var myPieceCount, theirPieceCount int

	for x := 0; x < b.Width(); x++ {
		for y := 0; y < b.Height(); y++ {
			sq, _ := b.Get(game.Position{x, y})
			if sq.PieceColor == c {
				myPieceCount += sq.PieceCount
				if sq.PieceCount == 1 || sq.PieceCount > 4 {
					score--
				}
			} else {
				theirPieceCount += sq.PieceCount
				if sq.PieceCount == 1 || sq.PieceCount > 4 {
					score++
				}
			}
		}
	}

	return score + 5*int(float64(myPieceCount-theirPieceCount)*float64(2*b.Width())/float64(myPieceCount)+0.5)
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
		}
		if val > alpha {
			alpha = val
		}
		if alpha >= beta {
			break
		}
	}

	return bestScore, bestMove
}
