// Copyright 2015 Muir Manders.  All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package game

import "testing"

func TestApplyMove(t *testing.T) {
	b := NewBoard(2, 4).(*board)

	checkSquare := func(x, y, count int, color Color) {
		sq, _ := b.Get(Position{x, y})
		if c := sq.PieceCount; c != count {
			t.Errorf("%d,%d was %d", x, y, c)
		}
		if c := sq.PieceColor; c != color {
			t.Errorf("%d,%d was %s", x, y, c)
		}
	}

	b.applyMove(Move{Position{0, 0}, Position{0, 1}, 1})
	checkSquare(0, 0, 0, emptySquare)
	checkSquare(0, 1, 2, White)

	b.applyMove(Move{Position{0, 1}, Position{0, 2}, 2})
	checkSquare(0, 1, 0, emptySquare)
	checkSquare(0, 2, 2, White)

	b.applyMove(Move{Position{0, 2}, Position{0, 1}, 1})
	checkSquare(0, 2, 1, White)
	checkSquare(0, 1, 1, White)
}
