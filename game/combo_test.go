// Copyright 2015 Muir Manders.  All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package game

import "testing"

func TestApplyMove(t *testing.T) {
	b := newBoard(2, 4)

	checkSquare := func(x, y, count int, color Color) {
		if c := b.squares[y][x].PieceCount; c != count {
			t.Errorf("%d,%y was %d", x, y, c)
		}
		if c := b.squares[y][x].PieceColor; c != color {
			t.Errorf("%d,%y was %s", x, y, c)
		}
	}

	applyMove(b, Move{Position{0, 0}, Position{0, 1}, false})
	checkSquare(0, 0, 0, White)
	checkSquare(0, 1, 2, White)

	applyMove(b, Move{Position{0, 1}, Position{0, 2}, false})
	checkSquare(0, 1, 0, White)
	checkSquare(0, 2, 2, White)

	applyMove(b, Move{Position{0, 2}, Position{0, 1}, true})
	checkSquare(0, 2, 1, White)
	checkSquare(0, 1, 1, White)
}
