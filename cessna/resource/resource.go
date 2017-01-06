package resource

import (
	"code.cloudfoundry.org/garden"
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
	RunGet(atc.Version, atc.Params, cessna.Volume) (versionResult, error)
}

type resourceContainer struct {
	container.Wrapper
	resource Resource
}

func NewResource(resourceType ResourceType, source atc.Source) Resource {
	return Resource{
		ResourceType: resourceType,
		Source:       source,
		logger:       lager.NewLogger("resource"),
	}
}

func (r *Resource) Check(worker cessna.Worker, version *atc.Version) ([]atc.Version, error) {
	resourceContainer, err := r.createContainer(worker, nil)
	if err != nil {
		return []atc.Version{}, err
	}

	return resourceContainer.RunCheck()
}

func (r *Resource) Get(worker cessna.Worker, version *atc.Version, params atc.Params) (cessna.Volume, error) {

	volumeForGet, err := r.createCacheVolume(worker)
	if err != nil {
		return nil, err
	}

	bindMounts := &[]garden.BindMount{
		garden.BindMount{
			SrcPath: volumeForGet.Path(),
			DstPath: "/tmp/resource/get",
			Mode:    garden.BindMountModeRW,
		},
	}

	resourceContainer, err := r.createContainer(worker, bindMounts)
	if err != nil {
		return nil, err
	}

	_, err = resourceContainer.RunGet(*version, params, volumeForGet)
	if err != nil {
		return nil, err
	}

	return volumeForGet, nil
}

func (r *Resource) createResourceVolume(worker cessna.Worker) (cessna.Volume, error) {
	provider := volume.NewProvider(worker)
	return provider.COWFromRootFS(r.ResourceType.RootFSPath)
}

func (r *Resource) createCacheVolume(worker cessna.Worker) (cessna.Volume, error) {
	provider := volume.NewProvider(worker)
	return provider.CreateEmptyVolume("123", false)
}

func (r *Resource) createContainer(worker cessna.Worker, bindMounts *[]garden.BindMount) (ResourceContainer, error) {
	//IDEA:	maybe we need this separate from create container..
	volume, err := r.createResourceVolume(worker)
	if err != nil {
		return nil, err
	}
	//IDEA //

	provider := container.NewProvider(worker)
	gardenContainer, err := provider.CreateContainer(volume.Path(), bindMounts)

	if err != nil {
		r.logger.Error("failed-to-create-container-for-resource", err)
		return nil, err
	}

	return &resourceContainer{container.Wrapper{Container: gardenContainer}, *r}, nil
}
