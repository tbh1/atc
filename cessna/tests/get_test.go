package tests_test

import (
	"archive/tar"
	"io/ioutil"

	"github.com/concourse/atc"
	"github.com/concourse/atc/cessna"
	. "github.com/concourse/atc/cessna/resource"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Get version of a resource", func() {

	var getVolume cessna.Volume
	var getErr error

	Context("whose type is a base resource type", func() {

		BeforeEach(func() {
			source := atc.Source{
				"versions": []map[string]string{
					{"ref": "123"},
					{"beep": "boop"},
				},
			}

			testBaseResource = NewResource(baseResourceType, source)
			resourceManager = NewResourceManagerFor(testWorker)
		})

		JustBeforeEach(func() {
			getVolume, getErr = resourceManager.Get(testBaseResource, &atc.Version{"beep": "boop"}, nil)
		})

		It("runs the get script", func() {
			Expect(getErr).ShouldNot(HaveOccurred())
		})

		It("returns a volume with the result of running the get script", func() {
			file, err := getVolume.StreamOut("/version")
			Expect(err).ToNot(HaveOccurred())
			defer file.Close()

			tarReader := tar.NewReader(file)
			tarReader.Next()

			bytes, err := ioutil.ReadAll(tarReader)
			Expect(err).NotTo(HaveOccurred())
			Expect(bytes).To(MatchJSON(`{"beep": "boop"}`))
		})

	})

})
