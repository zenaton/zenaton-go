package errors_test

import (
	. "github.com/zenaton/zenaton-go/v1/zenaton/errors"

	"strings"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/pkg/errors"
	"github.com/zenaton/zenaton-go/v1/zenaton"
)

var _ = Describe("Errors", func() {
	Context("New", func() {
		It("should create a new zenaton error", func() {
			err := New("testName", "testMessage")
			Expect(err.Name()).To(Equal("testName"))
			Expect(err.Error()).To(Equal("testMessage"))

			trace := err.Trace()
			tracePieces := strings.Fields(trace)
			secondLine := tracePieces[1]
			Expect(secondLine).To(ContainSubstring("zenaton/errors/errors_test.go"))
		})
	})

	Context("Wrap", func() {
		It("should create a new zenaton error from an error", func() {
			err := Wrap("testName", errors.New("testMessage"))
			Expect(err.Name()).To(Equal("testName"))
			Expect(err.Error()).To(Equal("testMessage"))

			trace := err.Trace()
			tracePieces := strings.Fields(trace)
			secondLine := tracePieces[1]
			Expect(secondLine).To(ContainSubstring("zenaton/errors/errors_test.go"))
		})
	})

	Context("New from zenaton service", func() {
		It("should create a new zenaton error", func() {
			err := zenaton.NewService().Errors.New("testName", "testMessage")
			Expect(err.Name()).To(Equal("testName"))
			Expect(err.Error()).To(Equal("testMessage"))

			trace := err.Trace()
			tracePieces := strings.Fields(trace)
			secondLine := tracePieces[1]

			Expect(secondLine).To(ContainSubstring("zenaton/errors/errors_test.go"))
		})
	})

	Context("Wrap from zenaton service", func() {
		It("should create a new zenaton error from an error", func() {
			err := zenaton.NewService().Errors.Wrap("testName", errors.New("testMessage"))
			Expect(err.Name()).To(Equal("testName"))
			Expect(err.Error()).To(Equal("testMessage"))

			trace := err.Trace()
			tracePieces := strings.Fields(trace)
			secondLine := tracePieces[1]
			Expect(secondLine).To(ContainSubstring("zenaton/errors/errors_test.go"))
		})
	})
})
