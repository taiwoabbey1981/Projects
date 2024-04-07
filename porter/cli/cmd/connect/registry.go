package connect

import (
	"context"
	"fmt"

	"github.com/porter-dev/porter/api/types"

	"github.com/fatih/color"
	api "github.com/porter-dev/porter/api/client"
	"github.com/porter-dev/porter/cli/cmd/utils"
)

// Helm connects a Helm repository using HTTP basic authentication
func Registry(
	ctx context.Context,
	client api.Client,
	projectID uint,
) (uint, error) {
	// if project ID is 0, ask the user to set the project ID or create a project
	if projectID == 0 {
		return 0, fmt.Errorf("no project set, please run porter project set [id]")
	}

	// query for helm repo name
	repoURL, err := utils.PromptPlaintext(fmt.Sprintf(`Provide the image registry URL (include the protocol). For example, https://my-custom-registry.getporter.dev.
Image registry URL: `))
	if err != nil {
		return 0, err
	}

	username, err := utils.PromptPlaintext(fmt.Sprintf(`Provide the username/password for authentication (press enter if no authenicaiton is required).
Username: `))
	if err != nil {
		return 0, err
	}

	password, err := utils.PromptPasswordWithConfirmation()
	if err != nil {
		return 0, err
	}

	// create the basic auth integration
	integration, err := client.CreateBasicAuthIntegration(
		ctx,
		projectID,
		&types.CreateBasicRequest{
			Username: username,
			Password: password,
		},
	)
	if err != nil {
		return 0, err
	}

	color.New(color.FgGreen).Printf("created basic auth integration with id %d\n", integration.ID)

	reg, err := client.CreateRegistry(
		ctx,
		projectID,
		&types.CreateRegistryRequest{
			URL:                repoURL,
			Name:               repoURL,
			BasicIntegrationID: integration.ID,
		},
	)
	if err != nil {
		return 0, err
	}

	color.New(color.FgGreen).Printf("created private registry with id %d and name %s\n", reg.ID, reg.Name)

	return reg.ID, nil
}
