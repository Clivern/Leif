// Copyright 2026 Clivern. All rights reserved.
// License can be found in the LICENSE file.

package cli

import (
	"context"
	"database/sql"

	"github.com/clivern/leif/core"
	"github.com/spf13/cobra"

	_ "github.com/lib/pq"
)

var serveCmd = &cobra.Command{
	Use:   "serve",
	Short: "Start the PostgreSQL MCP server over stdio",
	Run: func(_ *cobra.Command, _ []string) {
		conn, err := sql.Open("postgres", core.DSN(dsnFlag))
		if err != nil {
			panic(err)
		}

		err = conn.Ping()
		if err != nil {
			panic(err)
		}

		err = core.Serve(context.Background(), conn, Version)
		if err != nil {
			panic(err)
		}
	},
}

func init() {
	serveCmd.Flags().StringVarP(
		&dsnFlag,
		"dsn",
		"d",
		"",
		"PostgreSQL DSN (defaults to DATABASE_URL or PG* env vars)",
	)
	rootCmd.AddCommand(serveCmd)
}
