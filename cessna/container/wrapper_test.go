package container_test

import (
	"errors"

	"code.cloudfoundry.org/garden/gardenfakes"
	. "github.com/concourse/atc/cessna/container"
	"github.com/tedsuo/ifrit"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var (
	testWrapper   *Wrapper
	fakeContainer *gardenfakes.FakeContainer
	fakeProcess   *gardenfakes.FakeProcess
)

var _ = Describe("Container Wrapper", func() {
	BeforeEach(func() {
		fakeContainer = new(gardenfakes.FakeContainer)
		fakeProcess = new(gardenfakes.FakeProcess)

		testWrapper = &Wrapper{
			Container: fakeContainer,
		}

	})

	Describe("RunScript", func() {
		var (
			scriptPath string
			args       []string
			input      map[string]string
			output     map[string]string
			runner     ifrit.Runner
			invokeErr  error
		)

		BeforeEach(func() {
			scriptPath = "/opt/whatever"

		})

		JustBeforeEach(func() {
			runner = testWrapper.RunScript(scriptPath, args, input, output)

			scriptInvoke := ifrit.Invoke(runner)
			invokeErr = <-scriptInvoke.Wait()
		})

		Context("running the container process returns an error", func() {
			BeforeEach(func() {
				fakeContainer.RunReturns(nil, errors.New("failed to run process"))
			})

			It("returns the error when invoking the runner", func() {
				Expect(invokeErr.Error()).To(Equal("failed to run process"))
			})
		})

	})

})
