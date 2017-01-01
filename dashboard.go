// dashboard.go
//
// Author::    Chirantan Mitra
// Copyright:: Copyright (c) 2015-2017. All rights reserved
// License::   MIT

package gocd

import (
	"encoding/json"
	"fmt"
	"strings"
)

type DashboardStage struct {
	Name   string `json:"name"`
	Status string `json:"status"`
}
type DashboardPipeline struct {
	Name   string           `json:"name"`
	Stages []DashboardStage `json:"stages"`
	order  int
}
type Dashboard []DashboardPipeline

func (dashboard Dashboard) ToJSON() (output []byte, err error) {
	output, err = json.Marshal(dashboard)
	if err != nil {
		return nil, fmt.Errorf("error marshalling simple dashboard JSON :%s", err.Error())
	}

	return output, nil
}

func (dashboard Dashboard) FilteredSort(order []string) (sortedDashboard Dashboard, ignores []string) {
	for _, o := range order {
		found := dashboard.findPipelineWithName(o)
		if found != nil {
			sortedDashboard = append(sortedDashboard, *found)
		}
	}

	for _, pipeline := range dashboard {
		if !isStringInsideSlice(order, pipeline.Name) {
			ignores = append(ignores, pipeline.Name)
		}
	}

	return
}

func (dashboard Dashboard) MapNames(mapping map[string]string) (mappedDashboard Dashboard) {
	for _, pipeline := range dashboard {
		if val, ok := mapping[pipeline.Name]; ok {
			mappedDashboard = append(mappedDashboard, DashboardPipeline{Name: val, Stages: pipeline.Stages})
		} else {
			mappedDashboard = append(mappedDashboard, pipeline)
		}
	}

	return
}

func (dashboard Dashboard) findPipelineWithName(name string) *DashboardPipeline {
	for _, pipeline := range dashboard {
		if strings.EqualFold(pipeline.Name, name) {
			return &pipeline
		}
	}

	return nil
}

func isStringInsideSlice(slice []string, item string) bool {
	for _, s := range slice {
		if strings.EqualFold(s, item) {
			return true
		}
	}

	return false
}
