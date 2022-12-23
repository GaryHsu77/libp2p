package main

import (
	"fmt"
	"os"

	golog "github.com/ipfs/go-log/v2"
	"github.com/spf13/cobra"
)

var (
	debug = ""

	rootCmd = &cobra.Command{
		Use:   "lp",
		Short: "Peer Cli",
	}

	relayCmd = &cobra.Command{
		Use:   "relay",
		Short: "Relay Peer",
		Run:   relayRun,
	}

	listenCmd = &cobra.Command{
		Use:   "listen",
		Short: "Listener Peer",
		Run:   listenRun,
	}

	callCmd = &cobra.Command{
		Use:   "call",
		Short: "caller Peer",
		Run:   callRun,
	}
)

func init() {
	listenCmd.Flags().StringVarP(&debug, "debug", "d", "", "enable debug mode")
	rootCmd.AddCommand(relayCmd)
	rootCmd.AddCommand(listenCmd)
	rootCmd.AddCommand(callCmd)
}

func argInit() {
	if debug == "true" {
		golog.SetAllLoggers(golog.LevelDebug)
	}
}

func main() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
