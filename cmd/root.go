package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "handoff",
	Short: "Transfer knowledge between agent sessions",
	Long:  "A CLI tool for storing and retrieving knowledge packages across AI agent context windows.",
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
