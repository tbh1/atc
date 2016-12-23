package resource

import (
	"code.cloudfoundry.org/lager"
	"github.com/concourse/atc"
	"github.com/concourse/atc/cessna"
	"github.com/concourse/atc/cessna/container"
	"github.com/concourse/atc/cessna/volume"
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

type ResourceContainer interface {
	RunCheck() ([]atc.Version, error)
}

func NewResource(resourceType ResourceType, source atc.Source) Resource {
	return Resource{
		ResourceType: resourceType,
		Source:       source,
		logger:       lager.NewLogger("resource"),
	}
}

func (r *Resource) Check(worker cessna.Worker, version *atc.Version) ([]atc.Version, error) {
	resourceVolume, err := r.createCheckVolume(worker)
	if err != nil {
		return []atc.Version{}, err
	}

	resourceContainer, err := r.createContainer(worker, resourceVolume)
	if err != nil {
		return []atc.Version{}, err
	}

	return resourceContainer.RunCheck()
}

func (r *Resource) createCheckVolume(worker cessna.Worker) (cessna.Volume, error) {
	provider := volume.NewProvider(worker)
	return provider.COWFromRootFS(r.ResourceType.RootFSPath)
}

func (r *Resource) createContainer(worker cessna.Worker, volume cessna.Volume) (ResourceContainer, error) {
	provider := container.NewProvider(worker)
	gardenContainer, err := provider.CreateContainer(volume.Path())

	if err != nil {
		r.logger.Error("failed-to-create-container-for-resource", err)
		return nil, err
	}

	return &resourceContainer{container.Wrapper{gardenContainer}, *r}, nil
}
