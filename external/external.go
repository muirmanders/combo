// Copyright 2015 Muir Manders.  All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package external

import (
	"combo/game"
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
)

type player struct {
	color         game.Color
	cmd           *exec.Cmd
	stdinEncoder  *json.Encoder
	stdoutDecoder *json.Decoder
}

func NewExternalPlayer(color game.Color, path string, args ...string) (game.Player, error) {
	cmd := exec.Command(path, args...)

	stdin, err := cmd.StdinPipe()
	if err != nil {
		return nil, fmt.Errorf("error opening stdin pipe: %s", err)
	}

	stdout, err := cmd.StdoutPipe()
	if err != nil {
		return nil, fmt.Errorf("error opening stdout pipe: %s", err)
	}

	if err := cmd.Start(); err != nil {
		return nil, fmt.Errorf("error starting external player: %s", err)
	}

	return &player{
		color:         color,
		cmd:           cmd,
		stdinEncoder:  json.NewEncoder(stdin),
		stdoutDecoder: json.NewDecoder(stdout),
	}, nil
}

func (p *player) Color() game.Color {
	return p.color
}

func (p *player) Name() string {
	return "external"
}

type moveRequest struct {
	Board game.Board `json:"board"`
	Color game.Color `json:"color"`
}

func (p *player) Move(b game.Board) game.Move {
	var move game.Move

	if err := p.stdinEncoder.Encode(moveRequest{b, p.Color()}); err != nil {
		fmt.Fprintf(os.Stderr, "Error sending JSON board state: %s\n", err)
		return move
	}

	if err := p.stdoutDecoder.Decode(&move); err != nil {
		fmt.Fprintf(os.Stderr, "Error decoding move: %s\n", err)
	}

	return move
}
