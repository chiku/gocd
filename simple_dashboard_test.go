// app/simple_dashboard_test.go
//
// Author::    Chirantan Mitra
// Copyright:: Copyright (c) 2015-2016. All rights reserved
// License::   MIT

package gocd_test

import (
	"sort"

	"github.com/chiku/gocd"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("SimpleDashboard", func() {
	Context("marshal to JSON", func() {
		It("has key names of pipelines starting with lower-case", func() {
			simpleDashboard := gocd.SimpleDashboard{
				Pipelines: []gocd.SimplePipeline{
					gocd.SimplePipeline{
						Name: "Pipeline",
						Stages: []gocd.SimpleStage{
							gocd.SimpleStage{
								Name:   "Stage",
								Status: "Passed",
							},
						},
					},
				},
			}

			body, err := simpleDashboard.ToJSON()
			Expect(err).To(BeNil())
			Expect(string(body)).To(Equal(`[{"name":"Pipeline","stages":[{"name":"Stage","status":"Passed"}]}]`))
		})
	})

	Context("sort by order", func() {
		It("sorts in ascending order", func() {
			p1 := (gocd.SimplePipeline{
				Name:   "PipelineOne",
				Stages: []gocd.SimpleStage{gocd.SimpleStage{Name: "Stage", Status: "Passed"}},
			}).Order(2)
			p2 := (gocd.SimplePipeline{
				Name:   "PipelineTwo",
				Stages: []gocd.SimpleStage{gocd.SimpleStage{Name: "Stage", Status: "Passed"}},
			}).Order(1)
			pipelines := []gocd.SimplePipeline{p1, p2}

			sort.Sort(gocd.ByOrder(pipelines))

			Expect(pipelines).To(HaveLen(2))
			Expect(pipelines[0].Name).To(Equal("PipelineTwo"))
			Expect(pipelines[1].Name).To(Equal("PipelineOne"))
		})
	})
})
