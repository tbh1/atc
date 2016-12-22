package container

import (
	"code.cloudfoundry.org/garden"
	"code.cloudfoundry.org/lager"
	"github.com/concourse/atc/cessna"
)

type Sandwich struct {
	RootFSPath string
	logger     lager.Logger
}

func NewSandwich(rootFSPath string) Sandwich {
	return Sandwich{
		RootFSPath: rootFSPath,
		logger:     lager.NewLogger("sandwich"),
	}
}

func (cs *Sandwich) ContainerOn(worker cessna.Worker) (garden.Container, error) {
	gardenClient := worker.GardenClient()

	rootFSPath := cs.RootFSPath

	gardenSpec := garden.ContainerSpec{
		Privileged: false,
		RootFSPath: rootFSPath,
	}

	container, err := gardenClient.Create(gardenSpec)
	if err != nil {
		cs.logger.Error("failed-to-create-container-in-garden", err)
		return nil, err
	}

	return container, nil
}
