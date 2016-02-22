// gocd_test.go
//
// Author::    Chirantan Mitra
// Copyright:: Copyright (c) 2015-2016. All rights reserved
// License::   MIT

package gocd_test

import (
	"github.com/chiku/gocd"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("PipelineGroups", func() {
	Context("without pipeline-groups", func() {
		groups := gocd.PipelineGroups{}
		dashboard := groups.ToDashboard()

		It("has no simple-pipelines", func() {
			Expect(dashboard).To(BeEmpty())
		})
	})

	Context("with pipeline-group without pipelines", func() {
		groups := gocd.PipelineGroups{gocd.PipelineGroup{}}
		dashboard := groups.ToDashboard()

		It("has no simple-pipelines", func() {
			Expect(dashboard).To(BeEmpty())
		})
	})

	Context("with pipeline-group, pipeline without instances", func() {
		instances := []gocd.Instance{}
		pipelines := []gocd.Pipeline{{Instances: instances}}
		group := gocd.PipelineGroup{Pipelines: pipelines}
		groups := gocd.PipelineGroups{group}
		dashboard := groups.ToDashboard()

		It("has no simple-pipelines", func() {
			Expect(dashboard).To(BeEmpty())
		})
	})

	Context("with pipeline-group, pipeline, instance without stages", func() {
		stages := []gocd.Stage{}
		instances := []gocd.Instance{{Stages: stages}}
		pipelines := []gocd.Pipeline{{Instances: instances}}
		group := gocd.PipelineGroup{Pipelines: pipelines}
		groups := gocd.PipelineGroups{group}
		dashboard := groups.ToDashboard()

		It("has no simple-pipelines", func() {
			Expect(dashboard).To(BeEmpty())
		})
	})

	Context("with pipeline-group, pipeline, instance and stage", func() {
		stages := []gocd.Stage{{Name: "Stage One", Status: "Unknown"}}
		instances := []gocd.Instance{{Stages: stages}}
		pipelines := []gocd.Pipeline{{Name: "Pipeline One", Instances: instances}}
		group := gocd.PipelineGroup{Pipelines: pipelines}
		groups := gocd.PipelineGroups{group}
		dashboard := groups.ToDashboard()

		It("has a simple-pipeline", func() {
			Expect(dashboard).To(HaveLen(1))
			Expect(dashboard[0].Name).To(Equal("Pipeline One"))
			dashboardStages := dashboard[0].Stages
			Expect(dashboardStages).To(HaveLen(1))
			Expect(dashboardStages[0].Name).To(Equal("Stage One"))
			Expect(dashboardStages[0].Status).To(Equal("Unknown"))
		})
	})

	Context("with pipeline-group, pipeline, multiple instances and stage", func() {
		stagesForOldInstance := []gocd.Stage{{Name: "Stage Old", Status: "Failed"}}
		stagesForNewInstance := []gocd.Stage{{Name: "Stage New", Status: "Passed"}}
		oldInstance := gocd.Instance{Stages: stagesForOldInstance}
		newInstance := gocd.Instance{Stages: stagesForNewInstance}
		instances := []gocd.Instance{oldInstance, newInstance}
		pipelines := []gocd.Pipeline{{Name: "Pipeline One", Instances: instances}}
		groups := gocd.PipelineGroups{gocd.PipelineGroup{Pipelines: pipelines}}
		dashboard := groups.ToDashboard()

		It("ignores older instances", func() {
			Expect(dashboard).To(HaveLen(1))
			dashboardStages := dashboard[0].Stages
			Expect(dashboardStages).To(HaveLen(1))
			Expect(dashboardStages[0].Name).To(Equal("Stage New"))
			Expect(dashboardStages[0].Status).To(Equal("Passed"))
		})

		Context("with the current status as unknown", func() {
			stagesForLatestInstance := []gocd.Stage{{Name: "Stage X", Status: "Unknown"}}
			stagesForMinus1Instance := []gocd.Stage{{Name: "Stage X", Status: "Unknown"}}
			stagesForMinus2Instance := []gocd.Stage{{Name: "Stage X", Status: "Passed"}}
			latestInstance := gocd.Instance{Stages: stagesForLatestInstance}
			minus1Instance := gocd.Instance{Stages: stagesForMinus1Instance}
			minus2Instance := gocd.Instance{Stages: stagesForMinus2Instance}
			instances := []gocd.Instance{minus2Instance, minus1Instance, latestInstance}
			pipelines := []gocd.Pipeline{{Name: "Pipeline One", Instances: instances}}
			groups := gocd.PipelineGroups{gocd.PipelineGroup{Pipelines: pipelines}}
			dashboard := groups.ToDashboard()

			It("uses the status of the older build", func() {
				Expect(dashboard).To(HaveLen(1))
				dashboardStages := dashboard[0].Stages
				Expect(dashboardStages).To(HaveLen(1))
				Expect(dashboardStages[0].Name).To(Equal("Stage X"))
				Expect(dashboardStages[0].Status).To(Equal("Passed"))
			})
		})

		Context("with current and older statuses as unknown", func() {
			stagesForLatestInstance := []gocd.Stage{{Name: "Stage X", Status: "Unknown"}}
			stagesForMinus1Instance := []gocd.Stage{{Name: "Stage X", Status: "Unknown"}}
			latestInstance := gocd.Instance{Stages: stagesForLatestInstance}
			minus1Instance := gocd.Instance{Stages: stagesForMinus1Instance}
			instances := []gocd.Instance{minus1Instance, latestInstance}
			pipelines := []gocd.Pipeline{{Name: "Pipeline One", Instances: instances}}
			groups := gocd.PipelineGroups{gocd.PipelineGroup{Pipelines: pipelines}}
			dashboard := groups.ToDashboard()

			It("has unknown status", func() {
				Expect(dashboard).To(HaveLen(1))
				dashboardStages := dashboard[0].Stages
				Expect(dashboardStages).To(HaveLen(1))
				Expect(dashboardStages[0].Name).To(Equal("Stage X"))
				Expect(dashboardStages[0].Status).To(Equal("Unknown"))
			})

			Context("with known previous instance status", func() {
				stagesForLatestInstance := []gocd.Stage{{Name: "Stage X", Status: "Unknown"}}
				stagesForMinus1Instance := []gocd.Stage{{Name: "Stage X", Status: "Unknown"}}
				latestInstance := gocd.Instance{Stages: stagesForLatestInstance}
				minus1Instance := gocd.Instance{Stages: stagesForMinus1Instance}
				previousInstance := gocd.PreviousInstance{Result: "Passed"}
				instances := []gocd.Instance{minus1Instance, latestInstance}
				pipelines := []gocd.Pipeline{{Name: "Pipeline One", Instances: instances, PreviousInstance: previousInstance}}
				groups := gocd.PipelineGroups{gocd.PipelineGroup{Pipelines: pipelines}}
				dashboard := groups.ToDashboard()

				It("uses the status of previous instance", func() {
					Expect(dashboard).To(HaveLen(1))
					dashboardStages := dashboard[0].Stages
					Expect(dashboardStages).To(HaveLen(1))
					Expect(dashboardStages[0].Name).To(Equal("Stage X"))
					Expect(dashboardStages[0].Status).To(Equal("Passed"))
				})
			})
		})

		Context("with previous result as failed and current status as building", func() {
			stagesForLatestInstance := []gocd.Stage{{Name: "Stage X", Status: "Building"}}
			latestInstance := gocd.Instance{Stages: stagesForLatestInstance}
			previousInstance := gocd.PreviousInstance{Result: "Failed"}
			instances := []gocd.Instance{latestInstance}
			pipelines := []gocd.Pipeline{{Name: "Pipeline One", Instances: instances, PreviousInstance: previousInstance}}
			groups := gocd.PipelineGroups{gocd.PipelineGroup{Pipelines: pipelines}}
			dashboard := groups.ToDashboard()

			It("uses marks the status as recovering", func() {
				Expect(dashboard).To(HaveLen(1))
				dashboardStages := dashboard[0].Stages
				Expect(dashboardStages).To(HaveLen(1))
				Expect(dashboardStages[0].Name).To(Equal("Stage X"))
				Expect(dashboardStages[0].Status).To(Equal("Recovering"))
			})
		})
	})

	Context("with multiple pipeline-group, pipelines, instances and stages", func() {
		stage_1_old_1 := gocd.Stage{Name: "Stage 1.1.1", Status: "Passed"}
		stage_1_old_2 := gocd.Stage{Name: "Stage 1.1.2", Status: "Failed"}
		stage_1_new_1 := gocd.Stage{Name: "Stage 1.2.1", Status: "Cancelled"}
		stage_1_new_2 := gocd.Stage{Name: "Stage 1.2.2", Status: "Failing"}
		stage_2_old_1 := gocd.Stage{Name: "Stage 2.1.1", Status: "Building"}
		stage_2_old_2 := gocd.Stage{Name: "Stage 2.1.2", Status: "Unknown"}
		stage_2_new_1 := gocd.Stage{Name: "Stage 2.2.1", Status: "Passed"}
		stage_2_new_2 := gocd.Stage{Name: "Stage 2.2.2", Status: "Failed"}
		stages_1_old := []gocd.Stage{stage_1_old_1, stage_1_old_2}
		stages_1_new := []gocd.Stage{stage_1_new_1, stage_1_new_2}
		stages_2_old := []gocd.Stage{stage_2_old_1, stage_2_old_2}
		stages_2_new := []gocd.Stage{stage_2_new_1, stage_2_new_2}
		instance_1_old := gocd.Instance{Stages: stages_1_old}
		instance_1_new := gocd.Instance{Stages: stages_1_new}
		instance_2_old := gocd.Instance{Stages: stages_2_old}
		instance_2_new := gocd.Instance{Stages: stages_2_new}
		instances_1 := []gocd.Instance{instance_1_old, instance_1_new}
		instances_2 := []gocd.Instance{instance_2_old, instance_2_new}
		pipeline_1 := gocd.Pipeline{Instances: instances_1, Name: "Pipeline One"}
		pipeline_2 := gocd.Pipeline{Instances: instances_2, Name: "Pipeline Two"}
		pipelines := []gocd.Pipeline{pipeline_1, pipeline_2}
		groups := gocd.PipelineGroups{gocd.PipelineGroup{Pipelines: pipelines}}
		dashboard := groups.ToDashboard()

		It("has simple-pipelines", func() {
			Expect(dashboard).To(HaveLen(2))
			dashboard_1 := dashboard[0]
			dashboard_2 := dashboard[1]
			Expect(dashboard_1.Name).To(Equal("Pipeline One"))
			Expect(dashboard_2.Name).To(Equal("Pipeline Two"))
			dashboardStages_1 := dashboard_1.Stages
			dashboardStages_2 := dashboard_2.Stages
			Expect(dashboardStages_1).To(HaveLen(2))
			Expect(dashboardStages_2).To(HaveLen(2))
			Expect(dashboardStages_1[0]).To(Equal(gocd.DashboardStage{Name: "Stage 1.2.1", Status: "Cancelled"}))
			Expect(dashboardStages_1[1]).To(Equal(gocd.DashboardStage{Name: "Stage 1.2.2", Status: "Failing"}))
			Expect(dashboardStages_2[0]).To(Equal(gocd.DashboardStage{Name: "Stage 2.2.1", Status: "Passed"}))
			Expect(dashboardStages_2[1]).To(Equal(gocd.DashboardStage{Name: "Stage 2.2.2", Status: "Failed"}))
		})
	})

	Context("unmarshal from JSON", func() {
		const dashboardJSON = `[{
		  "name": "Group",
		  "pipelines": [
		    {
		      "name": "Pipeline",
		      "instances": [
		        {
		          "stages": [
		            { "name": "StageOne", "status": "Passed" },
		            { "name": "StageTwo", "status": "Building" }
		          ]
		        }
		      ],
		      "previous_instance": {
		        "result": "Passed"
		      }
		    }
		  ]
		}
		]`

		It("is created from byte array", func() {
			groups, err := gocd.NewPipelineGroups([]byte(dashboardJSON))
			Expect(err).To(Succeed())
			Expect(groups).To(HaveLen(1))
			Expect(groups[0].Pipelines).To(HaveLen(1))
			pipeline := groups[0].Pipelines[0]
			Expect(pipeline.Name).To(Equal("Pipeline"))
			Expect(pipeline.Instances).To(HaveLen(1))
			instance := pipeline.Instances[0]
			Expect(instance.Stages).To(HaveLen(2))
			stages := instance.Stages
			Expect(stages[0].Name).To(Equal("StageOne"))
			Expect(stages[0].Status).To(Equal("Passed"))
			Expect(stages[1].Name).To(Equal("StageTwo"))
			Expect(stages[1].Status).To(Equal("Building"))
			previousInstance := pipeline.PreviousInstance
			Expect(previousInstance.Result).To(Equal("Passed"))
		})

		Context("on failure", func() {
			It("has error", func() {
				groups, err := gocd.NewPipelineGroups([]byte(`Random`))
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("error unmarshalling Gocd JSON: "))
				Expect(groups).To(BeNil())
			})
		})
	})
})
