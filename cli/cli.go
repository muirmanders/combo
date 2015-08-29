// Copyright 2015 Muir Manders.  All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package cli

import (
	"bufio"
	"bytes"
	"combo/game"
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"strings"
)

type cliPlayer struct {
	color game.Color
}

func Go(whitePlayer, blackPlayer game.Player) {
	if blackPlayer == nil {
		blackPlayer = cliPlayer{game.Black}
	}

	config := game.Config{
		Black: blackPlayer,
		White: whitePlayer,

		Width:  8,
		Height: 8,
	}

	g, err := game.NewGame(config)
	if err != nil {
		panic(err)
	}

	result := g.Play()

	fmt.Fprintf(os.Stdout, "%s is the winner with %d remaining!\n", result.Winner.Color(), result.WinnerPiecesRemaining)
}

func (p cliPlayer) Color() game.Color {
	return p.color
}

func (p cliPlayer) Name() string {
	return p.color.String()
}

func (p cliPlayer) Move(b game.Board) game.Move {
	buf := new(bytes.Buffer)

	buf.WriteString(" " + strings.Repeat("-", 5*b.Width()-1) + "\n")

	for y := 0; y < b.Height(); y++ {
		buf.WriteByte('|')
		for x := 0; x < b.Width(); x++ {
			sq, _ := b.Get(game.Position{x, y})
			if sq.PieceCount == 0 {
				buf.WriteString("    ")
			} else {
				var c string
				if sq.PieceColor == game.Black {
					c = "b"
				} else {
					c = "w"
				}

				num := strconv.Itoa(sq.PieceCount)
				if len(num) == 1 {
					buf.WriteString(" " + num + c + " ")
				} else {
					buf.WriteString(num + c + " ")
				}
			}

			buf.WriteByte('|')
		}

		buf.WriteString(" " + strconv.Itoa(y) + "\n")
		buf.WriteString(" " + strings.Repeat("-", 5*b.Width()-1) + "\n")
	}

	buf.WriteString(" ")
	for x := 0; x < b.Width(); x++ {
		fmt.Fprintf(buf, " %2d  ", x)
	}
	buf.WriteString("\n\n")

	clear := exec.Command("clear")
	clear.Stdout = os.Stdout
	clear.Run()

	os.Stdout.Write(buf.Bytes())

	for {
		var (
			move game.Move

			n   int
			err error
		)

		for n < 4 || err != nil {
			fmt.Fprintf(os.Stdout, "Enter move for %s (fromx,fromy tox,toy count): ", p.Name())
			scanner := bufio.NewScanner(os.Stdin)
			scanner.Scan()
			n, err = fmt.Sscanf(scanner.Text(), "%d,%d %d,%d %d",
				&move.From.X,
				&move.From.Y,
				&move.To.X,
				&move.To.Y,
				&move.PieceCount,
			)
		}

		available := b.AvailableMoves(p.Color())
		for _, m := range available {
			if move == m {
				return move
			}
		}

		fmt.Fprintf(os.Stdout, "Move %+v is not legal!\n", move)
	}
}
