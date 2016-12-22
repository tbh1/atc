package resource

import (
	"code.cloudfoundry.org/lager"
	"github.com/concourse/atc"
	"github.com/concourse/atc/cessna"
	"github.com/concourse/atc/cessna/container"
)

type ResourceType struct {
	RootFSPath string
	Name       string
}

type Resource struct {
	ResourceType ResourceType
	Source       atc.Source
	logger       lager.Logger
}

func NewResource(resourceType ResourceType, source atc.Source) Resource {
	return Resource{
		ResourceType: resourceType,
		Source:       source,
		logger:       lager.NewLogger("resource"),
	}
}

func (r *Resource) Check(worker cessna.Worker, version *atc.Version) ([]atc.Version, error) {
	//resource type rootfs stuff
	resourceContainer, err := r.containerFor(worker)
	if err != nil {
		return []atc.Version{}, err
	}

	return resourceContainer.Check()
}

func (r *Resource) containerFor(worker cessna.Worker) (ResourceContainer, error) {

	cs := container.NewSandwich(r.ResourceType.RootFSPath)

	gardenContainer, err := cs.ContainerOn(worker)
	if err != nil {
		r.logger.Error("failed-to-create-container-for-resource", err)
		return nil, err
	}

	return &resourceContainer{container.Wrapper{gardenContainer}, *r}, nil
}
