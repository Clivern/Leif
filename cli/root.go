// Copyright 2026 Clivern. All rights reserved.
// License can be found in the LICENSE file.

package cli

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var dsnFlag string

var rootCmd = &cobra.Command{
	Use:   "leif",
	Short: "Stdio MCP server for a local PostgreSQL database",
}

// Execute runs the leif CLI.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
