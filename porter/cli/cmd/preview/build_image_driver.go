package preview

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/cli/cli/git"
	"github.com/docker/distribution/reference"
	"github.com/mitchellh/mapstructure"
	"github.com/porter-dev/porter/api/client"
	"github.com/porter-dev/porter/api/types"
	"github.com/porter-dev/porter/cli/cmd/config"
	"github.com/porter-dev/porter/cli/cmd/deploy"
	"github.com/porter-dev/porter/cli/cmd/docker"
	"github.com/porter-dev/porter/cli/cmd/utils"
	"github.com/porter-dev/porter/internal/integrations/preview"
	"github.com/porter-dev/switchboard/pkg/drivers"
	"github.com/porter-dev/switchboard/pkg/models"
)

type BuildDriver struct {
	source      *preview.Source
	target      *preview.Target
	config      *preview.BuildDriverConfig
	lookupTable *map[string]drivers.Driver
	output      map[string]interface{}
	apiClient   client.Client
	cliConfig   config.CLIConfig
}

// NewBuildDriver extends switchboard with the ability to build images and buildpacks
func NewBuildDriver(ctx context.Context, apiClient client.Client, cliConfig config.CLIConfig) func(resource *models.Resource, opts *drivers.SharedDriverOpts) (drivers.Driver, error) {
	return func(resource *models.Resource, opts *drivers.SharedDriverOpts) (drivers.Driver, error) {
		driver := &BuildDriver{
			lookupTable: opts.DriverLookupTable,
			output:      make(map[string]interface{}),
			cliConfig:   cliConfig,
			apiClient:   apiClient,
		}

		target, err := GetTarget(ctx, resource.Name, resource.Target, apiClient, cliConfig)
		if err != nil {
			return nil, err
		}

		driver.target = target

		source, err := GetSource(ctx, target.Project, resource.Name, resource.Source, apiClient)
		if err != nil {
			return nil, err
		}

		driver.source = source

		return driver, nil
	}
}

func (d *BuildDriver) ShouldApply(resource *models.Resource) bool {
	return true
}

func (d *BuildDriver) Apply(resource *models.Resource) (*models.Resource, error) {
	ctx := context.TODO() // switchboard blocks changing this for now

	buildDriverConfig, err := d.getConfig(resource)
	if err != nil {
		return nil, err
	}

	d.config = buildDriverConfig

	// FIXME: give tag option in config build, but override if PORTER_TAG is present
	tag := os.Getenv("PORTER_TAG")

	if tag == "" {
		commit, err := git.LastCommit()
		if err == nil {
			tag = commit.Sha[:7]
		}
	}

	// if the method is registry and a tag is defined, we use the provided tag
	if d.config.Build.Method == "registry" {
		imageSpl := strings.Split(d.config.Build.Image, ":")

		if len(imageSpl) == 2 {
			tag = imageSpl[1]
		}

		if tag == "" {
			tag = "latest"
		}
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
			repoSuffix = utils.SlugifyRepoSuffix(repoOwner, repoName)
		}
	}

	createAgent := &deploy.CreateAgent{
		Client: d.apiClient,
		CreateOpts: &deploy.CreateOpts{
			SharedOpts: &deploy.SharedOpts{
				ProjectID:       d.target.Project,
				ClusterID:       d.target.Cluster,
				OverrideTag:     tag,
				Namespace:       d.target.Namespace,
				LocalPath:       d.config.Build.Context,
				LocalDockerfile: d.config.Build.Dockerfile,
				Method:          deploy.DeployBuildType(d.config.Build.Method),
				EnvGroups:       d.config.EnvGroups,
				UseCache:        d.config.Build.UsePackCache,
			},
			Kind:        d.source.Name,
			ReleaseName: d.target.AppName,
			RegistryURL: registryURL,
			RepoSuffix:  repoSuffix,
		},
	}

	regID, imageURL, err := createAgent.GetImageRepoURL(ctx, d.target.AppName, d.target.Namespace)
	if err != nil {
		return nil, err
	}

	// create repository if it does not exist
	repoResp, err := d.apiClient.ListRegistryRepositories(ctx, d.target.Project, regID)
	if err != nil {
		return nil, err
	}

	repos := *repoResp

	found := false

	for _, repo := range repos {
		if repo.URI == imageURL {
			found = true
			break
		}
	}

	if !found {
		err = d.apiClient.CreateRepository(
			ctx,
			d.target.Project,
			regID,
			&types.CreateRegistryRepositoryRequest{
				ImageRepoURI: imageURL,
			},
		)

		if err != nil {
			return nil, err
		}
	}

	if d.config.Build.UsePackCache {
		err := config.SetDockerConfig(ctx, d.apiClient, d.target.Project)
		if err != nil {
			return nil, err
		}
	}

	if d.config.Build.Method != "" {
		if d.config.Build.Method == string(deploy.DeployBuildTypeDocker) {
			if d.config.Build.Dockerfile == "" {
				hasDockerfile := createAgent.HasDefaultDockerfile(d.config.Build.Context)

				if !hasDockerfile {
					return nil, fmt.Errorf("dockerfile not found")
				}

				d.config.Build.Dockerfile = "Dockerfile"
			}
		}
	} else {
		// try to detect dockerfile, otherwise fall back to `pack`
		hasDockerfile := createAgent.HasDefaultDockerfile(d.config.Build.Context)

		if !hasDockerfile {
			d.config.Build.Method = string(deploy.DeployBuildTypePack)
		} else {
			d.config.Build.Method = string(deploy.DeployBuildTypeDocker)
			d.config.Build.Dockerfile = "Dockerfile"
		}
	}

	// create docker agent
	agent, err := docker.NewAgentWithAuthGetter(ctx, d.apiClient, d.target.Project)
	if err != nil {
		return nil, err
	}

	_, mergedValues, err := createAgent.GetMergedValues(ctx, d.config.Values)
	if err != nil {
		return nil, err
	}

	env, err := deploy.GetEnvForRelease(
		ctx,
		d.apiClient,
		mergedValues,
		d.target.Project,
		d.target.Cluster,
		d.target.Namespace,
	)
	if err != nil {
		env = make(map[string]string)
	}

	envConfig, err := deploy.GetNestedMap(mergedValues, "container", "env")

	if err == nil {
		_, exists := envConfig["build"]

		if exists {
			buildEnv, err := deploy.GetNestedMap(mergedValues, "container", "env", "build")

			if err == nil {
				for key, val := range buildEnv {
					if valStr, ok := val.(string); ok {
						env[key] = valStr
					}
				}
			}
		}
	}

	for k, v := range d.config.Build.Env {
		env[k] = v
	}

	buildAgent := &deploy.BuildAgent{
		SharedOpts:  createAgent.CreateOpts.SharedOpts,
		APIClient:   d.apiClient,
		ImageRepo:   imageURL,
		Env:         env,
		ImageExists: false,
	}

	if d.config.Build.Method == string(deploy.DeployBuildTypeDocker) {
		var basePath string

		basePath, err = filepath.Abs(".")

		if err != nil {
			return nil, err
		}

		var currentTag string
		// implement caching for porter stack builds
		if os.Getenv("PORTER_STACK_NAME") != "" {
			currentTag = getCurrentImageTagIfExists(ctx, d.apiClient, d.target.Project, d.target.Cluster, os.Getenv("PORTER_STACK_NAME"))
		}

		err = buildAgent.BuildDocker(
			ctx,
			agent,
			basePath,
			d.config.Build.Context,
			d.config.Build.Dockerfile,
			tag,
			currentTag,
		)
	} else {
		var buildConfig *types.BuildConfig

		if d.config.Build.Builder != "" {
			buildConfig = &types.BuildConfig{
				Builder:    d.config.Build.Builder,
				Buildpacks: d.config.Build.Buildpacks,
			}
		}

		err = buildAgent.BuildPack(
			ctx,
			agent,
			d.config.Build.Context,
			tag,
			"",
			buildConfig,
		)
	}

	if err != nil {
		return nil, err
	}

	named, _ := reference.ParseNamed(imageURL)
	domain := reference.Domain(named)
	imageRepo := reference.Path(named)

	d.output["registry_url"] = domain
	d.output["image_repo"] = imageRepo
	d.output["image_tag"] = tag
	d.output["image"] = fmt.Sprintf("%s:%s", imageURL, tag)

	return resource, nil
}

func getCurrentImageTagIfExists(ctx context.Context, client client.Client, projectID, clusterID uint, stackName string) string {
	namespace := fmt.Sprintf("porter-stack-%s", stackName)
	release, err := client.GetRelease(
		ctx,
		projectID,
		clusterID,
		namespace,
		stackName,
	)
	if err != nil {
		return ""
	}

	if release == nil {
		return ""
	}

	if release.Config == nil {
		return ""
	}

	value, ok := release.Config["global"]
	if !ok {
		return ""
	}
	globalConfig := value.(map[string]interface{})
	if globalConfig == nil {
		return ""
	}

	value, ok = globalConfig["image"]
	if !ok {
		return ""
	}
	imageConfig := value.(map[string]interface{})
	if imageConfig == nil {
		return ""
	}

	value, ok = imageConfig["tag"]
	if !ok {
		return ""
	}
	tag := value.(string)

	return tag
}

func (d *BuildDriver) Output() (map[string]interface{}, error) {
	return d.output, nil
}

func (d *BuildDriver) getConfig(resource *models.Resource) (*preview.BuildDriverConfig, error) {
	populatedConf, err := drivers.ConstructConfig(&drivers.ConstructConfigOpts{
		RawConf:      resource.Config,
		LookupTable:  *d.lookupTable,
		Dependencies: resource.Dependencies,
	})
	if err != nil {
		return nil, err
	}

	config := &preview.BuildDriverConfig{}

	err = mapstructure.Decode(populatedConf, config)

	if err != nil {
		return nil, err
	}

	return config, nil
}
