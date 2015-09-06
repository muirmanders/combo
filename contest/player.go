// Copyright 2015 Muir Manders.  All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package contest

import (
	"combo/external"
	"combo/game"
	"fmt"
	"os/exec"
	"strings"
)

type player struct {
	name           string
	externalPlayer game.Player
	containerID    string
}

func (p *player) Name() string {
	return p.name
}

func (p *player) Move(b game.Board) game.Move {
	return p.externalPlayer.Move(b)
}

func (p *player) Color() game.Color {
	return p.externalPlayer.Color()
}

func (p *player) kill() {
	exec.Command("docker", "kill", p.containerID).Run()
}

func newPlayer(c contestant, color game.Color) (*player, error) {

	var cpus string
	if color == game.Black {
		cpus = "0,1"
	} else {
		cpus = "2,3"
	}

	args := []string{
		"run",
		"--interactive",
		"--tty",
		"--rm",
		"--net=none",
		"--memory=8G",
		"--cpuset-cpus=" + cpus,
		"--volume=" + c.submission + ":/combo",
		"combo",
		"/combo/start",
	}

	ep, err := external.NewExternalPlayer(color, "docker", args...)

	if err != nil {
		return nil, err
	}

	containerIDBytes, err := exec.Command("docker", "ps", "--latest", "--quiet").CombinedOutput()
	if err != nil {
		return nil, fmt.Errorf("error getting latest image id: %s (%s)", err, containerIDBytes)
	}

	return &player{
		name:           c.email,
		externalPlayer: ep,
		containerID:    strings.TrimSpace(string(containerIDBytes)),
	}, nil
}
