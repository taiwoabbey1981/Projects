package v2beta1

import (
	"context"
	"fmt"
	"regexp"

	api "github.com/porter-dev/porter/api/client"
	"github.com/porter-dev/porter/cli/cmd/config"
	"github.com/porter-dev/switchboard/pkg/types"
	"gopkg.in/yaml.v3"
)

type PreviewApplier struct {
	apiClient api.Client
	cliConfig config.CLIConfig
	rawBytes  []byte
	namespace string
	parsed    *PorterYAML
}

// NewApplier returns an applier for preview environments
func NewApplier(client api.Client, cliConfig config.CLIConfig, raw []byte, namespace string) (*PreviewApplier, error) {
	// replace all instances of ${{ porter.env.FOO }} with { .get-env.FOO }
	re := regexp.MustCompile(`\$\{\{\s*porter\.env\.(.*)\s*\}\}`)
	raw = re.ReplaceAll(raw, []byte("{.get-env.$1}"))

	parsed := &PorterYAML{}

	err := yaml.Unmarshal(raw, parsed)
	if err != nil {
		errMsg := composePreviewMessage("error parsing porter.yaml", Error)
		return nil, fmt.Errorf("%s: %w", errMsg, err)
	}

	err = cliConfig.ValidateCLIEnvironment()

	if err != nil {
		errMsg := composePreviewMessage("porter CLI is not configured correctly", Error)
		return nil, fmt.Errorf("%s: %w", errMsg, err)
	}

	return &PreviewApplier{
		apiClient: client,
		cliConfig: cliConfig,
		rawBytes:  raw,
		namespace: namespace,
		parsed:    parsed,
	}, nil
}

func (a *PreviewApplier) Apply() error {
	// for v2beta1, check if the namespace exists in the current project-cluster pair
	//
	// this is a sanity check to ensure that the user does not see any internal
	// errors that are caused by the namespace not existing
	nsList, err := a.apiClient.GetK8sNamespaces(
		context.TODO(), // can not change because of switchboard
		a.cliConfig.Project,
		a.cliConfig.Cluster,
	)
	if err != nil {
		errMsg := composePreviewMessage(fmt.Sprintf("error listing namespaces for project '%d', cluster '%d'",
			a.cliConfig.Project, a.cliConfig.Cluster), Error)
		return fmt.Errorf("%s: %w", errMsg, err)
	}

	namespaces := *nsList
	nsFound := false

	for _, ns := range namespaces {
		if ns.Name == a.namespace {
			nsFound = true
			break
		}
	}

	if !nsFound {
		// 	errMsg := composePreviewMessage(fmt.Sprintf("namespace '%s' does not exist in project '%d', cluster '%d'",
		// 		a.namespace, config.GetCLIConfig().Project, config.GetCLIConfig().Cluster), Error)
		// 	return fmt.Errorf("%s: %w", errMsg, err)
	}

	printInfoMessage(fmt.Sprintf("Applying porter.yaml with the following attributes:\n"+
		"\tHost: %s\n\tProject ID: %d\n\tCluster ID: %d\n\tNamespace: %s",
		a.cliConfig.Host,
		a.cliConfig.Project,
		a.cliConfig.Cluster,
		a.namespace),
	)

	// err = a.readOSEnv()

	// if err != nil {
	// 	errMsg := composePreviewMessage("error reading OS environment variables", Error)
	// 	return fmt.Errorf("%s: %w", errMsg, err)
	// }

	// err = a.processVariables()

	// if err != nil {
	// 	return err
	// }

	// err = a.processEnvGroups()

	// if err != nil {
	// 	return err
	// }

	return nil
}

func (a *PreviewApplier) DowngradeToV1() (*types.ResourceGroup, error) {
	err := a.Apply()
	if err != nil {
		return nil, err
	}

	v1File := &types.ResourceGroup{
		Version: "v1",
		Resources: []*types.Resource{
			{
				Name:   "get-env",
				Driver: "os-env",
			},
		},
	}

	buildRefs := make(map[string]*Build)

	for _, b := range a.parsed.Builds {
		if b == nil {
			continue
		}

		buildRefs[b.GetName()] = b

		bi, err := b.getV1BuildImage()
		if err != nil {
			return nil, err
		}

		pi, err := b.getV1PushImage()
		if err != nil {
			return nil, err
		}

		v1File.Resources = append(v1File.Resources, bi, pi)
	}

	for _, app := range a.parsed.Apps {
		if app == nil {
			continue
		}

		if _, ok := buildRefs[app.GetBuildRef()]; !ok {
			errMsg := composePreviewMessage(fmt.Sprintf("build_ref '%s' referenced by app '%s' does not exist",
				app.GetBuildRef(), app.GetName()), Error)
			return nil, fmt.Errorf("%s: %w", errMsg, err)
		}

		ai, err := app.getV1Resource(buildRefs[app.GetBuildRef()])
		if err != nil {
			return nil, err
		}

		v1File.Resources = append(v1File.Resources, ai)
	}

	for _, addon := range a.parsed.Addons {
		if addon == nil {
			continue
		}

		ai, err := addon.getV1Addon()
		if err != nil {
			return nil, err
		}

		v1File.Resources = append(v1File.Resources, ai)
	}

	return v1File, nil
}

// func (a *PreviewApplier) readOSEnv() error {
// 	printInfoMessage("Reading OS environment variables")

// 	env := os.Environ()
// 	osEnv := make(map[string]string)

// 	for _, e := range env {
// 		k, v, _ := strings.Cut(e, "=")
// 		kCopy := k

// 		if k != "" && v != "" && strings.HasPrefix(k, "PORTER_APPLY_") {
// 			// we only read in env variables that start with PORTER_APPLY_
// 			for strings.HasPrefix(k, "PORTER_APPLY_") {
// 				k = strings.TrimPrefix(k, "PORTER_APPLY_")
// 			}

// 			if k == "" {
// 				printWarningMessage(fmt.Sprintf("Ignoring invalid OS environment variable '%s'", kCopy))
// 			}

// 			osEnv[k] = v
// 		}
// 	}

// 	a.osEnv = osEnv

// 	return nil
// }

// func (a *PreviewApplier) processVariables() error {
// 	printInfoMessage("Processing variables")

// 	constantsMap := make(map[string]string)
// 	variablesMap := make(map[string]string)

// 	for _, v := range a.parsed.Variables {
// 		if v == nil {
// 			continue
// 		}

// 		if v.Once != nil && *v.Once {
// 			// a constant which should be stored in the env group on first run
// 			if exists, err := a.constantExistsInEnvGroup(*v.Name); err == nil {
// 				if exists == nil {
// 					// this should not happen
// 					return fmt.Errorf("internal error: please let the Porter team know about this and quote the following " +
// 						"error:\n-----\nERROR: checking for constant existence in env group returned nil with no error")
// 				}

// 				val := *exists

// 				if !val {
// 					// create the constant in the env group
// 					if *v.Value != "" {
// 						constantsMap[*v.Name] = *v.Value
// 					} else if v.Random != nil && *v.Random {
// 						constantsMap[*v.Name] = randomString(*v.Length, defaultCharset)
// 					} else {
// 						// this should not happen
// 						return fmt.Errorf("internal error: please let the Porter team know about this and quote the following "+
// 							"error:\n-----\nERROR: for variable '%s', random is false and value is empty", *v.Name)
// 					}
// 				}
// 			} else {
// 				return fmt.Errorf("error checking for existence of constant %s: %w", *v.Name, err)
// 			}
// 		} else {
// 			if v.Value != nil && *v.Value != "" {
// 				variablesMap[*v.Name] = *v.Value
// 			} else if v.Random != nil && *v.Random {
// 				variablesMap[*v.Name] = randomString(*v.Length, defaultCharset)
// 			} else {
// 				// this should not happen
// 				return fmt.Errorf("internal error: please let the Porter team know about this and quote the following "+
// 					"error:\n-----\nERROR: for variable '%s', random is false and value is empty", *v.Name)
// 			}
// 		}
// 	}

// 	if len(constantsMap) > 0 {
// 		// we need to create these constants in the env group
// 		_, err := a.apiClient.CreateEnvGroup(
// 			ctx,
// 			config.GetCLIConfig().Project,
// 			config.GetCLIConfig().Cluster,
// 			a.namespace,
// 			&apiTypes.CreateEnvGroupRequest{
// 				Name:      constantsEnvGroup,
// 				Variables: constantsMap,
// 			},
// 		)

// 		if err != nil {
// 			return fmt.Errorf("error creating constants (variables with once set to true) in env group: %w", err)
// 		}

// 		for k, v := range constantsMap {
// 			variablesMap[k] = v
// 		}
// 	}

// 	a.variablesMap = variablesMap

// 	return nil
// }

// func (a *PreviewApplier) constantExistsInEnvGroup(name string) (*bool, error) {
// 	apiResponse, err := a.apiClient.GetEnvGroup(
// 		ctx,
// 		config.GetCLIConfig().Project,
// 		config.GetCLIConfig().Cluster,
// 		a.namespace,
// 		&apiTypes.GetEnvGroupRequest{
// 			Name: constantsEnvGroup,
// 			// we do not care about the version because it always needs to be the latest
// 		},
// 	)

// 	if err != nil {
// 		if strings.Contains(err.Error(), "env group not found") {
// 			return booleanptr(false), nil
// 		}

// 		return nil, err
// 	}

// 	if _, ok := apiResponse.Variables[name]; ok {
// 		return booleanptr(true), nil
// 	}

// 	return booleanptr(false), nil
// }

// func (a *PreviewApplier) processEnvGroups() error {
// 	printInfoMessage("Processing env groups")

// 	for _, eg := range a.parsed.EnvGroups {
// 		if eg == nil {
// 			continue
// 		}

// 		if eg.Name == nil || *eg.Name == "" {

// 		}

// 		envGroup, err := a.apiClient.GetEnvGroup(
// 			ctx,
// 			config.GetCLIConfig().Project,
// 			config.GetCLIConfig().Cluster,
// 			a.namespace,
// 			&apiTypes.GetEnvGroupRequest{
// 				Name: *eg.Name,
// 			},
// 		)

// 		if err != nil && strings.Contains(err.Error(), "env group not found") {
// 			if eg.CloneFrom == nil {
// 				return fmt.Errorf(composePreviewMessage(fmt.Sprintf("empty clone_from for env group '%s'", *eg.Name), Error))
// 			}

// 			egNS, egName, found := strings.Cut(*eg.CloneFrom, "/")

// 			if !found {
// 				return fmt.Errorf("error parsing clone_from for env group '%s': invalid format", *eg.Name)
// 			}

// 			// clone the env group
// 			envGroup, err := a.apiClient.CloneEnvGroup(
// 				ctx,
// 				config.GetCLIConfig().Project,
// 				config.GetCLIConfig().Cluster,
// 				egNS,
// 				&apiTypes.CloneEnvGroupRequest{
// 					SourceName:      egName,
// 					TargetNamespace: a.namespace,
// 					TargetName:      *eg.Name,
// 				},
// 			)

// 			if err != nil {
// 				return fmt.Errorf("error cloning env group '%s' from '%s': %w", egName, egNS, err)
// 			}

// 			a.envGroups[*eg.Name] = &apiTypes.EnvGroup{
// 				Name:      envGroup.Name,
// 				Variables: envGroup.Variables,
// 			}
// 		} else if err != nil {
// 			return fmt.Errorf("error checking for env group '%s': %w", *eg.Name, err)
// 		} else {
// 			a.envGroups[*eg.Name] = &apiTypes.EnvGroup{
// 				Name:      envGroup.Name,
// 				Variables: envGroup.Variables,
// 			}
// 		}
// 	}

// 	return nil
// }
