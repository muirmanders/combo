// Copyright 2015 Muir Manders.  All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package game

import "fmt"

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

			possibleDistance := sq.PieceCount - 1
			if possibleDistance == 0 {
				possibleDistance = 1
			}

			from := sq.Position

			addMovesInLine := func(to Position) {
				for _, psq := range squaresInALine(b, from, to) {
					if psq.PieceColor == c && psq.PieceCount > 0 {
						moves = append(moves, Move{from, psq.Position, false})
						break
					}

					if sq.PieceCount == 1 {
						break
					}

					moves = append(moves, Move{from, psq.Position, false})

					if psq.PieceCount == 0 || psq.PieceColor == c {
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
		}
	}

	return moves
}

func squaresInALine(b Board, from, to Position) []Square {
	var (
		variable *int
		distance int
	)

	if from.X == to.X {
		variable = &from.Y
		distance = to.Y - from.Y
	} else if from.Y == to.Y {
		variable = &from.X
		distance = to.X - from.X
	} else {
		panic(fmt.Sprintf("%+v -> %+v not in a line!", from, to))
	}

	neg := 1
	if distance < 0 {
		distance = -distance
		neg = -1
	}

	var ret []Square

	initial := *variable
	for i := 1; i <= distance; i++ {
		*variable = initial + neg*i
		square, err := b.Get(from)
		if err == nil {
			ret = append(ret, square)
		}
	}

	return ret
}
