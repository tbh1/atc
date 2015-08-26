package factory_test

import (
	"github.com/concourse/atc"
	"github.com/concourse/atc/scheduler/factory"
	"github.com/concourse/atc/scheduler/factory/fakes"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Factory Do", func() {
	var (
		fakeLocationPopulator *fakes.FakeLocationPopulator
		buildFactory          factory.BuildFactory

		resources atc.ResourceConfigs
	)

	BeforeEach(func() {
		fakeLocationPopulator = &fakes.FakeLocationPopulator{}

		buildFactory = factory.NewBuildFactory(
			"some-pipeline",
			fakeLocationPopulator,
		)

		resources = atc.ResourceConfigs{
			{
				Name:   "some-resource",
				Type:   "git",
				Source: atc.Source{"uri": "git://some-resource"},
			},
		}
	})

	Context("when I have a nested do ", func() {
		It("returns the correct plan", func() {

			actual, err := buildFactory.Create(atc.JobConfig{
				Plan: atc.PlanSequence{
					{
						Do: &atc.PlanSequence{
							{
								Task: "some thing",
							},
							{
								Task: "some thing-2",
							},
							{
								Do: &atc.PlanSequence{
									{
										Task: "some other thing",
									},
								},
							},
						},
					},
				},
			}, resources, nil)
			Ω(err).ShouldNot(HaveOccurred())

			expected := atc.Plan{
				OnSuccess: &atc.OnSuccessPlan{
					Step: atc.Plan{
						Task: &atc.TaskPlan{
							Name: "some thing",
						},
					},
					Next: atc.Plan{
						OnSuccess: &atc.OnSuccessPlan{
							Step: atc.Plan{
								Task: &atc.TaskPlan{
									Name: "some thing-2",
								},
							},
							Next: atc.Plan{
								Task: &atc.TaskPlan{
									Name: "some other thing",
								},
							},
						},
					},
				},
			}
			Ω(actual).Should(Equal(expected))
		})
	})

	Context("when I have an aggregate inside a do", func() {
		It("returns the correct plan", func() {

			actual, err := buildFactory.Create(atc.JobConfig{
				Plan: atc.PlanSequence{
					{
						Do: &atc.PlanSequence{
							{
								Task: "some thing",
							},
							{
								Aggregate: &atc.PlanSequence{
									{
										Task: "some other thing",
									},
								},
							},
							{
								Task: "some thing-2",
							},
						},
					},
				},
			}, resources, nil)
			Ω(err).ShouldNot(HaveOccurred())

			expected := atc.Plan{
				OnSuccess: &atc.OnSuccessPlan{
					Step: atc.Plan{
						Task: &atc.TaskPlan{
							Name: "some thing",
						},
					},
					Next: atc.Plan{
						OnSuccess: &atc.OnSuccessPlan{
							Step: atc.Plan{

								Aggregate: &atc.AggregatePlan{
									{
										Task: &atc.TaskPlan{
											Name: "some other thing",
										},
									},
								},
							},
							Next: atc.Plan{
								Task: &atc.TaskPlan{
									Name: "some thing-2",
								},
							},
						},
					},
				},
			}
			Ω(actual).Should(Equal(expected))
		})
	})

	Context("when i have a do inside an aggregate inside a hook", func() {
		It("returns the correct plan", func() {

			actual, err := buildFactory.Create(atc.JobConfig{
				Plan: atc.PlanSequence{
					{
						Task: "starting-task",
						Success: &atc.PlanConfig{
							Aggregate: &atc.PlanSequence{
								{
									Task: "some thing",
								},
								{
									Do: &atc.PlanSequence{
										{
											Task: "some other thing",
										},
									},
								},
							},
						},
					},
				},
			}, resources, nil)
			Ω(err).ShouldNot(HaveOccurred())

			expected := atc.Plan{
				OnSuccess: &atc.OnSuccessPlan{
					Step: atc.Plan{
						Task: &atc.TaskPlan{
							Name: "starting-task",
						},
					},
					Next: atc.Plan{
						Aggregate: &atc.AggregatePlan{
							{
								Task: &atc.TaskPlan{
									Name: "some thing",
								},
							},
							{
								Task: &atc.TaskPlan{
									Name: "some other thing",
								},
							},
						},
					},
				},
			}

			Ω(actual).Should(Equal(expected))
		})
	})

	Context("when I have a do inside an aggregate", func() {
		It("returns the correct plan", func() {

			actual, err := buildFactory.Create(atc.JobConfig{
				Plan: atc.PlanSequence{
					{
						Aggregate: &atc.PlanSequence{
							{
								Task: "some thing",
							},
							{
								Do: &atc.PlanSequence{
									{
										Task: "some other thing",
									},
									{
										Task: "some other thing-2",
									},
								},
							},
							{
								Task: "some thing-2",
							},
						},
					},
				},
			}, resources, nil)
			Ω(err).ShouldNot(HaveOccurred())

			expected := atc.Plan{
				Aggregate: &atc.AggregatePlan{
					{
						Task: &atc.TaskPlan{
							Name: "some thing",
						},
					},
					{
						OnSuccess: &atc.OnSuccessPlan{
							Step: atc.Plan{
								Task: &atc.TaskPlan{
									Name: "some other thing",
								},
							},
							Next: atc.Plan{
								Task: &atc.TaskPlan{
									Name: "some other thing-2",
								},
							},
						},
					},
					{
						Task: &atc.TaskPlan{
							Name: "some thing-2",
						},
					},
				},
			}

			Ω(actual).Should(Equal(expected))
		})
	})
})