package transport

import (
	"net/http"
	"net/url"

	"github.com/concourse/atc"
	"github.com/concourse/atc/dbng"
)

type baggageclaimRoundTripper struct {
	worker                dbng.Worker
	innerRoundTripper     http.RoundTripper
	cachedBaggageclaimURL *string
}

func NewBaggageclaimRoundTripper(worker dbng.Worker, innerRoundTripper http.RoundTripper) http.RoundTripper {
	return &baggageclaimRoundTripper{
		innerRoundTripper:     innerRoundTripper,
		worker:                worker,
		cachedBaggageclaimURL: worker.BaggageclaimURL(),
	}
}

func (c *baggageclaimRoundTripper) RoundTrip(request *http.Request) (*http.Response, error) {
	if c.cachedBaggageclaimURL == nil {
		found, err := c.worker.Reload()
		if err != nil {
			return nil, err
		}

		if !found {
			return nil, ErrMissingWorker{WorkerName: c.worker.Name()}
		}

		if c.worker.State() == atc.WorkerStateStalled {
			return nil, ErrWorkerStalled{WorkerName: c.worker.Name()}
		}

		if c.worker.BaggageclaimURL() == nil {
			return nil, ErrWorkerBaggageclaimURLIsMissing{WorkerName: c.worker.Name()}
		}

		c.cachedBaggageclaimURL = c.worker.BaggageclaimURL()
	}

	baggageclaimURL, err := url.Parse(*c.cachedBaggageclaimURL)
	if err != nil {
		return nil, err
	}

	updatedURL := *request.URL
	updatedURL.Host = baggageclaimURL.Host

	updatedRequest := *request
	updatedRequest.URL = &updatedURL

	response, err := c.innerRoundTripper.RoundTrip(&updatedRequest)
	if err != nil {
		c.cachedBaggageclaimURL = nil
	}

	return response, err
}
