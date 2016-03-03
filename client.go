// dashboard.go
//
// Author::    Chirantan Mitra
// Copyright:: Copyright (c) 2015-2016. All rights reserved
// License::   MIT

package gocd

import (
	"fmt"
	"io/ioutil"
	"net/http"
)

const (
	maxRetries = 3
)

func fetchGocdDashboard(client *http.Client, request *http.Request) (response *http.Response, err error) {
	retries := 0

	for response == nil && retries < maxRetries {
		response, err = client.Do(request)
		if err != nil {
			err = fmt.Errorf("%s (after %d retries)", err, retries+1)
		}
		retries++
	}

	return
}

func parseHTTPResponse(response *http.Response) (PipelineGroups, error) {
	if response != nil {
		defer response.Body.Close()
	}

	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return nil, fmt.Errorf("error reading response: %s", err)
	}

	if response.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("error fetching response from Gocd: the HTTP status code was %d, body: %s", response.StatusCode, body)
	}

	groups, err := NewPipelineGroups(body)
	if err != nil {
		return nil, err
	}

	return groups, nil
}

type Client struct {
	client *http.Client
}

func NewClient() *Client {
	return &Client{client: &http.Client{}}
}

func (c *Client) Fetch(url string) (Dashboard, error) {
	request, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("error creating Gocd request: %s", err)
	}

	response, err := fetchGocdDashboard(c.client, request)
	if err != nil {
		return nil, fmt.Errorf("error fetching data from Gocd: %s", err)
	}

	groups, err := parseHTTPResponse(response)
	if err != nil {
		return nil, err
	}

	dashboard := groups.ToDashboard()
	return dashboard, nil
}
