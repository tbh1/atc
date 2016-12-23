package cessna

import (
	"net/http"

	"code.cloudfoundry.org/garden"
	gclient "code.cloudfoundry.org/garden/client"
	"github.com/concourse/baggageclaim"
	bclient "github.com/concourse/baggageclaim/client"

	"code.cloudfoundry.org/garden/client/connection"
)

type Worker struct {
	GardenAddr       string
	BaggageclaimAddr string
}

func (w *Worker) GardenClient() garden.Client {
	return gclient.New(connection.New("tcp", w.GardenAddr))
}

func (w *Worker) BaggageClaimClient() baggageclaim.Client {
	return bclient.New(w.BaggageclaimAddr, http.DefaultTransport)
}

type Container interface {
	garden.Container
}

type Volume interface {
	baggageclaim.Volume
	COWify(properties baggageclaim.VolumeProperties, privileged bool) (Volume, error)
}
