package config

import (
	"context"
	"encoding/base64"
	"fmt"
	"io/ioutil"
	"net/url"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/docker/cli/cli/config/configfile"
	"github.com/docker/cli/cli/config/types"
	"github.com/fatih/color"
	api "github.com/porter-dev/porter/api/client"
	"github.com/porter-dev/porter/cli/cmd/github"
)

// SetDockerConfig sets up the docker config.json
func SetDockerConfig(ctx context.Context, client api.Client, pID uint) error {
	// get all registries that should be added
	regToAdd := make([]string, 0)

	// get the list of namespaces
	resp, err := client.ListRegistries(
		ctx,
		pID,
	)
	if err != nil {
		return err
	}

	registries := *resp

	for _, registry := range registries {
		if registry.URL != "" {
			rURL := registry.URL

			if !strings.Contains(rURL, "http") {
				rURL = "http://" + rURL
			}

			// strip the protocol
			regURL, err := url.Parse(rURL)
			if err != nil {
				continue
			}

			regToAdd = append(regToAdd, regURL.Host)
		}
	}

	// create a docker dir if it does not exist
	dockerDir := filepath.Join(home, ".docker")

	if _, err := os.Stat(dockerDir); os.IsNotExist(err) {
		err = os.Mkdir(dockerDir, 0o700)

		if err != nil {
			return err
		}
	}

	dockerConfigFile := filepath.Join(home, ".docker", "config.json")

	// determine if configfile exists
	if _, err := os.Stat(dockerConfigFile); os.IsNotExist(err) {
		// if it does not exist, create it
		err := ioutil.WriteFile(dockerConfigFile, []byte("{}"), 0o700)
		if err != nil {
			return err
		}
	}

	// read the file bytes
	// // TODO: STEFAN - figure out why we are parsing the ~/.docker/config.json into the CLI config. Are we using the variables somewhere?
	// configBytes, err := ioutil.ReadFile(dockerConfigFile)
	// if err != nil {
	// 	return err
	// }

	// check if the docker credential helper exists
	if !commandExists("docker-credential-porter") {
		err := downloadCredMatchingRelease(ctx)
		if err != nil {
			color.New(color.FgRed).Println("Failed to download credential helper binary:", err.Error())
			os.Exit(1)
		}
	}

	// otherwise, check the version flag of the binary
	cmdVersionCred := exec.Command("docker-credential-porter", "--version")
	writer := &VersionWriter{}
	cmdVersionCred.Stdout = writer

	err = cmdVersionCred.Run()

	if err != nil || writer.Version != Version {
		err := downloadCredMatchingRelease(ctx)
		if err != nil {
			color.New(color.FgRed).Println("Failed to download credential helper binary:", err.Error())
			os.Exit(1)
		}
	}

	configFile := &configfile.ConfigFile{
		Filename: dockerConfigFile,
	}

	// // TODO: STEFAN - figure out why we are parsing the ~/.docker/config.json into the CLI config. Are we using the variables somewhere?
	// err = json.Unmarshal(configBytes, GetCLIConfig())
	// if err != nil {
	// 	return err
	// }

	if configFile.CredentialHelpers == nil {
		configFile.CredentialHelpers = make(map[string]string)
	}

	if configFile.AuthConfigs == nil {
		configFile.AuthConfigs = make(map[string]types.AuthConfig)
	}

	for _, regURL := range regToAdd {
		// if this is a dockerhub registry, see if an auth config has already been generated
		// for index.docker.io
		if strings.Contains(regURL, "index.docker.io") {
			isAuthenticated := false

			for key := range configFile.AuthConfigs {
				if key == "https://index.docker.io/v1/" {
					isAuthenticated = true
				}
			}

			if !isAuthenticated {
				// get a dockerhub token from the Porter API
				tokenResp, err := client.GetDockerhubAuthorizationToken(ctx, pID)
				if err != nil {
					return err
				}

				decodedToken, err := base64.StdEncoding.DecodeString(tokenResp.Token)
				if err != nil {
					return fmt.Errorf("Invalid token: %v", err)
				}

				parts := strings.SplitN(string(decodedToken), ":", 2)

				if len(parts) < 2 {
					return fmt.Errorf("Invalid token: expected two parts, got %d", len(parts))
				}

				configFile.AuthConfigs["https://index.docker.io/v1/"] = types.AuthConfig{
					Auth:     tokenResp.Token,
					Username: parts[0],
					Password: parts[1],
				}

				// since we're using token-based auth, unset the credstore
				configFile.CredentialsStore = ""
			}
		} else {
			configFile.CredentialHelpers[regURL] = "porter"
		}
	}

	return configFile.Save()
}

func commandExists(cmd string) bool {
	_, err := exec.LookPath(cmd)
	return err == nil
}

func downloadCredMatchingRelease(ctx context.Context) error {
	// download the porter cred helper
	z := &github.ZIPReleaseGetter{
		AssetName:           "docker-credential-porter",
		AssetFolderDest:     "/usr/local/bin",
		ZipFolderDest:       filepath.Join(home, ".porter"),
		ZipName:             "docker-credential-porter_latest.zip",
		EntityID:            "porter-dev",
		RepoName:            "porter",
		IsPlatformDependent: true,
		Downloader: &github.ZIPDownloader{
			ZipFolderDest:   filepath.Join(home, ".porter"),
			AssetFolderDest: "/usr/local/bin",
			ZipName:         "docker-credential-porter_latest.zip",
		},
	}

	return z.GetRelease(ctx, Version)
}
