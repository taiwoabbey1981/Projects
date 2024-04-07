package types

type StreamHelmReleaseRequest struct {
	Selectors string   `schema:"selectors"`
	Charts    []string `schema:"charts"`
	Namespace string   `schema:"namespace"`
}
