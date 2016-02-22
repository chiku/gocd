// gocd.go
//
// Author::    Chirantan Mitra
// Copyright:: Copyright (c) 2015-2016. All rights reserved
// License::   MIT

package gocd

import (
	"encoding/json"
	"fmt"
	"strings"
)

const (
	building   = "Building"
	unknown    = "Unknown"
	failed     = "Failed"
	recovering = "Recovering"
)

type Stage struct {
	Name   string `json:"name"`
	Status string `json:"status"`
}
type Instance struct {
	Stages []Stage `json:"stages"`
}
type PreviousInstance struct {
	Result string `json:"result"`
}
type Pipeline struct {
	Name             string           `json:"name"`
	Instances        []Instance       `json:"instances"`
	PreviousInstance PreviousInstance `json:"previous_instance"`
}
type PipelineGroup struct {
	Pipelines []Pipeline `json:"pipelines"`
}
type PipelineGroups []PipelineGroup

func NewPipelineGroups(body []byte) (PipelineGroups, error) {
	var dashboard []PipelineGroup
	err := json.Unmarshal(body, &dashboard)
	if err != nil {
		return nil, fmt.Errorf("error unmarshalling Gocd JSON: %s", err.Error())
	}
	return dashboard, nil
}

func (groups *PipelineGroups) ToDashboard() Dashboard {
	dashboard := Dashboard{}

	for _, group := range *groups {
		for _, pipeline := range group.Pipelines {
			displayName := pipeline.Name
			stages := []DashboardStage{}
			if len(pipeline.Instances) > 0 {

				instance := pipeline.Instances[len(pipeline.Instances)-1]
				for _, stage := range instance.Stages {
					status := traverseStatusInInstances(stage, pipeline.Instances, pipeline.PreviousInstance)
					stages = append(stages, DashboardStage{Name: stage.Name, Status: status})
				}
				if len(stages) > 0 {
					dashboard = append(dashboard, DashboardPipeline{Name: displayName, Stages: stages})
				}
			}
		}
	}

	return dashboard
}

func traverseStatusInInstances(currentStage Stage, instances []Instance, previousInstance PreviousInstance) string {
	selfStatus := currentStage.Status
	previousInstanceResult := previousInstance.Result

	if previousInstanceResult == failed && strings.EqualFold(selfStatus, building) {
		return recovering
	}

	if !strings.EqualFold(selfStatus, unknown) {
		return selfStatus
	}

	olderInstances := instances[0 : len(instances)-1]
	olderInstanceStatus := findKnownStatusInInstances(currentStage, olderInstances)
	if !strings.EqualFold(olderInstanceStatus, unknown) {
		return olderInstanceStatus
	}

	if previousInstanceResult != "" && !strings.EqualFold(previousInstanceResult, unknown) {
		return previousInstanceResult
	}

	return unknown
}

func findKnownStatusInInstances(currentStage Stage, instances []Instance) string {
	for i := len(instances) - 1; i >= 0; i-- {
		instance := instances[i]
		for j := len(instance.Stages) - 1; j >= 0; j-- {
			stage := instance.Stages[j]
			if currentStage.Name == stage.Name && !strings.EqualFold(stage.Status, unknown) {
				return stage.Status
			}
		}
	}

	return unknown
}
