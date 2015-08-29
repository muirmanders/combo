// Copyright 2015 Muir Manders.  All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package contest

import "combo/game"

type player struct {
	name  string
	color game.Color
}

func (p *player) Name() string {
	return p.name
}

func (p *player) Move(b game.Board) game.Move {
	return game.Move{}
}

func (p *player) Color() game.Color {
	return p.color
}

func (p *player) kill() {
}

func newPlayer(c contestant) *player {
	return new(player)
}
