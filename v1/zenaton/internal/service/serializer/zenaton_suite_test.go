package serializer

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"testing"
)

func TestZenaton(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Zenaton Suite")
}
