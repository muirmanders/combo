// Copyright 2015 Muir Manders.  All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package game

import "errors"

type board struct {
	squares []Square
	width   int
	height  int

	blackCount int
	whiteCount int
}

func NewBoard(width, height int) Board {
	b := &board{
		squares: make([]Square, width*height),
		width:   width,
		height:  height,
	}

	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			sq := &b.squares[y*width+x]

			sq.X = x
			sq.Y = y

			if y <= 1 {
				sq.PieceCount = 1
				sq.PieceColor = White
				b.whiteCount++
			} else if y >= height-2 {
				sq.PieceCount = 1
				sq.PieceColor = Black
				b.blackCount++
			}
		}
	}

	return b
}

func (b *board) PieceCount(c Color) int {
	if c == Black {
		return b.blackCount
	} else {
		return b.whiteCount
	}
}

var ErrOutOfBounds = errors.New("out of bounds")

func (b *board) Get(p Position) (Square, error) {
	if p.X < 0 || p.X >= b.width || p.Y < 0 || p.Y >= b.height {
		return Square{}, ErrOutOfBounds
	}
	return b.squares[p.Y*b.width+p.X], nil
}

func (b *board) IfMove(m Move) Board {
	dupe := *b
	dupe.squares = make([]Square, len(b.squares))
	copy(dupe.squares, b.squares)
	dupe.applyMove(m)
	return &dupe
}

func (b *board) AvailableMoves(c Color) []Move {
	var moves []Move

	for x := 0; x < b.width; x++ {
		for y := 0; y < b.height; y++ {

			sq := b.squares[y*b.width+x]

			if sq.PieceColor != c || sq.PieceCount == 0 {
				continue
			}

			for dx := -1; dx <= 1; dx++ {
				for dy := -1; dy <= 1; dy++ {

					if dx == 0 && dy == 0 {
						continue
					}

					x, y := sq.X, sq.Y

					for i := 1; i <= sq.PieceCount; i++ {
						x += dx
						y += dy

						if x < 0 || x >= b.width || y < 0 || y >= b.height {
							break
						}

						psq := b.squares[x+y*b.width]

						if sq.PieceCount > 1 || psq.PieceColor == c || psq.PieceCount == 0 {
							moves = append(moves, Move{sq.Position, psq.Position, false})
						}

						if sq.PieceCount > 1 && i == 1 && (psq.PieceCount == 0 || psq.PieceColor == c) {
							moves = append(moves, Move{sq.Position, psq.Position, true})
						}

						if psq.PieceCount > 0 {
							break
						}
					}
				}
			}
		}
	}

	return moves
}

func (b *board) applyMove(move Move) {
	fromSq := &b.squares[move.From.Y*b.width+move.From.X]
	toSq := &b.squares[move.To.Y*b.width+move.To.X]

	if move.Split {
		toSq.PieceColor = fromSq.PieceColor
		toSq.PieceCount++
		fromSq.PieceCount--
	} else {
		if toSq.PieceColor == fromSq.PieceColor {
			toSq.PieceCount += fromSq.PieceCount
		} else {
			if toSq.PieceColor == Black {
				b.blackCount -= toSq.PieceCount
			} else {
				b.whiteCount -= toSq.PieceCount
			}
			toSq.PieceCount = fromSq.PieceCount
			toSq.PieceColor = fromSq.PieceColor
		}
		fromSq.PieceCount = 0
	}
}

func (b *board) Width() int {
	return b.width
}

func (b *board) Height() int {
	return b.height
}
