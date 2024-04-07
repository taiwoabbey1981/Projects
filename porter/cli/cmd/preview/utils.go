package preview

import (
	"context"
	"fmt"
	"os"
	"strconv"

	api "github.com/porter-dev/porter/api/client"
	"github.com/porter-dev/porter/api/types"
	"github.com/porter-dev/porter/cli/cmd/config"
	"github.com/porter-dev/porter/internal/integrations/preview"
)

// GetSource extends switchboard
func GetSource(ctx context.Context, projectID uint, resourceName string, input map[string]interface{}, apiClient api.Client) (*preview.Source, error) {
	output := &preview.Source{}

	// first read from env vars
	output.Name = os.Getenv("PORTER_SOURCE_NAME")
	output.Repo = os.Getenv("PORTER_SOURCE_REPO")
	output.Version = os.Getenv("PORTER_SOURCE_VERSION")

	// next, check for values in the YAML file
	if output.Name == "" {
		if name, ok := input["name"]; ok {
			nameVal, ok := name.(string)
			if !ok {
				return nil, fmt.Errorf("error parsing source for resource '%s': invalid name provided", resourceName)
			}
			output.Name = nameVal
		}
	}

	if output.Name == "" {
		return nil, fmt.Errorf("error parsing source for resource '%s': source name required", resourceName)
	}

	if output.Repo == "" {
		if repo, ok := input["repo"]; ok {
			repoVal, ok := repo.(string)
			if !ok {
				return nil, fmt.Errorf("error parsing source for resource '%s': invalid repo provided", resourceName)
			}
			output.Repo = repoVal
		}
	}

	if output.Version == "" {
		if version, ok := input["version"]; ok {
			versionVal, ok := version.(string)
			if !ok {
				return nil, fmt.Errorf("error parsing source for resource '%s': invalid version provided", resourceName)
			}
			output.Version = versionVal
		}
	}

	// lastly, just put in the defaults
	if output.Version == "" {
		output.Version = "latest"
	}

	serverMetadata, err := apiClient.GetPorterInstanceMetadata(ctx)
	if err != nil {
		return nil, fmt.Errorf("error fetching Porter instance metadata: %w", err)
	}

	if output.Repo == "" {
		if serverMetadata.DefaultAppHelmRepoURL != "" {
			output.Repo = serverMetadata.DefaultAppHelmRepoURL
		} else {
			output.Repo = "https://charts.getporter.dev"
		}

		values, err := existsInRepo(ctx, projectID, output.Name, output.Version, output.Repo, apiClient)

		if err == nil {
			output.SourceValues = values
			output.IsApplication = true

			return output, nil
		}

		if serverMetadata.DefaultAddonHelmRepoURL != "" {
			output.Repo = serverMetadata.DefaultAddonHelmRepoURL
		} else {
			output.Repo = "https://chart-addons.getporter.dev"
		}

		values, err = existsInRepo(ctx, projectID, output.Name, output.Version, output.Repo, apiClient)

		if err == nil {
			output.SourceValues = values

			return output, nil
		}

		return nil, fmt.Errorf("error parsing source for resource '%s': source chart does not exist in the default "+
			"Helm repositories", resourceName)
	} else {
		// we look in the passed-in repo
		values, err := existsInRepo(ctx, projectID, output.Name, output.Version, output.Repo, apiClient)

		if err == nil {
			output.SourceValues = values
			output.IsApplication = output.Repo == serverMetadata.DefaultAppHelmRepoURL || output.Repo == "https://charts.getporter.dev"

			return output, nil
		}
	}

	return nil, fmt.Errorf("error parsing source for resource '%s': source '%s' does not exist in repo '%s'",
		resourceName, output.Name, output.Repo)
}

// GetTarget extends switchboard
func GetTarget(ctx context.Context, resourceName string, input map[string]interface{}, apiClient api.Client, cliConfig config.CLIConfig) (*preview.Target, error) {
	output := &preview.Target{}

	// first read from env vars
	if projectEnv := os.Getenv("PORTER_PROJECT"); projectEnv != "" {
		project, err := strconv.Atoi(projectEnv)
		if err != nil {
			return nil, fmt.Errorf("error parsing target for resource '%s': %w", resourceName, err)
		}
		output.Project = uint(project)
	}

	if clusterEnv := os.Getenv("PORTER_CLUSTER"); clusterEnv != "" {
		cluster, err := strconv.Atoi(clusterEnv)
		if err != nil {
			return nil, fmt.Errorf("error parsing target for resource '%s': %w", resourceName, err)
		}
		output.Cluster = uint(cluster)
	}

	output.Namespace = os.Getenv("PORTER_NAMESPACE")

	// next, check for values in the YAML file
	if output.Project == 0 {
		if project, ok := input["project"]; ok {
			projectVal, ok := project.(uint)
			if !ok {
				return nil, fmt.Errorf("error parsing target for resource '%s': project value must be an integer", resourceName)
			}
			output.Project = projectVal
		}
	}

	if output.Cluster == 0 {
		if cluster, ok := input["cluster"]; ok {
			clusterVal, ok := cluster.(uint)
			if !ok {
				return nil, fmt.Errorf("error parsing target for resource '%s': cluster value must be an integer",
					resourceName)
			}
			output.Cluster = clusterVal
		}
	}

	if output.Namespace == "" {
		if namespace, ok := input["namespace"]; ok {
			namespaceVal, ok := namespace.(string)
			if !ok {
				return nil, fmt.Errorf("error parsing target for resource '%s': invalid namespace provided", resourceName)
			}
			output.Namespace = namespaceVal
		}
	}

	if registryURL, ok := input["registry_url"]; ok {
		registryURLVal, ok := registryURL.(string)
		if !ok {
			return nil, fmt.Errorf("error parsing target for resource '%s': invalid registry_url provided", resourceName)
		}
		output.RegistryURL = registryURLVal
	}

	if appName, ok := input["app_name"]; ok {
		appNameVal, ok := appName.(string)
		if !ok {
			return nil, fmt.Errorf("error parsing target for resource '%s': invalid app_name provided", resourceName)
		}
		output.AppName = appNameVal
	}

	// lastly, just put in the defaults

	if output.Project == 0 {
		output.Project = cliConfig.Project
	}

	if output.Cluster == 0 {
		output.Cluster = cliConfig.Cluster
	}

	if output.Namespace == "" {
		output.Namespace = "default"
	}

	if output.RegistryURL == "" {
		if cliConfig.Registry == 0 {
			regList, err := apiClient.ListRegistries(ctx, output.Project)
			if err != nil {
				return nil, fmt.Errorf("for resource '%s', error listing registries in project: %w", resourceName, err)
			}

			if len(*regList) == 0 {
				return nil, fmt.Errorf("for resource '%s', no registries found in project", resourceName)
			}

			output.RegistryURL = (*regList)[0].URL
		} else {
			reg, err := apiClient.GetRegistry(ctx, output.Project, cliConfig.Registry)
			if err != nil {
				return nil, fmt.Errorf("for resource '%s', error getting registry from CLI config: %w", resourceName, err)
			}

			output.RegistryURL = reg.URL
		}
	}

	return output, nil
}

func existsInRepo(ctx context.Context, projectID uint, name, version, url string, apiClient api.Client) (map[string]interface{}, error) {
	chart, err := apiClient.GetTemplate(
		ctx,
		projectID,
		name, version,
		&types.GetTemplateRequest{
			TemplateGetBaseRequest: types.TemplateGetBaseRequest{
				RepoURL: url,
			},
		},
	)
	if err != nil {
		return nil, err
	}
	return chart.Values, nil
}
