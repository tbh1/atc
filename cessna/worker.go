package cessna

import (
	"code.cloudfoundry.org/garden"
	"code.cloudfoundry.org/garden/client"
	"code.cloudfoundry.org/garden/client/connection"
)

type Worker struct {
	GardenAddr       string
	BaggageclaimAddr string
}

func (w *Worker) GardenClient() garden.Client {
	return client.New(connection.New("tcp", w.GardenAddr))
}
