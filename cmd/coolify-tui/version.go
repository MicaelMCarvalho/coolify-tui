package main

import (
	"fmt"
	"runtime/debug"
	"strings"
)

var (
	version = "dev"
	commit  = "none"
	date    = "unknown"
)

func printVersion() {
	fmt.Printf(
		"coolify-tui %s\ncommit: %s\ndate: %s\n",
		version,
		commit,
		date,
	)
}

func resolvedVersion() string {
	if version != "dev" {
		return strings.TrimPrefix(version, "v")
	}

	buildInfo, ok := debug.ReadBuildInfo()
	if !ok {
		return version
	}

	if buildInfo.Main.Version == "" ||
		buildInfo.Main.Version == "(devel)" {
		return version
	}

	return strings.TrimPrefix(
		buildInfo.Main.Version,
		"v",
	)
}
