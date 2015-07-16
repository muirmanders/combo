// Copyright 2015 Muir Manders.  All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package main

import (
	"combo/ai"
	"combo/cli"
	"combo/external"
	"combo/game"
	"combo/http"

	"github.com/spf13/cobra"
)

//go:generate go-bindata -o http/resources.go -prefix http/resources -pkg http http/resources/...

func main() {
	rootCmd := &cobra.Command{
		Use:   "combo",
		Short: "CLI interface for starting Combo games",
	}

	var cpuPlayer string

	whitePlayer := func() game.Player {
		switch cpuPlayer {
		case "random":
			return ai.NewRandomPlayer(game.White)
		case "weak":
			return ai.NewWeakPlayer(game.White)
		case "medium":
			return ai.NewMediumPlayer(game.White)
		default:
			return external.NewExternalPlayer(cpuPlayer, game.White)
		}
	}

	cliCmd := &cobra.Command{
		Use:   "cli",
		Short: "CLI frontend for Combo",
		Run: func(cmd *cobra.Command, args []string) {
			cli.Go(whitePlayer())
		},
	}
	cliCmd.Flags().StringVarP(&cpuPlayer, "cpu", "c", "medium", "random|weak|medium|/path/to/external/player")
	rootCmd.AddCommand(cliCmd)

	var httpListenAddr, certFile, keyFile string
	httpCmd := &cobra.Command{
		Use:   "http",
		Short: "HTTP frontend for Combo",
		Run: func(cmd *cobra.Command, args []string) {
			http.Go(httpListenAddr, whitePlayer(), certFile, keyFile)
		},
	}
	httpCmd.Flags().StringVarP(&httpListenAddr, "listen", "l", "localhost:8080", "http server addr:port")
	httpCmd.Flags().StringVarP(&cpuPlayer, "cpu", "c", "medium", "random|weak|medium|/path/to/external/player")
	httpCmd.Flags().StringVarP(&certFile, "cert", "", "", "certificate file")
	httpCmd.Flags().StringVarP(&keyFile, "key", "", "", "key file")
	rootCmd.AddCommand(httpCmd)

	rootCmd.Execute()
}
