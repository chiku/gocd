// interest_test.go
//
// Author::    Chirantan Mitra
// Copyright:: Copyright (c) 2015-2016. All rights reserved
// License::   MIT

package gocd_test

import (
	"github.com/chiku/gocd"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Interest", func() {
	Context("when without a display name", func() {
		interest := gocd.NewInterest("Foo")

		It("extracts information given exact name", func() {
			match, displayName := interest.PipelineName("Foo")
			Expect(match).To(BeTrue())
			Expect(displayName).To(Equal("Foo"))
		})

		It("extracts information given name with different case", func() {
			match, displayName := interest.PipelineName("foo")
			Expect(match).To(BeTrue())
			Expect(displayName).To(Equal("Foo"))
		})

		It("doesn't extract information for a different name", func() {
			match, displayName := interest.PipelineName("foo1")
			Expect(match).To(BeFalse())
			Expect(displayName).To(BeEmpty())
		})
	})

	Context("when with a display name", func() {
		interest := gocd.NewInterest("Foo:>New Foo")

		It("extracts information given the original name", func() {
			match, displayName := interest.PipelineName("Foo")
			Expect(match).To(BeTrue())
			Expect(displayName).To(Equal("New Foo"))
		})

		It("doesn't extract information given the display name", func() {
			match, displayName := interest.PipelineName("New Foo")
			Expect(match).To(BeFalse())
			Expect(displayName).To(BeEmpty())
		})
	})

	Context("when with an empty name", func() {
		interest := gocd.NewInterest("")

		It("doesn't extract information for empty string", func() {
			match, displayName := interest.PipelineName("")
			Expect(match).To(BeFalse())
			Expect(displayName).To(BeEmpty())
		})
	})
})

var _ = Describe("Interests", func() {
	Context("when created with values", func() {
		interests := gocd.NewInterests().Add("Foo:>A Foo").Add("Bar:>A Bar").Add("Baz")

		Context("when retrieving by a name without a display name", func() {
			position, displayName := interests.PipelineName("Baz")

			It("fetches the position in the interest", func() {
				Expect(position).To(Equal(2))
			})

			It("fetches the name", func() {
				Expect(displayName).To(Equal("Baz"))
			})
		})

		Context("when retrieving by a name with a display name", func() {
			position, displayName := interests.PipelineName("Bar")

			It("fetches the position in the interest", func() {
				Expect(position).To(Equal(1))
			})

			It("fetches the display name", func() {
				Expect(displayName).To(Equal("A Bar"))
			})
		})

		Context("when retrieving by a name not in the interest list", func() {
			position, displayName := interests.PipelineName("NotInterested")

			It("has a position outside the list", func() {
				Expect(position).To(Equal(-1))
			})

			It("doesn't have a name", func() {
				Expect(displayName).To(BeEmpty())
			})
		})
	})

	Context("when created without values", func() {
		interests := gocd.NewInterests().Add("").Add("")

		Context("when retrieving by empty name", func() {
			position, displayName := interests.PipelineName("")

			It("has a position outside the list", func() {
				Expect(position).To(Equal(-1))
			})

			It("doesn't have a name", func() {
				Expect(displayName).To(BeEmpty())
			})
		})
	})
})
