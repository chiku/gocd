// gocd_test.go
//
// Author::    Chirantan Mitra
// Copyright:: Copyright (c) 2015-2017. All rights reserved
// License::   MIT

package gocd_test

import (
	"reflect"
	"strings"
	"testing"

	"github.com/chiku/gocd"
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

	dashboard0 := dashboard[0]
	if dashboard0.Name != "Pipeline One" {
		t.Errorf("Expected proper pipeline name, but was: %s", dashboard0.Name)
	}

	dashboardStages := dashboard[0].Stages
	if len(dashboardStages) != 1 {
		t.Fatalf("Expected 1 stage, but has %d stages: stages: %#v", len(dashboardStages), dashboardStages)
	}

	if dashboardStages[0].Name != "Stage One" || dashboardStages[0].Status != "Unknown" {
		t.Errorf("Expected first stage to be proper, but was: %#v", dashboardStages[0])
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

	dashboard0 := dashboard[0]
	if dashboard0.Name != "Pipeline One" {
		t.Errorf("Expected proper pipeline name, but was: %s", dashboard0.Name)
	}

	dashboardStages := dashboard[0].Stages
	if len(dashboardStages) != 1 {
		t.Fatalf("Expected 1 stage, but has %d stages: stages: %#v", len(dashboardStages), dashboardStages)
	}

	if dashboardStages[0].Name != "Stage New" || dashboardStages[0].Status != "Passed" {
		t.Errorf("Expected first stage to be proper, but was: %#v", dashboardStages[0])
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

	dashboard0 := dashboard[0]
	if dashboard0.Name != "Pipeline One" {
		t.Errorf("Expected proper pipeline name, but was: %s", dashboard0.Name)
	}

	dashboardStages := dashboard[0].Stages
	if len(dashboardStages) != 1 {
		t.Fatalf("Expected 1 stage, but has %d stages: stages: %#v", len(dashboardStages), dashboardStages)
	}

	if dashboardStages[0].Name != "Stage X" || dashboardStages[0].Status != "Passed" {
		t.Errorf("Expected first stage to use status from previous builds, but was: %#v", dashboardStages[0])
	}
}

func TestToDashboardWithCurrentAndOlderStatusesAsUnknown(t *testing.T) {
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

	dashboard0 := dashboard[0]
	if dashboard0.Name != "Pipeline One" {
		t.Errorf("Expected proper pipeline name, but was: %s", dashboard0.Name)
	}

	dashboardStages := dashboard[0].Stages
	if len(dashboardStages) != 1 {
		t.Fatalf("Expected 1 stage, but has %d stages: stages: %#v", len(dashboardStages), dashboardStages)
	}

	if dashboardStages[0].Name != "Stage X" || dashboardStages[0].Status != "Unknown" {
		t.Errorf("Expected first stage to have unknown status, but was: %#v", dashboardStages[0])
	}
}

func TestToDashboardWithCurrentAndOlderStatusesUnknownButPreviousInstance(t *testing.T) {
	stagesForLatestInstance := []gocd.Stage{{Name: "Stage X", Status: "Unknown"}}
	stagesForMinus1Instance := []gocd.Stage{{Name: "Stage X", Status: "Unknown"}}
	latestInstance := gocd.Instance{Stages: stagesForLatestInstance}
	minus1Instance := gocd.Instance{Stages: stagesForMinus1Instance}
	previousInstance := gocd.PreviousInstance{Result: "Passed"}
	instances := []gocd.Instance{minus1Instance, latestInstance}
	pipelines := []gocd.Pipeline{{Name: "Pipeline One", Instances: instances, PreviousInstance: previousInstance}}
	groups := gocd.PipelineGroups{gocd.PipelineGroup{Pipelines: pipelines}}

	dashboard := groups.ToDashboard()

	if len(dashboard) != 1 {
		t.Fatalf("Expected 1 item in dashboard for latest instance, but has %d items: dashboard: %#v", len(dashboard), dashboard)
	}

	dashboard0 := dashboard[0]
	if dashboard0.Name != "Pipeline One" {
		t.Errorf("Expected proper pipeline name, but was: %s", dashboard0.Name)
	}

	dashboardStages := dashboard[0].Stages
	if len(dashboardStages) != 1 {
		t.Fatalf("Expected 1 stage, but has %d stages: stages: %#v", len(dashboardStages), dashboardStages)
	}

	if dashboardStages[0].Name != "Stage X" || dashboardStages[0].Status != "Passed" {
		t.Errorf("Expected first stage to have status from previous instance run, but was: %#v", dashboardStages[0])
	}
}

func TestToDashboardWithCurrentStatusBuildingButPreviousFailed(t *testing.T) {
	stagesForLatestInstance := []gocd.Stage{{Name: "Stage X", Status: "Building"}}
	latestInstance := gocd.Instance{Stages: stagesForLatestInstance}
	previousInstance := gocd.PreviousInstance{Result: "Failed"}
	instances := []gocd.Instance{latestInstance}
	pipelines := []gocd.Pipeline{{Name: "Pipeline One", Instances: instances, PreviousInstance: previousInstance}}
	groups := gocd.PipelineGroups{gocd.PipelineGroup{Pipelines: pipelines}}

	dashboard := groups.ToDashboard()

	if len(dashboard) != 1 {
		t.Fatalf("Expected 1 item in dashboard for latest instance, but has %d items: dashboard: %#v", len(dashboard), dashboard)
	}

	dashboard0 := dashboard[0]
	if dashboard0.Name != "Pipeline One" {
		t.Errorf("Expected proper pipeline name, but was: %s", dashboard0.Name)
	}

	dashboardStages := dashboard[0].Stages
	if len(dashboardStages) != 1 {
		t.Fatalf("Expected 1 stage, but has %d stages: stages: %#v", len(dashboardStages), dashboardStages)
	}

	if dashboardStages[0].Name != "Stage X" || dashboardStages[0].Status != "Recovering" {
		t.Errorf("Expected first stage to have recovering status, but was: %#v", dashboardStages[0])
	}
}

func TestToDashboardWithMultiplePipelineGroupsPipelinesInstancesAndStages(t *testing.T) {
	stage1old1 := gocd.Stage{Name: "Stage 1.1.1", Status: "Passed"}
	stage1old2 := gocd.Stage{Name: "Stage 1.1.2", Status: "Failed"}
	stage1new1 := gocd.Stage{Name: "Stage 1.2.1", Status: "Cancelled"}
	stage1new2 := gocd.Stage{Name: "Stage 1.2.2", Status: "Failing"}
	stage2old1 := gocd.Stage{Name: "Stage 2.1.1", Status: "Building"}
	stage2old2 := gocd.Stage{Name: "Stage 2.1.2", Status: "Unknown"}
	stage2new1 := gocd.Stage{Name: "Stage 2.2.1", Status: "Passed"}
	stage2new2 := gocd.Stage{Name: "Stage 2.2.2", Status: "Failed"}
	stages1old := []gocd.Stage{stage1old1, stage1old2}
	stages1new := []gocd.Stage{stage1new1, stage1new2}
	stages2old := []gocd.Stage{stage2old1, stage2old2}
	stages2new := []gocd.Stage{stage2new1, stage2new2}
	instance1old := gocd.Instance{Stages: stages1old}
	instance1new := gocd.Instance{Stages: stages1new}
	instance2old := gocd.Instance{Stages: stages2old}
	instance2new := gocd.Instance{Stages: stages2new}
	instances1 := []gocd.Instance{instance1old, instance1new}
	instances2 := []gocd.Instance{instance2old, instance2new}
	pipeline1 := gocd.Pipeline{Instances: instances1, Name: "Pipeline One"}
	pipeline2 := gocd.Pipeline{Instances: instances2, Name: "Pipeline Two"}
	pipelines := []gocd.Pipeline{pipeline1, pipeline2}
	groups := gocd.PipelineGroups{gocd.PipelineGroup{Pipelines: pipelines}}

	dashboard := groups.ToDashboard()

	if len(dashboard) != 2 {
		t.Fatalf("Expected 2 items in dashboard, but has %d item(s): dashboard: %#v", len(dashboard), dashboard)
	}

	dashboard0 := dashboard[0]
	if dashboard0.Name != "Pipeline One" {
		t.Errorf("Expected proper pipeline name for first pipeline, but was: %s", dashboard0.Name)
	}

	dashboard0Stages := dashboard[0].Stages
	if len(dashboard0Stages) != 2 {
		t.Fatalf("Expected 2 stages, but has %d stage(s) in first pipeline: stages: %#v", len(dashboard0Stages), dashboard0Stages)
	}

	if !reflect.DeepEqual(dashboard0Stages[0], gocd.DashboardStage{Name: "Stage 1.2.1", Status: "Cancelled"}) {
		t.Errorf("Expected proper pipeline 1 stage 1, but was: #%v", dashboard0Stages[0])
	}
	if !reflect.DeepEqual(dashboard0Stages[1], gocd.DashboardStage{Name: "Stage 1.2.2", Status: "Failing"}) {
		t.Errorf("Expected proper pipeline 1 stage 2, but was: #%v", dashboard0Stages[1])
	}

	dashboard1 := dashboard[1]
	if dashboard1.Name != "Pipeline Two" {
		t.Errorf("Expected proper pipeline name for second pipeline, but was: %s", dashboard1.Name)
	}

	dashboard1Stages := dashboard[1].Stages
	if len(dashboard1Stages) != 2 {
		t.Fatalf("Expected 2 stage, but has %d stage(s) in second pipeline: stages: %#v", len(dashboard1Stages), dashboard1Stages)
	}

	if !reflect.DeepEqual(dashboard1Stages[0], gocd.DashboardStage{Name: "Stage 2.2.1", Status: "Passed"}) {
		t.Errorf("Expected proper pipeline 1 stage 1, but was: #%v", dashboard1Stages[0])
	}
	if !reflect.DeepEqual(dashboard1Stages[1], gocd.DashboardStage{Name: "Stage 2.2.2", Status: "Failed"}) {
		t.Errorf("Expected proper pipeline 1 stage 2, but was: #%v", dashboard1Stages[1])
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

func TestNewPipelineGroups(t *testing.T) {
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

	groups, err := gocd.NewPipelineGroups([]byte(dashboardJSON))

	if err != nil {
		t.Fatalf("Expected no error when creating pipeline groups from valid JSON, but was: %s", err)
	}

	if len(groups) != 1 {
		t.Fatalf("Expected 1 groups, but has %d groups: groups: %#v", len(groups), groups)
	}

	pipelines := groups[0].Pipelines
	if len(pipelines) != 1 {
		t.Fatalf("Expected first groups to have 1 pipeline, but had %d pipelines: pipelines: %#v", len(pipelines), pipelines)
	}

	pipeline := pipelines[0]
	if pipeline.Name != "Pipeline" {
		t.Errorf("Expected pipeline to have proper name, but was: %s", pipeline.Name)
	}

	instances := pipeline.Instances
	if len(instances) != 1 {
		t.Fatalf("Expected pipeline to have 1 instance, but had %d instances: instances: %#v", len(instances), instances)
	}

	stages := instances[0].Stages
	if len(stages) != 2 {
		t.Fatalf("Expected instance to have 2 stages, but had %d stages: stages: %#v", len(stages), stages)
	}

	stage0 := stages[0]
	if stage0.Name != "StageOne" || stage0.Status != "Passed" {
		t.Errorf("Expected first stage to have proper name and status, but was: %#v", stage0)
	}

	stage1 := stages[1]
	if stage1.Name != "StageTwo" || stage1.Status != "Building" {
		t.Errorf("Expected second stage to have proper name and status, but was: %#v", stage1)
	}

	previousInstance := pipeline.PreviousInstance
	if previousInstance.Result != "Passed" {
		t.Errorf("Expected proper previous pipeline instance, but was: %#v", previousInstance)
	}
}

func TestNewPipelineGroupsOnError(t *testing.T) {
	groups, err := gocd.NewPipelineGroups([]byte(`Random`))

	if err == nil {
		t.Fatalf("Expected error when creating pipeline groups from malformed JSON, but was: %s", err)
	}

	if !strings.Contains(err.Error(), "error unmarshalling Gocd JSON: ") {
		t.Errorf("Expected error message about JSON unmarshall error, but was: %s", err.Error())
	}

	if groups != nil {
		t.Fatalf("Expected no invalid groups, but was: %#v", groups)
	}
}
