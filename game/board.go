// Copyright 2015 Muir Manders.  All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package game

import "errors"

type board struct {
	squares []Square `json:"squares"`
	width   int      `json:"width"`
	height  int      `json:"height"`
}

func newBoard(width, height int) *board {
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
			} else if y >= height-2 {
				sq.PieceCount = 1
				sq.PieceColor = Black
			}
		}
	}

	return b
}

func (b *board) mustGet(p Position) Square {
	sq, err := b.Get(p)
	if err != nil {
		panic(err)
	}
	return sq
}

var ErrOutOfBounds = errors.New("out of bounds")

func (b *board) Get(p Position) (Square, error) {
	if p.X < 0 || p.X >= b.width || p.Y < 0 || p.Y >= b.height {
		return Square{}, ErrOutOfBounds
	}
	return b.squares[p.Y*b.height+p.X], nil
}

func (b *board) IfMove(m Move) Board {
	dupe := *b
	dupe.squares = make([]Square, len(b.squares))
	copy(dupe.squares, b.squares)
	dupe.applyMove(m)
	return &dupe
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
