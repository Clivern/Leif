// Copyright 2026 Clivern. All rights reserved.
// License can be found in the LICENSE file.

package core

import (
	"fmt"
	"os"
)

// DSN returns a PostgreSQL connection string. An explicit override wins,
// then DATABASE_URL, then PG* environment variables.
func DSN(override string) string {
	if override != "" {
		return override
	}

	if url := os.Getenv("DATABASE_URL"); url != "" {
		return url
	}

	user := env("PGUSER", os.Getenv("USER"))
	if user == "" {
		user = "postgres"
	}

	return fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		env("PGHOST", "127.0.0.1"),
		env("PGPORT", "5432"),
		user,
		os.Getenv("PGPASSWORD"),
		env("PGDATABASE", "postgres"),
	)
}

func env(key, fallback string) string {
	val := os.Getenv(key)
	if val == "" {
		return fallback
	}
	return val
}
