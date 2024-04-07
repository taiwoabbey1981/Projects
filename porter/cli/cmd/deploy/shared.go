package deploy

import (
	"context"
	"fmt"

	api "github.com/porter-dev/porter/api/client"
	"github.com/porter-dev/porter/api/types"
)

// SharedOpts are common options for build, create, and deploy agents
type SharedOpts struct {
	ProjectID       uint
	ClusterID       uint
	Namespace       string
	LocalPath       string
	LocalDockerfile string
	OverrideTag     string
	Method          DeployBuildType
	AdditionalEnv   map[string]string
	EnvGroups       []types.EnvGroupMeta
	UseCache        bool
}

func coalesceEnvGroups(
	ctx context.Context,
	client api.Client,
	projectID, clusterID uint,
	namespace string,
	envGroups []types.EnvGroupMeta,
	config map[string]interface{},
) error {
	for _, group := range envGroups {
		if group.Name == "" {
			return fmt.Errorf("env group name cannot be empty")
		}

		envGroup, err := client.GetEnvGroup(
			ctx,
			projectID,
			clusterID,
			namespace,
			&types.GetEnvGroupRequest{
				Name:    group.Name,
				Version: group.Version,
			},
		)
		if err != nil {
			return err
		}

		envConfig, err := GetNestedMap(config, "container", "env", "normal")

		if err != nil || envConfig == nil {
			envConfig = make(map[string]interface{})
		}

		for k, v := range envGroup.Variables {
			envConfig[k] = v
		}

		containerMap, _ := config["container"].(map[string]interface{})
		envMap, _ := containerMap["env"].(map[string]interface{})

		envMap["normal"] = envConfig
	}

	return nil
}
