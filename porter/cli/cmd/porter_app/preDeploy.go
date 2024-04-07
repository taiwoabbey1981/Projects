package porter_app

import (
	"context"
	"fmt"

	"github.com/mitchellh/mapstructure"
	api "github.com/porter-dev/porter/api/client"
	"github.com/porter-dev/porter/cli/cmd/deploy"
	"github.com/porter-dev/porter/internal/integrations/preview"

	switchboardTypes "github.com/porter-dev/switchboard/pkg/types"
)

func createPreDeployResource(ctx context.Context, client api.Client, release *Service, stackName, buildResourceName, pushResourceName string, projectID, clusterID uint, env map[string]string) (*switchboardTypes.Resource, string, error) {
	var finalCmd string
	if release != nil && release.Run != nil {
		finalCmd = *release.Run
	} else {
		finalCmd = getPredeployStartCommandFromRelease(ctx, client, stackName, projectID, clusterID)
		if finalCmd == "" {
			return nil, "", nil
		}
	}

	config := &preview.ApplicationConfig{}

	config.Build.Method = "registry"
	config.Build.Image = fmt.Sprintf("{ .%s.image }", buildResourceName)
	config.Build.Env = CopyEnv(env)
	config.WaitForJob = true
	config.InjectBuild = true

	helm_values := make(map[string]interface{})
	if release != nil && release.Config != nil {
		helm_values = release.Config
	}
	helm_values["container"] = map[string]interface{}{
		"command": finalCmd,
		"env": map[string]interface{}{
			"normal": CopyEnv(env),
		},
	}
	helm_values["paused"] = false
	config.Values = convertMap(helm_values).(map[string]interface{})

	rawConfig := make(map[string]any)
	err := mapstructure.Decode(config, &rawConfig)
	if err != nil {
		return nil, "", fmt.Errorf("could not decode config for release: %w", err)
	}

	return &switchboardTypes.Resource{
		Name:      fmt.Sprintf("%s-r", stackName),
		DependsOn: []string{"get-env", buildResourceName, pushResourceName},
		Source: map[string]any{
			"name": "job",
		},
		Target: map[string]any{
			"app_name":  fmt.Sprintf("%s-r", stackName),
			"namespace": fmt.Sprintf("porter-stack-%s", stackName),
		},
		Config: rawConfig,
	}, finalCmd, nil
}

func getPredeployStartCommandFromRelease(ctx context.Context, client api.Client, stackName string, projectID uint, clusterID uint) string {
	namespace := fmt.Sprintf("porter-stack-%s", stackName)
	releaseName := fmt.Sprintf("%s-r", stackName)
	release, err := client.GetRelease(
		ctx,
		projectID,
		clusterID,
		namespace,
		releaseName,
	)

	if err != nil || release == nil || release.Config == nil {
		return ""
	}

	containerMap, err := deploy.GetNestedMap(release.Config, "container")
	if err == nil {
		if command, ok := containerMap["command"]; ok {
			if commandString, ok := command.(string); ok {
				return commandString
			}
		}
	}

	return ""
}

func convertMap(m interface{}) interface{} {
	switch m := m.(type) {
	case map[string]interface{}:
		for k, v := range m {
			m[k] = convertMap(v)
		}
	case map[interface{}]interface{}:
		result := map[string]interface{}{}
		for k, v := range m {
			result[k.(string)] = convertMap(v)
		}
		return result
	case []interface{}:
		for i, v := range m {
			m[i] = convertMap(v)
		}
	}
	return m
}
