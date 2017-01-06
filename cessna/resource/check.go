package resource

import (
	"github.com/concourse/atc"
	"github.com/tedsuo/ifrit"
)

type checkRequest struct {
	Source  atc.Source  `json:"source"`
	Version atc.Version `json:"version"`
}

func (rc *resourceContainer) RunCheck() ([]atc.Version, error) {
	var versions []atc.Version

	runner := rc.RunScript("/opt/resource/check", nil, checkRequest{rc.resource.Source, nil}, &versions)

	checking := ifrit.Invoke(runner)

	err := <-checking.Wait()
	if err != nil {
		return nil, err
	}

	return versions, nil
}

// func (rc *resourceContainer) In() {
//
// }
//
// func (rc *resourceContainer) Out() {
//
// }
