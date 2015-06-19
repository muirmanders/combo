// Copyright 2015 Muir Manders.  All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package main

import (
	"combo/cli"
	"combo/http"

	"github.com/spf13/cobra"
)

//go:generate go-bindata -o http/resources.go -prefix http/resources -pkg http http/resources/...

func main() {
	rootCmd := &cobra.Command{
		Use:   "combo",
		Short: "CLI interface for starting Combo frontends",
	}

	cliCmd := &cobra.Command{
		Use:   "cli",
		Short: "CLI frontend for Combo",
		Run: func(cmd *cobra.Command, args []string) {
			cli.Go()
		},
	}
	rootCmd.AddCommand(cliCmd)

	var httpListenAddr string
	httpCmd := &cobra.Command{
		Use:   "http",
		Short: "HTTP frontend for Combo",
		Run: func(cmd *cobra.Command, args []string) {
			http.Go(httpListenAddr)
		},
	}
	httpCmd.Flags().StringVarP(&httpListenAddr, "listen", "l", ":8080", "http server addr:port")
	rootCmd.AddCommand(httpCmd)

	rootCmd.Execute()
}
