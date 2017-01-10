// gocd.go
//
// Author::    Chirantan Mitra
// Copyright:: Copyright (c) 2015-2017. All rights reserved
// License::   MIT

package gocd

import "encoding/json"

func Fetch() func(string, []string, map[string]string) ([]byte, []string, error) {
	client := NewClient()

	return func(url string, filters []string, transforms map[string]string) (output []byte, ignores []string, err error) {
		dashboard, err := client.Fetch(url)
		if err != nil {
			return nil, nil, err
		}

		dashboard, ignores = dashboard.FilteredSort(filters)
		dashboard = dashboard.MapNames(transforms)

		output, err = json.Marshal(dashboard)
		if err != nil {
			return nil, nil, err
		}

		return output, ignores, nil
	}
}
