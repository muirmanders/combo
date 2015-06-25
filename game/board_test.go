// Copyright 2015 Muir Manders.  All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package game

import (
	"reflect"
	"sort"
	"testing"
)

type moveList []Move

func (m moveList) Swap(i, j int) {
	m[i], m[j] = m[j], m[i]
}

func (m moveList) Len() int {
	return len(m)
}

func (m moveList) Less(i, j int) bool {
	if m[i].From.X != m[j].From.X {
		return m[i].From.X < m[j].From.X
	}
	if m[i].From.Y != m[j].From.Y {
		return m[i].From.Y < m[j].From.Y
	}
	if m[i].To.X != m[j].To.X {
		return m[i].To.X < m[j].To.X
	}
	if m[i].To.Y != m[j].To.Y {
		return m[i].To.Y < m[j].To.Y
	}
	return m[i].PieceCount < m[j].PieceCount
}

func TestAvailableMoves(t *testing.T) {
	b := NewBoard(2, 5).(*board)

	expected := []Move{
		{Position{0, 0}, Position{0, 1}, 1},
		{Position{0, 0}, Position{1, 0}, 1},
		{Position{0, 0}, Position{1, 1}, 1},

		{Position{1, 0}, Position{0, 0}, 1},
		{Position{1, 0}, Position{0, 1}, 1},
		{Position{1, 0}, Position{1, 1}, 1},

		{Position{0, 1}, Position{0, 0}, 1},
		{Position{0, 1}, Position{1, 0}, 1},
		{Position{0, 1}, Position{1, 1}, 1},
		{Position{0, 1}, Position{0, 2}, 1},
		{Position{0, 1}, Position{1, 2}, 1},

		{Position{1, 1}, Position{0, 0}, 1},
		{Position{1, 1}, Position{1, 0}, 1},
		{Position{1, 1}, Position{0, 1}, 1},
		{Position{1, 1}, Position{0, 2}, 1},
		{Position{1, 1}, Position{1, 2}, 1},
	}

	got := b.AvailableMoves(White)

	sort.Sort(moveList(expected))
	sort.Sort(moveList(got))

	if !reflect.DeepEqual(expected, got) {
		t.Errorf("got %+v, expected %+v", got, expected)
	}

	b.applyMove(Move{Position{0, 0}, Position{0, 1}, 1})

	expected = []Move{
		{Position{1, 0}, Position{0, 0}, 1},
		{Position{1, 0}, Position{0, 1}, 1},
		{Position{1, 0}, Position{1, 1}, 1},

		{Position{0, 1}, Position{0, 0}, 1},
		{Position{0, 1}, Position{0, 0}, 2},
		{Position{0, 1}, Position{1, 0}, 1},
		{Position{0, 1}, Position{1, 0}, 2},
		{Position{0, 1}, Position{1, 1}, 1},
		{Position{0, 1}, Position{1, 1}, 2},
		{Position{0, 1}, Position{0, 2}, 1},
		{Position{0, 1}, Position{0, 2}, 2},
		{Position{0, 1}, Position{1, 2}, 1},
		{Position{0, 1}, Position{1, 2}, 2},
		{Position{0, 1}, Position{0, 3}, 1},

		{Position{1, 1}, Position{0, 0}, 1},
		{Position{1, 1}, Position{1, 0}, 1},
		{Position{1, 1}, Position{0, 1}, 1},
		{Position{1, 1}, Position{0, 2}, 1},
		{Position{1, 1}, Position{1, 2}, 1},
	}
}
