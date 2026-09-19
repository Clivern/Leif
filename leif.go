// Copyright 2026 Clivern. All rights reserved.
// License can be found in the LICENSE file.

package main

import "github.com/clivern/leif/cli"

var (
	version = "dev"
	commit  = "none"
	date    = "unknown"
	builtBy = "unknown"
)

func main() {
	cli.Version = version
	cli.Commit = commit
	cli.Date = date
	cli.BuiltBy = builtBy
	cli.Execute()
}
