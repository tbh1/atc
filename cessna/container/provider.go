package container

import (
	"code.cloudfoundry.org/garden"
	"code.cloudfoundry.org/lager"
	"github.com/concourse/atc/cessna"
)

type Provider struct {
	gardenClient garden.Client
	logger       lager.Logger
}

func NewProvider(worker cessna.Worker) *Provider {
	return &Provider{
		gardenClient: worker.GardenClient(),
		logger:       lager.NewLogger("worker-container-provider"),
	}
}

func (p *Provider) CreateContainer(rootFSPath string, bindMounts *[]garden.BindMount) (garden.Container, error) {

	gardenSpec := garden.ContainerSpec{
		Privileged: false,
		RootFSPath: rootFSPath,
	}

	if bindMounts != nil {
		gardenSpec.BindMounts = *bindMounts
	}

	container, err := p.gardenClient.Create(gardenSpec)
	if err != nil {
		p.logger.Error("failed-to-create-container-in-garden", err)
		return nil, err
	}

	return container, nil
}
