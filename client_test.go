// client_test.go
//
// Author::    Chirantan Mitra
// Copyright:: Copyright (c) 2015-2017. All rights reserved
// License::   MIT

package gocd_test

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/chiku/gocd"
)

func TestClientFetch(t *testing.T) {
	const serverResponse = `[{
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
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(serverResponse))
	}))
	defer ts.Close()

	client := gocd.NewClient()
	dashboard, err := client.Fetch(ts.URL)

	if err != nil {
		t.Fatalf("Expected no error fetching valid response: %s", err)
	}

	if len(dashboard) != 1 {
		t.Fatalf("Expected dashboard to contain 1 item, but it had %d items: dashboard: %#v", len(dashboard), dashboard)
	}

	pipeline := dashboard[0]
	if pipeline.Name != "Pipeline" {
		t.Errorf("Expected proper pipeline name, but was: %v", pipeline.Name)
	}

	stages := pipeline.Stages
	if len(stages) != 2 {
		t.Fatalf("Expected stages to contain 2 items, but it had %d items: stages: %#v", len(stages), stages)
	}

	stage0 := stages[0]
	if stage0.Name != "StageOne" || stage0.Status != "Passed" {
		t.Errorf("Expected first stage to be proper, but was: %#v", stage0)
	}

	stage1 := stages[1]
	if stage1.Name != "StageTwo" || stage1.Status != "Building" {
		t.Errorf("Expected second stage to be proper, but was: %#v", stage1)
	}
}

func TestClientFetchWhenServerResponseNot200(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusForbidden)
		w.Write([]byte("Forbidden!"))
	}))
	defer ts.Close()

	client := gocd.NewClient()
	dashboard, err := client.Fetch(ts.URL)

	if err == nil {
		t.Fatalf("Expected error fetching invalid response: %s", err)
	}
	if err.Error() != "error fetching response from Gocd: the HTTP status code was 403, body: Forbidden!" {
		t.Errorf("Expected proper error message but was: %s", err.Error())
	}

	if dashboard != nil {
		t.Errorf("Expected no invalid dashboard, but was: %#v", dashboard)
	}
}

func TestClientFetchWhenServerResponseMalformed(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Bad response"))
	}))
	defer ts.Close()

	client := gocd.NewClient()
	dashboard, err := client.Fetch(ts.URL)

	if err == nil {
		t.Fatalf("Expected error fetching invalid response: %s", err)
	}
	if !strings.Contains(err.Error(), "error unmarshalling Gocd JSON: ") {
		t.Errorf("Expected proper error message but was: %s", err.Error())
	}

	if dashboard != nil {
		t.Errorf("Expected no invalid dashboard, but was: %#v", dashboard)
	}
}

func TestClientFetchWhenServerDoesNotRespondAfterRetries(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Bad response"))
	}))
	ts.Close()

	client := gocd.NewClient()
	dashboard, err := client.Fetch(ts.URL)

	if err == nil {
		t.Fatalf("Expected error fetching invalid response: %s", err)
	}
	if !strings.Contains(err.Error(), "error fetching data from Gocd: ") || !strings.Contains(err.Error(), "(after 3 retries)") {
		t.Errorf("Expected proper error message with retries but was: %s", err.Error())
	}

	if dashboard != nil {
		t.Errorf("Expected no invalid dashboard, but was: %#v", dashboard)
	}
}

func TestClientFetchWhenServerDoesNotRespond(t *testing.T) {
	client := gocd.NewClient()
	dashboard, err := client.Fetch("<>")

	if err == nil {
		t.Fatalf("Expected error fetching invalid response: %s", err)
	}
	if !strings.Contains(err.Error(), "error fetching data from Gocd: ") || !strings.Contains(err.Error(), "(after 3 retries)") {
		t.Errorf("Expected proper error message with retries but was: %s", err.Error())
	}

	if dashboard != nil {
		t.Errorf("Expected no invalid dashboard, but was: %#v", dashboard)
	}
}

func TestClientFetchWhenRequestCreationFails(t *testing.T) {
	client := gocd.NewClient()
	dashboard, err := client.Fetch("::")

	if err == nil {
		t.Fatalf("Expected error fetching invalid response: %s", err)
	}
	if !strings.Contains(err.Error(), "error creating Gocd request: ") {
		t.Errorf("Expected proper error message but was: %s", err.Error())
	}

	if dashboard != nil {
		t.Errorf("Expected no invalid dashboard, but was: %#v", dashboard)
	}
}
