// gocd_test.go
//
// Author::    Chirantan Mitra
// Copyright:: Copyright (c) 2015-2017. All rights reserved
// License::   MIT

package gocd_test

import (
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"

	"github.com/chiku/gocd"
)

func TestFetch(t *testing.T) {
	const serverResponse = `[{
		"name": "Group",
		"pipelines": [{
		    "name": "pipeline1",
		    "instances": [{ "stages": [{ "name": "StageOne", "status": "Passed" }] }]
	    }, {
		    "name": "pipeline2",
		    "instances": [{ "stages": [{ "name": "StageOne", "status": "Failed" }] }]
	    }, {
		    "name": "pipeline3",
		    "instances": [{ "stages": [{ "name": "StageOne", "status": "Building" }] }]
	    }]
	}]`
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(serverResponse))
	}))
	defer ts.Close()

	filters := []string{"pipeline1", "pipeline2"}
	transforms := map[string]string{"pipeline1": "p1"}

	fetcher := gocd.Fetch()
	output, ignores, err := fetcher(ts.URL, filters, transforms)

	if err != nil {
		t.Fatalf("Expected no error fetching valid response: %s", err)
	}

	expectedOutput := []byte(`[{"name":"p1","stages":[{"name":"StageOne","status":"Passed"}]},{"name":"pipeline2","stages":[{"name":"StageOne","status":"Failed"}]}]`)
	if !reflect.DeepEqual(output, expectedOutput) {
		t.Errorf("Incorrect output JSON (%s != %s)", output, expectedOutput)
	}

	expectedIgnores := []string{"pipeline3"}
	if !reflect.DeepEqual(ignores, expectedIgnores) {
		t.Errorf("Incorrect ignore (%s != %s)", ignores, expectedIgnores)
	}
}

func TestFetchWhenServerError(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusForbidden)
		w.Write([]byte("Forbidden!"))
	}))
	defer ts.Close()

	filters := []string{"pipeline1", "pipeline2"}
	transforms := map[string]string{"pipeline1": "p1"}

	fetcher := gocd.Fetch()
	output, ignores, err := fetcher(ts.URL, filters, transforms)

	if err == nil {
		t.Fatalf("Expected error fetching invalid response: %s", err)
	}

	if err.Error() != "error fetching response from Gocd: the HTTP status code was 403, body: Forbidden!" {
		t.Errorf("Expected proper error message but was: %s", err.Error())
	}

	if output != nil || ignores != nil {
		t.Errorf("Expected no output but output=%s and ignores=%v", output, ignores)
	}
}
