package features

import (
	"strconv"
	"strings"

	"github.com/porter-dev/porter/api/server/handlers/cluster"
	"github.com/porter-dev/porter/internal/kubernetes"
)

// isPorterAgentUpdated checks if the agent version is at least the version specified by the major, minor, and patch arguments
func isPorterAgentUpdated(agent *kubernetes.Agent, major, minor, patch int) bool {
	res, err := cluster.GetAgentVersionResponse(agent)
	if err != nil {
		return false
	}
	image := res.Image
	parsed := strings.Split(image, ":")

	if len(parsed) != 2 {
		return false
	}

	tag := parsed[1]
	if tag == "dev" {
		return true
	}

	if !strings.HasPrefix(tag, "v") {
		return false
	}

	tag = strings.TrimPrefix(tag, "v")
	parsedTag := strings.Split(tag, ".")
	if len(parsedTag) != 3 {
		return false
	}

	parsedMajor, _ := strconv.Atoi(parsedTag[0])
	parsedMinor, _ := strconv.Atoi(parsedTag[1])
	parsedPatch, _ := strconv.Atoi(parsedTag[2])
	if parsedMajor < major {
		return false
	} else if parsedMajor > major {
		return true
	}
	if parsedMinor < minor {
		return false
	} else if parsedMinor > minor {
		return true
	}
	return parsedPatch >= patch
}

// Only create the PROGRESSING event if the cluster's agent is updated, because only the updated agent can update the status
func AreAgentDeployEventsEnabled(agent *kubernetes.Agent) bool {
	return isPorterAgentUpdated(agent, 3, 1, 6)
}
