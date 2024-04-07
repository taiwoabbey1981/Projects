package commands

import (
	"fmt"
	"os"

	"github.com/fatih/color"
	api "github.com/porter-dev/porter/api/client"
	"github.com/porter-dev/porter/cli/cmd/config"
	"github.com/porter-dev/porter/cli/cmd/utils"

	"github.com/spf13/cobra"
)

func registerCommand_Open(cliConf config.CLIConfig) *cobra.Command {
	openCmd := &cobra.Command{
		Use:   "open",
		Short: "Opens the browser at the currently set Porter instance",
		Run: func(cmd *cobra.Command, args []string) {
			ctx := cmd.Context()

			client, err := api.NewClientWithConfig(ctx, api.NewClientInput{
				BaseURL:        fmt.Sprintf("%s/api", cliConf.Host),
				BearerToken:    cliConf.Token,
				CookieFileName: "cookie.json",
			})
			if err != nil {
				_, _ = color.New(color.FgRed).Fprintf(os.Stderr, "error creating porter API client: %v\n", err)
				os.Exit(1)
			}

			user, err := client.AuthCheck(ctx)
			if err != nil {
				_ = utils.OpenBrowser(fmt.Sprintf("%s/register", cliConf.Host))
				return
			}

			_ = utils.OpenBrowser(fmt.Sprintf("%s/login?email=%s", cliConf.Host, user.Email)) //nolint:errcheck,gosec // do not want to change logic of CLI. New linter error
		},
	}

	return openCmd
}
