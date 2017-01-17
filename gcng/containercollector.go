package gcng

import (
	"errors"
	"time"

	"code.cloudfoundry.org/garden"
	"code.cloudfoundry.org/garden/client"
	"code.cloudfoundry.org/garden/client/connection"
	"code.cloudfoundry.org/lager"
	"github.com/concourse/atc/dbng"
)

const HijackedContainerTimeout = 5 * time.Minute

//go:generate counterfeiter . ContainerProvider

type ContainerProvider interface {
	MarkContainersForDeletion() error
	FindContainersMarkedForDeletion() ([]dbng.DestroyingContainer, error)
	FindHijackedContainersForDeletion() ([]dbng.CreatedContainer, error)
}

type ContainerCollector struct {
	Logger              lager.Logger
	ContainerProvider   ContainerProvider
	WorkerProvider      dbng.WorkerFactory
	GardenClientFactory GardenClientFactory
}

type GardenClientFactory func(*dbng.Worker) (garden.Client, error)

func NewGardenClientFactory() GardenClientFactory {
	return func(w *dbng.Worker) (garden.Client, error) {
		if w.GardenAddr == nil {
			return nil, errors.New("worker-does-not-have-garden-address")
		}

		gconn := connection.New("tcp", *w.GardenAddr)
		return client.New(gconn), nil
	}
}

func (c *ContainerCollector) Run() error {
	workers, err := c.WorkerProvider.Workers()
	if err != nil {
		c.Logger.Error("failed-to-get-workers", err)
		return err
	}
	workersByName := map[string]*dbng.Worker{}
	for _, w := range workers {
		workersByName[w.Name] = w
	}

	hijackedContainersForDeletion, err := c.ContainerProvider.FindHijackedContainersForDeletion()
	if err != nil {
		c.Logger.Error("failed-to-get-hijacked-containers-for-deletion", err)
		return err
	}

	for _, hijackedContainer := range hijackedContainersForDeletion {
		w, found := workersByName[hijackedContainer.WorkerName()]
		if !found {
			c.Logger.Info("worker-not-found", lager.Data{
				"worker-name": container.WorkerName(),
			})
			continue
		}

		gclient, err := c.GardenClientFactory(w)
		if err != nil {
			c.Logger.Error("failed-to-get-garden-client-for-worker", err, lager.Data{
				"worker": w,
			})
			continue
		}

		gardenContainer, err := gclient.Lookup(hijackedContainer.Handle())
		if err != nil {
			if _, ok := err.(garden.ContainerNotFoundError); ok {
				c.Logger.Debug("hijacked-container-not-found-in-garden", lager.Data{
					"worker": w,
					"handle": container.Handle(),
				})

				err = hijackedContainer.Destroying()
				if err != nil {
					c.Logger.Error("failed-to-mark-container-as-destroying", err, lager.Data{
						"worker": w,
						"handle": container.Handle(),
					})
					continue
				}

				continue
			}

			c.Logger.Error("failed-to-lookup-garden-container", err, lager.Data{
				"worker": w,
				"handle": container.Handle(),
			})
			continue
		}

		err = gardenContainer.SetGraceTime(HijackedContainerTimeout)
		if err != nil {
			c.Logger.Error("failed-to-set-grace-time-on-hijacked-container", err, lager.Data{
				"worker": w,
				"handle": container.Handle(),
			})
			continue
		}

		err = hijackedContainer.Destroying()
		if err != nil {
			c.Logger.Error("failed-to-mark-container-as-destroying", err, lager.Data{
				"worker": w,
				"handle": container.Handle(),
			})
			continue
		}
	}

	err = c.ContainerProvider.MarkContainersForDeletion()
	if err != nil {
		c.Logger.Error("marking-build-containers-for-deletion", err)
		return err
	}

	cs, err := c.ContainerProvider.FindContainersMarkedForDeletion()
	if err != nil {
		c.Logger.Error("find-build-containers-for-deletion", err)
		return err
	}
	containerHandles := []string{}
	for _, container := range cs {
		containerHandles = append(containerHandles, container.Handle())
	}
	c.Logger.Debug("found-build-containers-for-deletion", lager.Data{
		"containers": containerHandles,
	})

	for _, container := range cs {
		w, found := workersByName[container.WorkerName()]
		if !found {
			c.Logger.Info("worker-not-found", lager.Data{
				"workername": container.WorkerName(),
			})
			continue
		}

		gclient, err := c.GardenClientFactory(w)
		if err != nil {
			c.Logger.Error("failed-to-get-garden-client-for-worker", err, lager.Data{
				"worker": w,
			})
			continue
		}

		// if container.IsHijacked() {
		// lookup
		// if found continue
		// }

		err = gclient.Destroy(container.Handle())
		if err != nil {
			c.Logger.Error("failed-to-destroy-garden-container", err, lager.Data{
				"worker": w,
				"handle": container.Handle(),
			})
			continue
		}

		ok, err := container.Destroy()
		if err != nil {
			c.Logger.Error("failed-to-destroy-database-container", err, lager.Data{
				"handle": container.Handle(),
			})
			continue
		}

		if !ok {
			c.Logger.Info("container-provider-container-not-found", lager.Data{
				"handle": container.Handle(),
			})
			continue
		}

		c.Logger.Debug("completed-deleting-container", lager.Data{
			"handle": container.Handle(),
		})
	}

	c.Logger.Debug("completed-deleting-containers")

	return nil
}
