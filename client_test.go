// dashboard.go
//
// Author::    Chirantan Mitra
// Copyright:: Copyright (c) 2015-2017. All rights reserved
// License::   MIT

package gocd_test

import (
	"net/http"
	"net/http/httptest"

	"github.com/chiku/gocd"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Client", func() {
	Context("fetching dashboard.json from gocd", func() {
		Context("when success", func() {
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

			It("returns dashboard", func() {
				Expect(dashboard).To(HaveLen(1))
				pipeline := dashboard[0]
				Expect(pipeline.Name).To(Equal("Pipeline"))
				stages := pipeline.Stages
				Expect(stages).To(HaveLen(2))
				Expect(stages[0].Name).To(Equal("StageOne"))
				Expect(stages[0].Status).To(Equal("Passed"))
				Expect(stages[1].Name).To(Equal("StageTwo"))
				Expect(stages[1].Status).To(Equal("Building"))
			})

			It("has no errors", func() {
				Expect(err).To(Succeed())
			})
		})

		Context("when server response in not 200 OK", func() {
			ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusForbidden)
				w.Write([]byte("Forbidden!"))
			}))
			defer ts.Close()

			client := gocd.NewClient()
			dashboard, err := client.Fetch(ts.URL)

			It("has no dashboard", func() {
				Expect(dashboard).To(BeNil())
			})

			It("has errors", func() {
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(Equal("error fetching response from Gocd: the HTTP status code was 403, body: Forbidden!"))
			})
		})

		Context("when server response body is malformed", func() {
			ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.Write([]byte("Bad response"))
			}))
			defer ts.Close()

			client := gocd.NewClient()
			dashboard, err := client.Fetch(ts.URL)

			It("has no dashboard", func() {
				Expect(dashboard).To(BeNil())
			})

			It("has errors", func() {
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("error unmarshalling Gocd JSON: "))
			})
		})

		Context("when server doesn't respond", func() {
			ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.Write([]byte("Bad response"))
			}))
			ts.Close()

			client := gocd.NewClient()
			dashboard, err := client.Fetch(ts.URL)

			It("has no dashboard", func() {
				Expect(dashboard).To(BeNil())
			})

			It("has errors", func() {
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("error fetching data from Gocd: "))
				Expect(err.Error()).To(ContainSubstring("(after 3 retries)"))
			})
		})

		Context("when <>", func() {
			client := gocd.NewClient()
			dashboard, err := client.Fetch("<>")

			It("has no dashboard", func() {
				Expect(dashboard).To(BeNil())
			})

			It("has errors", func() {
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("error fetching data from Gocd: "))
				Expect(err.Error()).To(ContainSubstring("(after 3 retries)"))
			})
		})

		Context("when request creation fails", func() {
			client := gocd.NewClient()
			dashboard, err := client.Fetch("::")

			It("has no dashboard", func() {
				Expect(dashboard).To(BeNil())
			})

			It("has errors", func() {
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("error creating Gocd request: "))
			})
		})
	})
})
