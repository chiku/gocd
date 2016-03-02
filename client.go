// dashboard.go
//
// Author::    Chirantan Mitra
// Copyright:: Copyright (c) 2015-2016. All rights reserved
// License::   MIT

package gocd

import (
	"fmt"
	"io/ioutil"
	"log"
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
			log.Printf("error fetching data from Gocd (retry #%d): %s", retries+1, err)
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
	url    string
	client *http.Client
}

func NewClient(url string) *Client {
	return &Client{
		url:    url,
		client: &http.Client{},
	}
}

func (c *Client) Fetch() (Dashboard, error) {
	request, err := http.NewRequest("GET", c.url, nil)
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
