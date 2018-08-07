package zenaton_test

import (
	"encoding/json"

	. "github.com/onsi/ginkgo"

	"fmt"

	. "github.com/onsi/gomega"
	"github.com/zenaton/zenaton-go/v1/zenaton"
)

type Person struct {
	Parent *Person `json:"parent"`
	Child  *Person `json:"child"`
}

type SerializeMe struct {
	Initialized bool
}

var _ = Describe("Encode", func() {
	s := zenaton.Serializer{}
	Context("with a string", func() {
		It("represents the string as a data", func() {
			encoded, err := s.Encode("e")
			Expect(err).ToNot(HaveOccurred())
			var value map[string]interface{}
			json.Unmarshal([]byte(encoded), &value)
			Expect(value["d"]).To(Equal("e"))
			Expect(value["s"]).To(Equal([]interface{}{}))
		})
	})

	Context("with a string", func() {
		It("represents the integer as a data", func() {
			encoded, err := s.Encode(1)
			Expect(err).ToNot(HaveOccurred())

			fmt.Println("encoded: ", encoded)
			Expect(encoded).To(Equal(encoded))
			//Expect(value["s"]).To(Equal([]interface{}{}))
		})
	})

	Context("with a float", func() {
		It("represents the integer as a data", func() {
			encoded, err := s.Encode(1.8)
			Expect(err).ToNot(HaveOccurred())
			fmt.Println(encoded)
			Expect(encoded).To(Equal(`{"d":1.8,"s":[]}`))
		})
	})

	Context("with true", func() {
		It("represents the boolean as a data", func() {
			encoded, err := s.Encode(true)
			Expect(err).ToNot(HaveOccurred())
			fmt.Println(encoded)
			Expect(encoded).To(Equal(`{"d":true,"s":[]}`))
		})
	})

	Context("with false", func() {
		It("represents the boolean as a data", func() {
			encoded, err := s.Encode(false)
			Expect(err).ToNot(HaveOccurred())
			fmt.Println(encoded)
			Expect(encoded).To(Equal(`{"d":false,"s":[]}`))
		})
	})

	Context("with nil", func() {
		It("represents the nil as a data", func() {
			var nilInterface interface{}
			encoded, err := s.Encode(nilInterface)
			Expect(err).ToNot(HaveOccurred())
			fmt.Println(encoded)
			Expect(encoded).To(Equal(`{"d":null,"s":[]}`))
		})
	})

	Context("with a simple array", func() {
		XIt("represents the array as an array", func() {
			arr := []interface{}{1, "e"}
			encoded, err := s.Encode(arr)
			Expect(err).ToNot(HaveOccurred())
			fmt.Println(encoded)
			Expect(encoded).To(Equal(`{"a":[1,"e"],"s":[]}`))
		})
	})

	Context("with recursive arrays", func() {
		XIt("represents the array as an array", func() {
			var arr1 []interface{}
			var arr2 []interface{}
			arr1 = append(arr1, &arr2)
			arr2 = append(arr2, &arr1)

			encoded, err := s.Encode(arr1)
			Expect(err).ToNot(HaveOccurred())
			fmt.Println(encoded)
			Expect(encoded).To(Equal(`{"a":[1,"e"],"s":[]}`))
		})
	})

	//todo: this needs to be different, as maps can have pointer keys and values
	Context("with an map", func() {
		It("represents the map as an map", func() {
			m := map[string]string{
				"key": "value",
			}
			encoded, err := s.Encode(m)
			Expect(err).ToNot(HaveOccurred())
			fmt.Println(encoded)
			Expect(encoded).To(Equal(`{"a":{"key":"value"},"s":[]}`))
		})
	})

	Context("with a simple struct", func() {
		It("represents the struct as an object", func() {
			sm := SerializeMe{
				Initialized: true,
			}

			encoded, err := s.Encode(sm)
			Expect(err).ToNot(HaveOccurred())
			fmt.Println(encoded)
			Expect(encoded).To(Equal(`{"o":"@zenaton#0","s":[{"n":"SerializeMe","p":{"Initialized":true}}]}`))
		})
	})

	type Person struct {
		Parent *Person
		Child  *Person
	}

	Context("with a struct with circular dependencies", func() {
		It("represents the struct as an object", func() {

			expectedSerialized := `{"o":"@zenaton#0","s":[{"n":"Person","p":{"Child":null,"Parent":"@zenaton#0"}},{"n":"Person","p":{"Child":"@zenaton#1","Parent":null}}]}`
			parent := Person{}
			child := Person{
				Parent: &parent,
			}
			parent.Child = &child

			encoded, err := s.Encode(parent)
			Expect(err).ToNot(HaveOccurred())
			Expect(encoded).To(Equal(expectedSerialized))
		})
	})

	Describe("decode", func() {
		s := &zenaton.Serializer{}
		Context("with a simple object", func() {
			XIt("returns an instance with the correct instance variables", func() {
				//type MyInt struct {
				//	Int int8
				//}
				//
				//myInt := MyInt{5}
				//
				//jsonMyInt, err := json.Marshal(myInt)
				//Expect(err).NotTo(HaveOccurred())
				//
				//var myInt2 MyInt
				//err = json.Unmarshal(jsonMyInt, &myInt2)
				//Expect(err).NotTo(HaveOccurred())

				str := `{"o":"@zenaton#0","s":[{"n":"SerializeMe","p":{"Initialized":true}}]}`
				var value SerializeMe
				err := s.Decode(str, &value)
				Expect(err).ToNot(HaveOccurred())
				Expect(value).To(Equal(SerializeMe{true}))
			})
		})
	})
})
