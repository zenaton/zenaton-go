package zenaton_test

import (
	"github.com/davecgh/go-spew/spew"
	. "github.com/onsi/ginkgo"
	//. "github.com/onsi/gomega"
)

type Person struct {
	Parent *Person `json:"parent"`
	Child  *Person `json:"child"`
}

var _ = Describe("Serializer", func() {
	Context("When using recursive structures", func() {
		It("Should still serialize and unserialize", func() {
			grandFather := &Person{}
			father := &Person{
				Parent: grandFather,
			}
			son := &Person{
				Parent: father,
			}
			father.Child = son
			grandFather.Child = father

			spew.Dump(grandFather)
		})
	})
})
