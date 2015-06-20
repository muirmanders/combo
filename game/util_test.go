// Copyright 2015 Muir Manders.  All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package game

import (
	"reflect"
	"testing"
)

func TestSquaresInALine(t *testing.T) {
	b := newBoard(3, 3)

	var expected []Square

	got := squaresInALine(b, Position{0, 0}, Position{0, 0})
	if !reflect.DeepEqual(got, expected) {
		t.Errorf("got %v, expected %v", got, expected)
	}

	got = squaresInALine(b, Position{0, 0}, Position{-2, 0})
	if !reflect.DeepEqual(got, expected) {
		t.Errorf("got %v, expected %v", got, expected)
	}

	expected = []Square{
		b.mustGet(Position{0, 1}),
	}
	got = squaresInALine(b, Position{0, 0}, Position{0, 1})
	if !reflect.DeepEqual(got, expected) {
		t.Errorf("got %v, expected %v", got, expected)
	}

	expected = []Square{
		b.mustGet(Position{1, 0}),
	}
	got = squaresInALine(b, Position{1, 1}, Position{1, 0})
	if !reflect.DeepEqual(got, expected) {
		t.Errorf("got %v, expected %v", got, expected)
	}

	expected = []Square{
		b.mustGet(Position{1, 2}),
		b.mustGet(Position{2, 2}),
	}
	got = squaresInALine(b, Position{0, 2}, Position{4, 2})
	if !reflect.DeepEqual(got, expected) {
		t.Errorf("got %v, expected %v", got, expected)
	}
}

func TestAvailableMoves(t *testing.T) {
	b := newBoard(2, 4)

	expected := []Move{
		{Position{0, 0}, Position{0, 1}, false},
		{Position{0, 0}, Position{1, 0}, false},
		{Position{0, 1}, Position{0, 0}, false},
		{Position{0, 1}, Position{1, 1}, false},
		{Position{1, 0}, Position{1, 1}, false},
		{Position{1, 0}, Position{0, 0}, false},
		{Position{1, 1}, Position{1, 0}, false},
		{Position{1, 1}, Position{0, 1}, false},
	}

	got := AvailableMoves(b, White)

	if !reflect.DeepEqual(expected, got) {
		t.Errorf("got %+v, expected %+v", got, expected)
	}

	b.applyMove(Move{Position{0, 0}, Position{0, 1}, false})

	expected = []Move{
		{Position{0, 1}, Position{0, 2}, false},
		{Position{0, 1}, Position{0, 0}, false},
		{Position{0, 1}, Position{0, 0}, true},
		{Position{0, 1}, Position{1, 1}, false},
		{Position{1, 0}, Position{1, 1}, false},
		{Position{1, 1}, Position{1, 0}, false},
		{Position{1, 1}, Position{0, 1}, false},
	}

	got = AvailableMoves(b, White)

	if !reflect.DeepEqual(expected, got) {
		t.Errorf("got %+v, expected %+v", got, expected)
	}

	// merge to triple
	b.applyMove(Move{Position{0, 1}, Position{1, 1}, false})

	// make room on black side for longer distance white move
	b.applyMove(Move{Position{1, 2}, Position{1, 3}, false})

	expected = []Move{
		{Position{1, 0}, Position{1, 1}, false},
		{Position{1, 1}, Position{1, 2}, false},
		{Position{1, 1}, Position{1, 2}, true},
		{Position{1, 1}, Position{1, 3}, false},
		{Position{1, 1}, Position{1, 0}, false},
		{Position{1, 1}, Position{0, 1}, false},
		{Position{1, 1}, Position{0, 1}, true},
	}

	got = AvailableMoves(b, White)

	if !reflect.DeepEqual(expected, got) {
		t.Errorf("got %+v, expected %+v", got, expected)
	}

}
