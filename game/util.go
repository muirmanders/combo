// Copyright 2015 Muir Manders.  All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package game

func AvailableMoves(b Board, c Color) []Move {
	var moves []Move

	for x := 0; x < b.Width(); x++ {
		for y := 0; y < b.Height(); y++ {

			sq, err := b.Get(Position{x, y})
			if err != nil {
				panic(err)
			}

			if sq.PieceColor != c || sq.PieceCount == 0 {
				continue
			}

			possibleDistance := sq.PieceCount
			from := sq.Position

			addMovesInLine := func(to Position) {
				for pi, psq := range squaresInALine(b, from, to) {
					if sq.PieceCount > 1 || psq.PieceColor == c || psq.PieceCount == 0 {
						moves = append(moves, Move{from, psq.Position, false})
					}

					if sq.PieceCount > 1 && pi == 0 && (psq.PieceCount == 0 || psq.PieceColor == c) {
						moves = append(moves, Move{from, psq.Position, true})
					}

					if psq.PieceCount > 0 {
						break
					}
				}
			}

			addMovesInLine(Position{x, y + possibleDistance})
			addMovesInLine(Position{x, y - possibleDistance})
			addMovesInLine(Position{x + possibleDistance, y})
			addMovesInLine(Position{x - possibleDistance, y})

			addMovesInLine(Position{x + possibleDistance, y + possibleDistance})
			addMovesInLine(Position{x + possibleDistance, y - possibleDistance})
			addMovesInLine(Position{x - possibleDistance, y + possibleDistance})
			addMovesInLine(Position{x - possibleDistance, y - possibleDistance})
		}
	}

	return moves
}

func squaresInALine(b Board, from, to Position) []Square {
	dx := to.X - from.X
	dy := to.Y - from.Y

	distance := dx
	if distance == 0 {
		distance = dy
	}

	if distance < 0 {
		distance = -distance
	}

	dx /= distance
	dy /= distance

	var ret []Square
	for i := 0; i < distance; i++ {
		from.X += dx
		from.Y += dy
		square, err := b.Get(from)
		if err == nil {
			ret = append(ret, square)
		}
	}

	return ret
}
