package preview

import (
	"context"
	"fmt"
	"os"
	"strings"

	"github.com/cli/cli/git"
	"github.com/fatih/color"
	"github.com/mitchellh/mapstructure"
	api "github.com/porter-dev/porter/api/client"
	"github.com/porter-dev/porter/cli/cmd/config"
	"github.com/porter-dev/porter/cli/cmd/deploy"
	"github.com/porter-dev/porter/cli/cmd/deploy/wait"
	cliUtils "github.com/porter-dev/porter/cli/cmd/utils"
	"github.com/porter-dev/porter/internal/integrations/preview"
	"github.com/porter-dev/porter/internal/templater/utils"
	"github.com/porter-dev/switchboard/pkg/drivers"
	"github.com/porter-dev/switchboard/pkg/models"
)

type UpdateConfigDriver struct {
	source      *preview.Source
	target      *preview.Target
	config      *preview.UpdateConfigDriverConfig
	lookupTable *map[string]drivers.Driver
	output      map[string]interface{}
	apiClient   api.Client
	cliConfig   config.CLIConfig
}

// NewUpdateConfigDriver extends switchboard with config updating for an app
func NewUpdateConfigDriver(ctx context.Context, apiClient api.Client, cliConfig config.CLIConfig) func(resource *models.Resource, opts *drivers.SharedDriverOpts) (drivers.Driver, error) {
	return func(resource *models.Resource, opts *drivers.SharedDriverOpts) (drivers.Driver, error) {
		driver := &UpdateConfigDriver{
			lookupTable: opts.DriverLookupTable,
			output:      make(map[string]interface{}),
			apiClient:   apiClient,
			cliConfig:   cliConfig,
		}

		target, err := GetTarget(ctx, resource.Name, resource.Target, apiClient, cliConfig)
		if err != nil {
			return nil, err
		}

		driver.target = target

		source, err := GetSource(ctx, driver.target.Project, resource.Name, resource.Source, apiClient)
		if err != nil {
			return nil, err
		}

		driver.source = source

		return driver, nil
	}
}

func (d *UpdateConfigDriver) ShouldApply(resource *models.Resource) bool {
	return true
}

func (d *UpdateConfigDriver) Apply(resource *models.Resource) (*models.Resource, error) {
	ctx := context.TODO() // switchboard blocks changing this for now

	updateConfigDriverConfig, err := d.getConfig(resource)
	if err != nil {
		return nil, err
	}

	d.config = updateConfigDriverConfig

	_, err = d.apiClient.GetRelease(
		ctx,
		d.target.Project,
		d.target.Cluster,
		d.target.Namespace,
		d.target.AppName,
	)

	shouldCreate := err != nil

	color.New(color.FgBlue).Println("checking for the existence of PORTER_TAG environment variable for the image tag")

	tag := os.Getenv("PORTER_TAG")

	if tag == "" {
		color.New(color.FgBlue).Println("PORTER_TAG environment variable not found, checking for update_config.tag in porter.yaml for the image tag")

		tag = d.config.UpdateConfig.Tag
	}

	if tag == "" {
		color.New(color.FgBlue).Println("update_config.tag not found in porter.yaml, falling back to the latest git commit SHA as the image tag")

		commit, err := git.LastCommit()
		if err != nil {
			return nil, err
		}

		tag = commit.Sha[:7]
	}

	regList, err := d.apiClient.ListRegistries(ctx, d.target.Project)
	if err != nil {
		return nil, err
	}

	var registryURL string

	if len(*regList) == 0 {
		return nil, fmt.Errorf("no registry found")
	} else {
		registryURL = (*regList)[0].URL
	}

	var repoSuffix string

	if repoName := os.Getenv("PORTER_REPO_NAME"); repoName != "" {
		if repoOwner := os.Getenv("PORTER_REPO_OWNER"); repoOwner != "" {
			repoSuffix = cliUtils.SlugifyRepoSuffix(repoOwner, repoName)
		}
	}

	sharedOpts := &deploy.SharedOpts{
		ProjectID:   d.target.Project,
		ClusterID:   d.target.Cluster,
		OverrideTag: tag,
		Namespace:   d.target.Namespace,
		Method:      "registry",
		EnvGroups:   d.config.EnvGroups,
	}

	if shouldCreate {
		color.New(color.FgYellow).Printf("Could not read release %s/%s: attempting creation\n", d.target.Namespace, d.target.AppName)

		createAgent := &deploy.CreateAgent{
			Client: d.apiClient,
			CreateOpts: &deploy.CreateOpts{
				SharedOpts:  sharedOpts,
				Kind:        d.source.Name,
				ReleaseName: d.target.AppName,
				RegistryURL: registryURL,
				RepoSuffix:  repoSuffix,
			},
		}

		image := fmt.Sprintf("%s:%s", strings.Split(d.config.UpdateConfig.Image, ":")[0], tag)

		_, err = createAgent.CreateFromRegistry(ctx, image, d.config.Values)

		if err != nil {
			return nil, err
		}
	} else if !updateConfigDriverConfig.OnlyCreate {
		updateAgent, err := deploy.NewDeployAgent(ctx, d.apiClient, d.target.AppName, &deploy.DeployOpts{
			SharedOpts: sharedOpts,
			Local:      false,
		})
		if err != nil {
			return nil, err
		}

		err = updateAgent.UpdateImageAndValues(ctx, d.config.Values)

		if err != nil {
			return nil, err
		}
	}

	if d.source.Name == "job" && updateConfigDriverConfig.WaitForJob && (shouldCreate || !updateConfigDriverConfig.OnlyCreate) {
		color.New(color.FgYellow).Printf("Waiting for job '%s' to finish\n", resource.Name)

		err = wait.WaitForJob(ctx, d.apiClient, &wait.WaitOpts{
			ProjectID: d.target.Project,
			ClusterID: d.target.Cluster,
			Namespace: d.target.Namespace,
			Name:      d.target.AppName,
		})
		if err != nil {
			return nil, err
		}
	}

	err = d.assignOutput(ctx, resource, d.apiClient)
	if err != nil {
		return nil, err
	}

	return resource, nil
}

func (d *UpdateConfigDriver) Output() (map[string]interface{}, error) {
	return d.output, nil
}

func (d *UpdateConfigDriver) getConfig(resource *models.Resource) (*preview.UpdateConfigDriverConfig, error) {
	populatedConf, err := drivers.ConstructConfig(&drivers.ConstructConfigOpts{
		RawConf:      resource.Config,
		LookupTable:  *d.lookupTable,
		Dependencies: resource.Dependencies,
	})
	if err != nil {
		return nil, err
	}

	config := &preview.UpdateConfigDriverConfig{}

	err = mapstructure.Decode(populatedConf, config)

	if err != nil {
		return nil, err
	}

	return config, nil
}

func (d *UpdateConfigDriver) assignOutput(ctx context.Context, _ *models.Resource, client api.Client) error {
	release, err := client.GetRelease(
		ctx,
		d.target.Project,
		d.target.Cluster,
		d.target.Namespace,
		d.target.AppName,
	)
	if err != nil {
		return err
	}

	d.output = utils.CoalesceValues(d.source.SourceValues, release.Config)

	return nil
}
