package cessna

import (
	"net/http"

	"code.cloudfoundry.org/garden"
	gclient "code.cloudfoundry.org/garden/client"
	"github.com/concourse/baggageclaim"
	bclient "github.com/concourse/baggageclaim/client"

	"code.cloudfoundry.org/garden/client/connection"
)

func NewWorker(gardenAddr string, baggageclaimAddr string) Worker {
	return &worker{
		gardenAddr:       gardenAddr,
		baggageclaimAddr: baggageclaimAddr,
	}
}

type worker struct {
	gardenAddr       string
	baggageclaimAddr string
}

func (w *worker) GardenClient() garden.Client {
	return gclient.New(connection.New("tcp", w.gardenAddr))
}

func (w *worker) BaggageClaimClient() baggageclaim.Client {
	return bclient.New(w.baggageclaimAddr, http.DefaultTransport)
}

//go:generate counterfeiter . Worker

type Worker interface {
	GardenClient() garden.Client
	BaggageClaimClient() baggageclaim.Client
}

//go:generate counterfeiter . Container

type Container interface {
	garden.Container
}

//go:generate counterfeiter . Volume

type Volume interface {
	baggageclaim.Volume
	COWify(privileged bool) (Volume, error)
}
