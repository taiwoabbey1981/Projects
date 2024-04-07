package connect

import (
	"context"
	"fmt"
	"strings"
	"time"

	api "github.com/porter-dev/porter/api/client"
	"github.com/porter-dev/porter/api/types"
	"github.com/porter-dev/porter/cli/cmd/utils"
)

// DOCR creates a DOCR integration
func DOCR(
	ctx context.Context,
	client api.Client,
	projectID uint,
) (uint, error) {
	// if project ID is 0, ask the user to set the project ID or create a project
	if projectID == 0 {
		return 0, fmt.Errorf("no project set, please run porter project set [id]")
	}

	// list oauth integrations and make sure DO exists
	resp, err := client.ListOAuthIntegrations(context.TODO(), projectID)
	if err != nil {
		return 0, err
	}

	oauthInts := *resp

	linkedDO := false
	var doAuth *types.OAuthIntegration

	// iterate through oauth integrations to find do
	for _, oauthInt := range oauthInts {
		if oauthInt.Client == types.OAuthDigitalOcean {
			linkedDO = true
			doAuth = oauthInt
			break
		}
	}

	if !linkedDO {
		doAuth, err = triggerDigitalOceanOAuth(client, projectID)

		if err != nil {
			return 0, err
		}
	}

	// use the digital ocean oauth to create a registry
	regURL, err := utils.PromptPlaintext(fmt.Sprintf(`Please provide the registry URL, in the form registry.digitalocean.com/[REGISTRY_NAME]. For example, registry.digitalocean.com/porter-test. 
Registry URL: `))
	if err != nil {
		return 0, err
	}

	urlArr := strings.Split(regURL, "/")
	regName := urlArr[len(urlArr)-1]

	if err != nil {
		return 0, err
	}

	reg, err := client.CreateRegistry(
		ctx,
		projectID,
		&types.CreateRegistryRequest{
			Name:            regName,
			DOIntegrationID: doAuth.ID,
			URL:             regURL,
		},
	)

	return reg.ID, nil
}

func triggerDigitalOceanOAuth(client api.Client, projectID uint) (*types.OAuthIntegration, error) {
	var doAuth *types.OAuthIntegration

	oauthURL := fmt.Sprintf("%s/projects/%d/oauth/digitalocean", client.BaseURL, projectID)

	fmt.Printf("Please visit %s in your browser to connect to Digital Ocean (it should open automatically).", oauthURL)
	utils.OpenBrowser(oauthURL)

	for {
		resp, err := client.ListOAuthIntegrations(context.TODO(), projectID)
		if err != nil {
			return doAuth, err
		}

		linkedDO := false

		oauthInts := *resp

		// iterate through oauth integrations to find do
		for _, oauthInt := range oauthInts {
			if oauthInt.Client == types.OAuthDigitalOcean {
				linkedDO = true
				doAuth = oauthInt
				break
			}
		}

		if linkedDO {
			break
		}

		time.Sleep(2 * time.Second)
	}

	return doAuth, nil
}
