//go:build cli
// +build cli

package main

import (
	"context"
	"os"
	"time"

	"github.com/fatih/color"
	"github.com/getsentry/sentry-go"
	"github.com/porter-dev/porter/cli/cmd/commands"
	"github.com/porter-dev/porter/cli/cmd/config"
	"github.com/porter-dev/porter/cli/cmd/errors"
)

func main() {
	ctx := context.Background()
	if errors.SentryDSN != "" {
		err := sentry.Init(sentry.ClientOptions{
			Dsn:         errors.SentryDSN,
			Environment: "cli",
			Debug:       config.Version == "dev",
			Release:     config.Version,
			IgnoreErrors: []string{
				"Forbidden",
			},
		})
		if err != nil {
			color.New(color.FgRed).Fprintf(os.Stderr, "error initialising sentry: %s\n", err)
			os.Exit(1)
		}

		defer sentry.Flush(2 * time.Second)
	}

	err := commands.Execute(ctx)
	if err != nil {
		color.New(color.FgRed).Fprintf(os.Stderr, "error executing command: %s\n", err)
		os.Exit(1)
	}
}
