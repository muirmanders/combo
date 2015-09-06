// Copyright 2015 Muir Manders.  All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package main

import (
	"combo/ai"
	"combo/cli"
	"combo/contest"
	"combo/external"
	"combo/game"
	"combo/http"
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

//go:generate go-bindata -o http/resources.go -prefix http/resources -pkg http http/resources/...

func main() {
	rootCmd := &cobra.Command{
		Use:   "combo",
		Short: "CLI interface for starting Combo games",
	}

	var whitePlayer, blackPlayer string

	makePlayer := func(playerName string, c game.Color) game.Player {
		switch playerName {
		case "human":
			return nil
		case "random":
			return ai.NewRandomPlayer(c)
		case "weak":
			return ai.NewWeakPlayer(c)
		case "medium":
			return ai.NewMediumPlayer(c)
		default:
			player, err := external.NewExternalPlayer(c, playerName)
			if err != nil {
				fmt.Fprintf(os.Stderr, "Error creating external player: %s", err)
				os.Exit(1)
			}
			return player
		}
	}

	cliCmd := &cobra.Command{
		Use:   "cli",
		Short: "CLI frontend for Combo",
		Run: func(cmd *cobra.Command, args []string) {
			cli.Go(makePlayer(whitePlayer, game.White), makePlayer(blackPlayer, game.Black))
		},
	}
	cliCmd.Flags().StringVarP(&whitePlayer, "cpu", "c", "medium", "random|weak|medium|/path/to/external/player")
	cliCmd.Flags().StringVarP(&whitePlayer, "white", "w", "medium", "random|weak|medium|/path/to/external/player")
	cliCmd.Flags().StringVarP(&blackPlayer, "black", "b", "human", "human|random|weak|medium|/path/to/external/player")
	rootCmd.AddCommand(cliCmd)

	var httpListenAddr, certFile, keyFile string
	httpCmd := &cobra.Command{
		Use:   "http",
		Short: "HTTP frontend for Combo",
		Run: func(cmd *cobra.Command, args []string) {
			http.Go(httpListenAddr, makePlayer(whitePlayer, game.White), certFile, keyFile)
		},
	}
	httpCmd.Flags().StringVarP(&httpListenAddr, "listen", "l", "localhost:8080", "http server addr:port")
	httpCmd.Flags().StringVarP(&whitePlayer, "cpu", "c", "medium", "random|weak|medium|/path/to/external/player")
	httpCmd.Flags().StringVarP(&certFile, "cert", "", "", "certificate file")
	httpCmd.Flags().StringVarP(&keyFile, "key", "", "", "key file")
	rootCmd.AddCommand(httpCmd)

	var submissionsDir string
	contestCmd := &cobra.Command{
		Use:   "contest",
		Short: "Round-robin tournament",
		Run: func(cmd *cobra.Command, args []string) {
			if err := contest.PlayTournament(submissionsDir); err != nil {
				fmt.Printf("error playing tournament: %s\n", err)
				os.Exit(1)
			}
		},
	}
	contestCmd.Flags().StringVarP(&submissionsDir, "submissions", "s", "", "directory containing extracted submissions")
	rootCmd.AddCommand(contestCmd)

	rootCmd.Execute()
}
