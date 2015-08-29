// Copyright 2015 Muir Manders.  All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package contest

import (
	"combo/game"
	"fmt"
	"os"
	"sort"
	"time"
)

// download all submissions from s3, do round-robin tournament, add some sort of game log

type contestant struct {
	email      string
	submission string

	wins            int
	losses          int
	piecesRemaining int
	clockRemaining  time.Duration
}

type contestantSlice []contestant

func (c contestantSlice) Len() int {
	return len(c)
}

func (c contestantSlice) Swap(i, j int) {
	c[i], c[j] = c[j], c[i]
}

func (c contestantSlice) Less(i, j int) bool {
	c1 := c[i]
	c2 := c[j]

	if c1.wins != c2.wins {
		return c1.wins < c2.wins
	}

	if c1.piecesRemaining != c2.piecesRemaining {
		return c1.piecesRemaining < c2.piecesRemaining
	}

	return c1.clockRemaining > c2.clockRemaining
}

func PlayTournament() {
	var contestants []contestant

	logger, err := os.Create("/tmp/contest.log")
	if err != nil {
		panic(fmt.Sprintf("error creating log file: %s", err))
	}

	for i := range contestants {
		for j := range contestants {
			if i == j {
				continue
			}

			c1 := &contestants[i]
			c2 := &contestants[j]

			p1 := newPlayer(*c1)
			p2 := newPlayer(*c2)

			g, err := game.NewGame(game.Config{
				Black:     p1,
				White:     p2,
				Width:     8,
				Height:    8,
				Logger:    logger,
				GameClock: 5 * time.Minute,
			})

			if err != nil {
				panic(fmt.Sprintf("error creating game: %s", err))
			}

			gameRes := g.Play()

			var winner, loser *contestant
			if gameRes.Winner == p1 {
				winner = c1
				loser = c2
			} else {
				winner = c2
				loser = c1
			}

			winner.wins += 1
			winner.piecesRemaining += gameRes.WinnerPiecesRemaining
			winner.clockRemaining += gameRes.WinnerClock

			loser.losses += 1
			loser.piecesRemaining += gameRes.LoserPiecesRemaining
			loser.clockRemaining += gameRes.LoserClock

			p1.kill()
			p2.kill()
		}
	}

	sort.Sort(contestantSlice(contestants))

	for i, c := range contestants {
		fmt.Fprintf(logger, "%2d. %s %d wins, %d losses, %d pieces remaining, %s clock left\n", i+1, c.email, c.wins, c.losses, c.piecesRemaining, c.clockRemaining)
	}
}
