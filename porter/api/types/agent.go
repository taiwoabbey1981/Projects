package types

type DetectAgentResponse struct {
	Version       string `json:"version"`
	LatestVersion string `json:"latest_version"`
	ShouldUpgrade bool   `json:"should_upgrade"`
	Image         string `json:"image"`
}

type GetAgentStatusResponse struct {
	Loki string `json:"loki"`
}
