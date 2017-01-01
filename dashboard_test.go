// dashboard_test.go
//
// Author::    Chirantan Mitra
// Copyright:: Copyright (c) 2015-2017. All rights reserved
// License::   MIT

package gocd_test

import (
	"github.com/chiku/gocd"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Dashboard", func() {
	Context("marshal to JSON", func() {
		It("has key names of pipelines starting with lower-case", func() {
			p1 := gocd.DashboardPipeline{
				Name: "Pipeline",
				Stages: []gocd.DashboardStage{
					gocd.DashboardStage{
						Name:   "Stage",
						Status: "Passed",
					},
				},
			}
			dashboard := gocd.Dashboard{p1}

			body, err := dashboard.ToJSON()
			Expect(err).To(Succeed())
			Expect(string(body)).To(Equal(`[{"name":"Pipeline","stages":[{"name":"Stage","status":"Passed"}]}]`))
		})
	})

	Context("filtered sort", func() {
		s1 := []gocd.DashboardStage{gocd.DashboardStage{Name: "Stage One", Status: "Passed"}}
		s2 := []gocd.DashboardStage{gocd.DashboardStage{Name: "Stage Two", Status: "Passed"}}
		s3 := []gocd.DashboardStage{gocd.DashboardStage{Name: "Stage Three", Status: "Passed"}}
		p1 := gocd.DashboardPipeline{Name: "Pipeline One", Stages: s1}
		p2 := gocd.DashboardPipeline{Name: "Pipeline Two", Stages: s2}
		p3 := gocd.DashboardPipeline{Name: "Pipeline Three", Stages: s3}
		dashboard := gocd.Dashboard{p2, p1, p3}

		Context("when all pipelines are demanded", func() {
			order := []string{"Pipeline One", "Pipeline Two", "Pipeline Three"}
			sortedDashboard, ignores := dashboard.FilteredSort(order)

			It("arranges pipelines by name in the given order", func() {
				Expect(sortedDashboard).To(HaveLen(3))

				Expect(sortedDashboard[0].Name).To(Equal("Pipeline One"))
				Expect(sortedDashboard[0].Stages).To(Equal(s1))

				Expect(sortedDashboard[1].Name).To(Equal("Pipeline Two"))
				Expect(sortedDashboard[1].Stages).To(Equal(s2))

				Expect(sortedDashboard[2].Name).To(Equal("Pipeline Three"))
				Expect(sortedDashboard[2].Stages).To(Equal(s3))
			})

			It("has no ignores", func() {
				Expect(ignores).To(BeEmpty())
			})
		})

		Context("when not all pipelines are not mentioned", func() {
			order := []string{"Pipeline One", "Pipeline Three"}
			sortedDashboard, ignores := dashboard.FilteredSort(order)

			It("leaves out non-matching pipelines", func() {
				Expect(sortedDashboard).To(HaveLen(2))
				Expect(sortedDashboard[0].Name).To(Equal("Pipeline One"))
				Expect(sortedDashboard[1].Name).To(Equal("Pipeline Three"))
			})

			It("collects the ignored pipelines", func() {
				Expect(ignores).To(HaveLen(1))
				Expect(ignores).To(ContainElement("Pipeline Two"))
			})
		})

		Context("when pipelines are in lower-case", func() {
			order := []string{"pipeline one", "pipeline three"}
			sortedDashboard, ignores := dashboard.FilteredSort(order)

			It("ignores case when matching pipelines by name", func() {
				Expect(sortedDashboard).To(HaveLen(2))
				Expect(sortedDashboard[0].Name).To(Equal("Pipeline One"))
				Expect(sortedDashboard[1].Name).To(Equal("Pipeline Three"))
			})

			It("maintains the original case of the ignored pipeline names", func() {
				Expect(ignores).To(HaveLen(1))
				Expect(ignores).To(ContainElement("Pipeline Two"))
			})
		})

		Context("when order contains an item that is not a pipeline name", func() {
			order := []string{"Pipeline One", "Pipeline Two", "Pipeline Four"}
			sortedDashboard, ignores := dashboard.FilteredSort(order)

			It("ignores the extra order", func() {
				Expect(sortedDashboard).To(HaveLen(2))
				Expect(sortedDashboard[0].Name).To(Equal("Pipeline One"))
				Expect(sortedDashboard[1].Name).To(Equal("Pipeline Two"))
				Expect(ignores).To(Equal([]string{"Pipeline Three"}))
			})
		})
	})

	Context("map names", func() {
		s1 := []gocd.DashboardStage{gocd.DashboardStage{Name: "Stage One", Status: "Passed"}}
		s2 := []gocd.DashboardStage{gocd.DashboardStage{Name: "Stage Two", Status: "Passed"}}
		p1 := gocd.DashboardPipeline{Name: "Pipeline One", Stages: s1}
		p2 := gocd.DashboardPipeline{Name: "Pipeline Two", Stages: s2}
		dashboard := gocd.Dashboard{p1, p2}

		Context("when all pipelines are mapped", func() {
			mapping := map[string]string{
				"Pipeline One": "Pipeline A",
				"Pipeline Two": "Pipeline B",
			}
			mappedDashboard := dashboard.MapNames(mapping)

			It("changes the names of the pipelines", func() {
				Expect(mappedDashboard).To(HaveLen(2))

				Expect(mappedDashboard[0].Name).To(Equal("Pipeline A"))
				Expect(mappedDashboard[0].Stages).To(Equal(s1))

				Expect(mappedDashboard[1].Name).To(Equal("Pipeline B"))
				Expect(mappedDashboard[1].Stages).To(Equal(s2))
			})
		})

		Context("when no pipelines are mapped", func() {
			mapping := map[string]string{}
			mappedDashboard := dashboard.MapNames(mapping)

			It("retains the original names of the pipelines", func() {
				Expect(mappedDashboard).To(HaveLen(2))

				Expect(mappedDashboard[0].Name).To(Equal("Pipeline One"))
				Expect(mappedDashboard[0].Stages).To(Equal(s1))

				Expect(mappedDashboard[1].Name).To(Equal("Pipeline Two"))
				Expect(mappedDashboard[1].Stages).To(Equal(s2))
			})
		})
	})
})
