// gocd_test.go
//
// Author::    Chirantan Mitra
// Copyright:: Copyright (c) 2015-2017. All rights reserved
// License::   MIT

package gocd_test

import (
	"testing"

	"github.com/chiku/gocd"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func TestToDashboard(t *testing.T) {
	stages := []gocd.Stage{{Name: "Stage One", Status: "Unknown"}}
	instances := []gocd.Instance{{Stages: stages}}
	pipelines := []gocd.Pipeline{{Name: "Pipeline One", Instances: instances}}
	group := gocd.PipelineGroup{Pipelines: pipelines}
	groups := gocd.PipelineGroups{group}

	dashboard := groups.ToDashboard()

	if len(dashboard) != 1 {
		t.Fatalf("Expected 1 item in dashboard, but has %d items: dashboard: %#v", len(dashboard), dashboard)
	}

	item0 := dashboard[0]
	if item0.Name != "Pipeline One" {
		t.Errorf("Expected proper pipeline name, but was: %s", item0.Name)
	}

	dashboardStages := dashboard[0].Stages
	if len(dashboardStages) != 1 {
		t.Fatalf("Expected 1 stage, but has %d stages: stages: %#v", len(dashboardStages), dashboardStages)
	}

	if dashboardStages[0].Name != "Stage One" || dashboardStages[0].Status != "Unknown" {
		t.Fatalf("Expected first stage to be proper, but was: %#v", dashboardStages[0])
	}
}

func TestToDashboardWithMultipleInstances(t *testing.T) {
	stagesForOldInstance := []gocd.Stage{{Name: "Stage Old", Status: "Failed"}}
	stagesForNewInstance := []gocd.Stage{{Name: "Stage New", Status: "Passed"}}
	oldInstance := gocd.Instance{Stages: stagesForOldInstance}
	newInstance := gocd.Instance{Stages: stagesForNewInstance}
	instances := []gocd.Instance{oldInstance, newInstance}
	pipelines := []gocd.Pipeline{{Name: "Pipeline One", Instances: instances}}
	groups := gocd.PipelineGroups{gocd.PipelineGroup{Pipelines: pipelines}}

	dashboard := groups.ToDashboard()

	if len(dashboard) != 1 {
		t.Fatalf("Expected 1 item in dashboard for latest instance, but has %d items: dashboard: %#v", len(dashboard), dashboard)
	}

	item0 := dashboard[0]
	if item0.Name != "Pipeline One" {
		t.Errorf("Expected proper pipeline name, but was: %s", item0.Name)
	}

	dashboardStages := dashboard[0].Stages
	if len(dashboardStages) != 1 {
		t.Fatalf("Expected 1 stage, but has %d stages: stages: %#v", len(dashboardStages), dashboardStages)
	}

	if dashboardStages[0].Name != "Stage New" || dashboardStages[0].Status != "Passed" {
		t.Fatalf("Expected first stage to be proper, but was: %#v", dashboardStages[0])
	}
}

func TestToDashboardWithCurrentStatusAsUnknown(t *testing.T) {
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

	if len(dashboard) != 1 {
		t.Fatalf("Expected 1 item in dashboard for latest instance, but has %d items: dashboard: %#v", len(dashboard), dashboard)
	}

	item0 := dashboard[0]
	if item0.Name != "Pipeline One" {
		t.Errorf("Expected proper pipeline name, but was: %s", item0.Name)
	}

	dashboardStages := dashboard[0].Stages
	if len(dashboardStages) != 1 {
		t.Fatalf("Expected 1 stage, but has %d stages: stages: %#v", len(dashboardStages), dashboardStages)
	}

	if dashboardStages[0].Name != "Stage X" || dashboardStages[0].Status != "Passed" {
		t.Fatalf("Expected first stage to use status from previous run, but was: %#v", dashboardStages[0])
	}
}

func TestToDashboardWithCurrentAndOlderStatusAsUnknown(t *testing.T) {
	stagesForLatestInstance := []gocd.Stage{{Name: "Stage X", Status: "Unknown"}}
	stagesForMinus1Instance := []gocd.Stage{{Name: "Stage X", Status: "Unknown"}}
	latestInstance := gocd.Instance{Stages: stagesForLatestInstance}
	minus1Instance := gocd.Instance{Stages: stagesForMinus1Instance}
	instances := []gocd.Instance{minus1Instance, latestInstance}
	pipelines := []gocd.Pipeline{{Name: "Pipeline One", Instances: instances}}
	groups := gocd.PipelineGroups{gocd.PipelineGroup{Pipelines: pipelines}}

	dashboard := groups.ToDashboard()

	if len(dashboard) != 1 {
		t.Fatalf("Expected 1 item in dashboard for latest instance, but has %d items: dashboard: %#v", len(dashboard), dashboard)
	}

	item0 := dashboard[0]
	if item0.Name != "Pipeline One" {
		t.Errorf("Expected proper pipeline name, but was: %s", item0.Name)
	}

	dashboardStages := dashboard[0].Stages
	if len(dashboardStages) != 1 {
		t.Fatalf("Expected 1 stage, but has %d stages: stages: %#v", len(dashboardStages), dashboardStages)
	}

	if dashboardStages[0].Name != "Stage X" || dashboardStages[0].Status != "Unknown" {
		t.Fatalf("Expected first stage to have unknown status, but was: %#v", dashboardStages[0])
	}
}

func TestToDashboardWithoutPipelineGroups(t *testing.T) {
	groups := gocd.PipelineGroups{}
	dashboard := groups.ToDashboard()

	if len(dashboard) != 0 {
		t.Errorf("Expected empty dashboard when no pipeline-groups: dashboard: %#v", dashboard)
	}
}

func TestToDashboardWithoutPipelines(t *testing.T) {
	groups := gocd.PipelineGroups{gocd.PipelineGroup{}}
	dashboard := groups.ToDashboard()

	if len(dashboard) != 0 {
		t.Errorf("Expected empty dashboard when pipeline-groups have no pipeline: dashboard: %#v", dashboard)
	}
}

func TestToDashboardWithoutInstances(t *testing.T) {
	instances := []gocd.Instance{}
	pipelines := []gocd.Pipeline{{Instances: instances}}
	group := gocd.PipelineGroup{Pipelines: pipelines}
	groups := gocd.PipelineGroups{group}

	dashboard := groups.ToDashboard()

	if len(dashboard) != 0 {
		t.Errorf("Expected empty dashboard when pipeline-groups have pipeline, but pipeline have no instances: dashboard: %#v", dashboard)
	}
}

func TestToDashboardWithoutStages(t *testing.T) {
	stages := []gocd.Stage{}
	instances := []gocd.Instance{{Stages: stages}}
	pipelines := []gocd.Pipeline{{Instances: instances}}
	group := gocd.PipelineGroup{Pipelines: pipelines}
	groups := gocd.PipelineGroups{group}

	dashboard := groups.ToDashboard()

	if len(dashboard) != 0 {
		t.Errorf("Expected empty dashboard when pipeline-groups have pipeline and pipeline has intances, but instances have no stages: dashboard: %#v", dashboard)
	}
}

var _ = Describe("PipelineGroups", func() {
	Context("with pipeline-group, pipeline, multiple instances and stage", func() {
		Context("with current and older statuses as unknown", func() {
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
