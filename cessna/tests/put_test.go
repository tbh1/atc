package tests_test

import (
	"github.com/concourse/atc"
	. "github.com/concourse/atc/cessna/resource"
	"github.com/concourse/baggageclaim"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Put a resource", func() {

	var (
		getBaseResource Resource
		getVolume       baggageclaim.Volume
		getErr          error

		putBaseResource Resource
		putResponse     OutResponse
		putErr          error
	)

	Context("whose type is a base resource type", func() {

		BeforeEach(func() {
			source := atc.Source{
				"versions": []map[string]string{
					{"ref": "123"},
					{"beep": "boop"},
				},
			}

			getBaseResource = NewResource(baseResourceType, source)
			resourceManager = NewResourceManagerFor(testWorker)

			getVolume, getErr = resourceManager.Get(getBaseResource, &atc.Version{"beep": "boop"}, nil)

			putBaseResource = NewResource(baseResourceType, source)
		})

		JustBeforeEach(func() {
			putResponse, putErr = resourceManager.Put(putBaseResource, atc.Params{
				"path": "inputresource/version",
			}, NamedArtifacts{
				"inputresource": getVolume,
			})
		})

		It("runs the out script", func() {
			Expect(putErr).ShouldNot(HaveOccurred())
		})

		It("outputs the version that was in the path param file", func() {
			Expect(putResponse.Version).To(Equal(atc.Version{"beep": "boop"}))
		})

	})

})
