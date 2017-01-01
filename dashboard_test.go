// dashboard_test.go
//
// Author::    Chirantan Mitra
// Copyright:: Copyright (c) 2015-2017. All rights reserved
// License::   MIT

package gocd_test

import (
	"reflect"
	"testing"

	"github.com/chiku/gocd"
)

func TestDashboardToJSON(t *testing.T) {
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

	if err != nil {
		t.Fatalf("Expected no error marshalling dashboard to JSON: %s", err)
	}

	if string(body) != `[{"name":"Pipeline","stages":[{"name":"Stage","status":"Passed"}]}]` {
		t.Errorf("Expected valid JSON output, but was: %s", body)
	}
}

func TestDashboardFilteredSort(t *testing.T) {
	s1 := []gocd.DashboardStage{gocd.DashboardStage{Name: "Stage One", Status: "Passed"}}
	s2 := []gocd.DashboardStage{gocd.DashboardStage{Name: "Stage Two", Status: "Passed"}}
	s3 := []gocd.DashboardStage{gocd.DashboardStage{Name: "Stage Three", Status: "Passed"}}
	p1 := gocd.DashboardPipeline{Name: "Pipeline One", Stages: s1}
	p2 := gocd.DashboardPipeline{Name: "Pipeline Two", Stages: s2}
	p3 := gocd.DashboardPipeline{Name: "Pipeline Three", Stages: s3}
	dashboard := gocd.Dashboard{p2, p1, p3}

	order := []string{"Pipeline One", "Pipeline Two", "Pipeline Three"}
	sortedDashboard, ignores := dashboard.FilteredSort(order)

	if len(sortedDashboard) != 3 {
		t.Fatalf("Expected sorted dashboard to have 3 entries, but it had %d entries: dashboard: %#v", len(sortedDashboard), sortedDashboard)
	}

	if len(ignores) != 0 {
		t.Fatalf("Expected nothing to be ignored, but %d were ignored: ignores: %#v", len(ignores), ignores)
	}

	item0 := sortedDashboard[0]
	if item0.Name != "Pipeline One" || !reflect.DeepEqual(item0.Stages, s1) {
		t.Errorf("Expected first stage to be proper, but was: %#v", item0)
	}

	item1 := sortedDashboard[1]
	if item1.Name != "Pipeline Two" || !reflect.DeepEqual(item1.Stages, s2) {
		t.Errorf("Expected second stage to be proper, but was: %#v", item1)
	}

	item2 := sortedDashboard[2]
	if item2.Name != "Pipeline Three" || !reflect.DeepEqual(item2.Stages, s3) {
		t.Errorf("Expected third stage to be proper, but was: %#v", item2)
	}
}

func TestDashboardFilteredSortWithIgnores(t *testing.T) {
	s1 := []gocd.DashboardStage{gocd.DashboardStage{Name: "Stage One", Status: "Passed"}}
	s2 := []gocd.DashboardStage{gocd.DashboardStage{Name: "Stage Two", Status: "Passed"}}
	s3 := []gocd.DashboardStage{gocd.DashboardStage{Name: "Stage Three", Status: "Passed"}}
	p1 := gocd.DashboardPipeline{Name: "Pipeline One", Stages: s1}
	p2 := gocd.DashboardPipeline{Name: "Pipeline Two", Stages: s2}
	p3 := gocd.DashboardPipeline{Name: "Pipeline Three", Stages: s3}
	dashboard := gocd.Dashboard{p2, p1, p3}

	order := []string{"Pipeline One", "Pipeline Three"}
	sortedDashboard, ignores := dashboard.FilteredSort(order)

	if len(sortedDashboard) != 2 {
		t.Fatalf("Expected sorted dashboard to have 2 entries, but it had %d entries: dashboard: %#v", len(sortedDashboard), sortedDashboard)
	}

	if len(ignores) != 1 {
		t.Fatalf("Expected one pipeline to be ignored, but %d were ignored: ignores: %#v", len(ignores), ignores)
	}

	item0 := sortedDashboard[0]
	if item0.Name != "Pipeline One" || !reflect.DeepEqual(item0.Stages, s1) {
		t.Errorf("Expected first stage to be proper, but was: %#v", item0)
	}

	item1 := sortedDashboard[1]
	if item1.Name != "Pipeline Three" || !reflect.DeepEqual(item1.Stages, s3) {
		t.Errorf("Expected second stage to be proper, but was: %#v", item1)
	}

	if ignores[0] != "Pipeline Two" {
		t.Fatalf("Expected proper ignores, but was: %#v", ignores[0])
	}
}

func TestDashboardFilteredSortWithDifferentCase(t *testing.T) {
	s1 := []gocd.DashboardStage{gocd.DashboardStage{Name: "Stage One", Status: "Passed"}}
	s2 := []gocd.DashboardStage{gocd.DashboardStage{Name: "Stage Two", Status: "Passed"}}
	s3 := []gocd.DashboardStage{gocd.DashboardStage{Name: "Stage Three", Status: "Passed"}}
	p1 := gocd.DashboardPipeline{Name: "Pipeline One", Stages: s1}
	p2 := gocd.DashboardPipeline{Name: "Pipeline Two", Stages: s2}
	p3 := gocd.DashboardPipeline{Name: "Pipeline Three", Stages: s3}
	dashboard := gocd.Dashboard{p2, p1, p3}

	order := []string{"pipeline one", "pipeline three"}
	sortedDashboard, ignores := dashboard.FilteredSort(order)

	if len(sortedDashboard) != 2 {
		t.Fatalf("Expected sorted dashboard to have 2 entries, but it had %d entries: dashboard: %#v", len(sortedDashboard), sortedDashboard)
	}

	if len(ignores) != 1 {
		t.Fatalf("Expected one pipeline to be ignored, but %d were ignored: ignores: %#v", len(ignores), ignores)
	}

	item0 := sortedDashboard[0]
	if item0.Name != "Pipeline One" || !reflect.DeepEqual(item0.Stages, s1) {
		t.Errorf("Expected first stage to be proper, but was: %#v", item0)
	}

	item1 := sortedDashboard[1]
	if item1.Name != "Pipeline Three" || !reflect.DeepEqual(item1.Stages, s3) {
		t.Errorf("Expected second stage to be proper, but was: %#v", item1)
	}

	if ignores[0] != "Pipeline Two" {
		t.Fatalf("Expected proper ignores, but was: %#v", ignores[0])
	}
}

func TestDashboardFilteredSortWithUnknownPipelines(t *testing.T) {
	s1 := []gocd.DashboardStage{gocd.DashboardStage{Name: "Stage One", Status: "Passed"}}
	s2 := []gocd.DashboardStage{gocd.DashboardStage{Name: "Stage Two", Status: "Passed"}}
	s3 := []gocd.DashboardStage{gocd.DashboardStage{Name: "Stage Three", Status: "Passed"}}
	p1 := gocd.DashboardPipeline{Name: "Pipeline One", Stages: s1}
	p2 := gocd.DashboardPipeline{Name: "Pipeline Two", Stages: s2}
	p3 := gocd.DashboardPipeline{Name: "Pipeline Three", Stages: s3}
	dashboard := gocd.Dashboard{p2, p1, p3}

	order := []string{"Pipeline One", "Pipeline Two", "Pipeline Four"}
	sortedDashboard, ignores := dashboard.FilteredSort(order)

	if len(sortedDashboard) != 2 {
		t.Fatalf("Expected sorted dashboard to have 2 entries, but it had %d entries: dashboard: %#v", len(sortedDashboard), sortedDashboard)
	}

	if len(ignores) != 1 {
		t.Fatalf("Expected one pipeline to be ignored, but %d were ignored: ignores: %#v", len(ignores), ignores)
	}

	item0 := sortedDashboard[0]
	if item0.Name != "Pipeline One" || !reflect.DeepEqual(item0.Stages, s1) {
		t.Errorf("Expected first stage to be proper, but was: %#v", item0)
	}

	item1 := sortedDashboard[1]
	if item1.Name != "Pipeline Two" || !reflect.DeepEqual(item1.Stages, s2) {
		t.Errorf("Expected second stage to be proper, but was: %#v", item1)
	}

	if ignores[0] != "Pipeline Three" {
		t.Fatalf("Expected proper ignores, but was: %#v", ignores[0])
	}
}

func TestDashboardMapNames(t *testing.T) {
	s1 := []gocd.DashboardStage{gocd.DashboardStage{Name: "Stage One", Status: "Passed"}}
	s2 := []gocd.DashboardStage{gocd.DashboardStage{Name: "Stage Two", Status: "Passed"}}
	p1 := gocd.DashboardPipeline{Name: "Pipeline One", Stages: s1}
	p2 := gocd.DashboardPipeline{Name: "Pipeline Two", Stages: s2}
	dashboard := gocd.Dashboard{p1, p2}

	mapping := map[string]string{
		"Pipeline One": "Pipeline A",
		"Pipeline Two": "Pipeline B",
	}
	mappedDashboard := dashboard.MapNames(mapping)

	if len(mappedDashboard) != 2 {
		t.Fatalf("Expected mapped dashboard to have 2 entries, but it had %d entries: dashboard: %#v", len(mappedDashboard), mappedDashboard)
	}

	item0 := mappedDashboard[0]
	if item0.Name != "Pipeline A" || !reflect.DeepEqual(item0.Stages, s1) {
		t.Errorf("Expected first stage to have new name, but was: %#v", item0)
	}

	item1 := mappedDashboard[1]
	if item1.Name != "Pipeline B" || !reflect.DeepEqual(item1.Stages, s2) {
		t.Errorf("Expected second stage to have new name, but was: %#v", item1)
	}
}

func TestDashboardMapNamesWhenNoMatch(t *testing.T) {
	s1 := []gocd.DashboardStage{gocd.DashboardStage{Name: "Stage One", Status: "Passed"}}
	s2 := []gocd.DashboardStage{gocd.DashboardStage{Name: "Stage Two", Status: "Passed"}}
	p1 := gocd.DashboardPipeline{Name: "Pipeline One", Stages: s1}
	p2 := gocd.DashboardPipeline{Name: "Pipeline Two", Stages: s2}
	dashboard := gocd.Dashboard{p1, p2}

	mapping := map[string]string{}
	mappedDashboard := dashboard.MapNames(mapping)

	if len(mappedDashboard) != 2 {
		t.Fatalf("Expected mapped dashboard to have 2 entries, but it had %d entries: dashboard: %#v", len(mappedDashboard), mappedDashboard)
	}

	item0 := mappedDashboard[0]
	if item0.Name != "Pipeline One" || !reflect.DeepEqual(item0.Stages, s1) {
		t.Errorf("Expected first stage to have original name, but was: %#v", item0)
	}

	item1 := mappedDashboard[1]
	if item1.Name != "Pipeline Two" || !reflect.DeepEqual(item1.Stages, s2) {
		t.Errorf("Expected second stage to have original name, but was: %#v", item1)
	}
}
