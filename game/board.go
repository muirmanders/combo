// Copyright 2015 Muir Manders.  All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package game

import "errors"

type board struct {
	squares [][]Square `json:"squares"`
	width   int        `json:"width"`
	height  int        `json:"height"`
}

func newBoard(width, height int) *board {
	b := &board{
		squares: make([][]Square, height),
		width:   width,
		height:  height,
	}

	for y := 0; y < height; y++ {
		b.squares[y] = make([]Square, width)

		for x := 0; x < width; x++ {
			b.squares[y][x].X = x
			b.squares[y][x].Y = y
		}

		if y <= 1 {
			for x := 0; x < width; x++ {
				b.squares[y][x].PieceCount = 1
				b.squares[y][x].PieceColor = White
			}
		}

		if y >= height-2 {
			for x := 0; x < width; x++ {
				b.squares[y][x].PieceCount = 1
				b.squares[y][x].PieceColor = Black
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
	return b.squares[p.Y][p.X], nil
}

func (b *board) Width() int {
	return b.width
}

func (b *board) Height() int {
	return b.height
}
