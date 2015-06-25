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
	color game.Color
	path  string

	cmd           *exec.Cmd
	stdinEncoder  *json.Encoder
	stdoutDecoder *json.Decoder
}

func NewExternalPlayer(path string, color game.Color) game.Player {
	return &player{
		color: color,
		path:  path,
	}
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
	if p.cmd == nil {
		var err error
		p.cmd = exec.Command(p.path)

		stdin, err := p.cmd.StdinPipe()
		if err != nil {
			fmt.Printf("Error opening stdin pipe: %s\n", err)
			os.Exit(1)
		}
		p.stdinEncoder = json.NewEncoder(stdin)

		stdout, err := p.cmd.StdoutPipe()
		if err != nil {
			fmt.Printf("Error opening stdout pipe: %s\n", err)
			os.Exit(1)
		}
		p.stdoutDecoder = json.NewDecoder(stdout)

		if err = p.cmd.Start(); err != nil {
			fmt.Printf("Error starting external player: %s\n", err)
			os.Exit(1)
		}
	}

	if err := p.stdinEncoder.Encode(moveRequest{b, p.Color()}); err != nil {
		fmt.Printf("Error sending JSON board state: %s\n", err)
		os.Exit(1)
	}

	var move game.Move
	if err := p.stdoutDecoder.Decode(&move); err != nil {
		fmt.Printf("Error decoding move: %s\n", err)
		os.Exit(1)
	}

	return move
}
