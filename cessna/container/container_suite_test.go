package container_test

import (
	"code.cloudfoundry.org/garden/gardenfakes"
	"github.com/concourse/atc/cessna/cessnafakes"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"testing"
)

var (
	fakeWorker       *cessnafakes.FakeWorker
	fakeGardenClient *gardenfakes.FakeClient
)

func TestContainer(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Container Suite")
}
