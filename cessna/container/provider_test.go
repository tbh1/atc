package container_test

import (
	"code.cloudfoundry.org/garden/gardenfakes"
	"github.com/concourse/atc/cessna/cessnafakes"
	. "github.com/concourse/atc/cessna/container"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var testProvider *Provider

var _ = Describe("Provider", func() {
	BeforeEach(func() {
		fakeWorker = new(cessnafakes.FakeWorker)
		fakeGardenClient = new(gardenfakes.FakeClient)
		fakeWorker.GardenClientReturns(fakeGardenClient)

		testProvider = NewProvider(fakeWorker)
	})

	Describe("CreateContainer", func() {
		var rootFsPath string

		BeforeEach(func() {
			rootFsPath = "/some/rootfs/path/"
			testProvider.CreateContainer(rootFsPath)
		})

		It("Creates a container using the worker's garden client", func() {
			Expect(fakeGardenClient.CreateCallCount()).To(Equal(1))
		})

		It("uses the rootfs path provided in the garden container spec", func() {
			spec := fakeGardenClient.CreateArgsForCall(0)
			Expect(spec.RootFSPath).To(Equal("/some/rootfs/path/"))

		})

	})

})
