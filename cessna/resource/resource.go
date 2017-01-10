package resource

import (
	"github.com/concourse/atc"
	"github.com/concourse/baggageclaim"
)

type ResourceType struct {
	RootFSPath string
	Name       string
}

type Resource struct {
	ResourceType ResourceType
	Source       atc.Source
}

func NewResource(resourceType ResourceType, source atc.Source) Resource {
	return Resource{
		ResourceType: resourceType,
		Source:       source,
	}
}

type CheckRequest struct {
	Source  atc.Source  `json:"source"`
	Version atc.Version `json:"version"`
}

type CheckResponse []atc.Version

type InRequest struct {
	Source  atc.Source  `json:"source"`
	Params  atc.Params  `json:"params"`
	Version atc.Version `json:"version"`
}

type InResponse struct {
	Version  atc.Version         `json:"version"`
	Metadata []atc.MetadataField `json:"metadata,omitempty"`
}

type OutRequest struct {
	Source atc.Source `json:"source"`
	Params atc.Params `json:"params"`
}

type OutResponse struct {
	Version  atc.Version         `json:"version"`
	Metadata []atc.MetadataField `json:"metadata,omitempty"`
}

type NamedArtifacts map[string]baggageclaim.Volume
