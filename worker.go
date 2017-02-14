package atc

import "time"

type WorkerState string

const (
	WorkerStateRunning  = WorkerState("running")
	WorkerStateStalled  = WorkerState("stalled")
	WorkerStateLanding  = WorkerState("landing")
	WorkerStateLanded   = WorkerState("landed")
	WorkerStateRetiring = WorkerState("retiring")
)

type Worker struct {
	Name  string      `json:"name"`
	State WorkerState `json:"state"`

	GardenAddr      *string `json:"addr"` // not garden_addr, for backwards-compatibility
	BaggageclaimURL *string `json:"baggageclaim_url"`

	HTTPProxyURL  string `json:"http_proxy_url,omitempty"`
	HTTPSProxyURL string `json:"https_proxy_url,omitempty"`
	NoProxy       string `json:"no_proxy,omitempty"`

	ActiveContainers int `json:"active_containers"`

	ResourceTypes []WorkerResourceType `json:"resource_types"`

	Platform  string   `json:"platform"`
	Tags      []string `json:"tags"`
	Team      string   `json:"team"`
	StartTime int64    `json:"start_time"`

	ExpiresIn time.Duration `json:"expires"`
}

type WorkerResourceType struct {
	Type    string `json:"type"`
	Image   string `json:"image"`
	Version string `json:"version"`
}

type PruneWorkerResponseBody struct {
	Stderr string `json:"stderr"`
}
