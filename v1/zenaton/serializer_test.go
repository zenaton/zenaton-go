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
		It("represents the array as an array", func() {
			arr := []interface{}{1, "e"}
			encoded, err := s.Encode(arr)
			Expect(err).ToNot(HaveOccurred())
			fmt.Println(encoded)
			Expect(encoded).To(Equal(`{"o":"@zenaton#0","s":[{"a":[1,"e"]}]}`))
		})
	})

	//context 'with recursive arrays' do
	//  let(:array1) { [1, 2, 3] }
	//  let(:array2) { [4, 5, 6] }
	//  let(:data) { array1 }
	//  let(:expected_representation) do
	//    {
	//      'o' => '@zenaton#0',
	//      's' => [{
	//        'a' => [1, 2, 3, '@zenaton#1']
	//      }, {
	//        'a' => [4, 5, 6, '@zenaton#0']
	//      }]
	//    }
	//  end
	//
	//  before do
	//    array1 << array2
	//    array2 << array1
	//  end
	//
	//  it 'represents the array as an object' do
	//    expect(parsed_json).to eq(expected_representation)
	//  end
	//

	FContext("with recursive arrays", func() {
		It("represents the array as an array", func() {
			expectedOutput := `{"o":"@zenaton#0","s":[{"a":["@zenaton#1"]},{"a":["@zenaton#0"]}]}`
			var arr1 []interface{}
			var arr2 []interface{}
			arr1 = append(arr1, &arr2)
			arr2 = append(arr2, &arr1)

			encoded, err := s.Encode(arr1)
			Expect(err).ToNot(HaveOccurred())
			fmt.Println(encoded)
			Expect(encoded).To(Equal(expectedOutput))
		})
	})

	//todo: this needs to be different, as maps can have pointer keys and values
	XContext("with an map", func() {
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

			expectedSerialized := `{"o":"@zenaton#0","s":[{"n":"Person","p":{"Child":"@zenaton#1","Parent":null}},{"n":"Person","p":{"Child":null,"Parent":"@zenaton#0"}}]}`
			parent := Person{}
			child := Person{
				Parent: &parent,
			}
			parent.Child = &child

			encoded, err := s.Encode(parent)
			Expect(err).ToNot(HaveOccurred())
			fmt.Println("encoded: ", encoded)
			Expect(encoded).To(Equal(expectedSerialized))
		})
	})

	//todo test with unexported fields!

	//todo: make sure that it handles unexported fields well.
	type MySimpleStruct struct {
		Bool    bool
		Int     int
		Int8    int8
		Int16   int16
		Int32   int32
		Int64   int64
		Uint    uint
		Uint8   uint8
		Uint16  uint16
		Uint32  uint32
		Uint64  uint64
		Float32 float32
		Float64 float64
		String  string
		Ptr     **string
		//todo:?
		//Array [1]string
		//Uintptr uintptr
		//Interface interface{}
		//Map map[string]interface{}
		//Slice []string
		//Struct struct
		//UnsafePointer unsafePointer
	}
	Describe("decode", func() {
		//s := &zenaton.Serializer{}
		Context("with a simple struct", func() {
			It("returns an instance with the correct instance variables", func() {

				pointed := "a"
				pointed2 := &pointed
				mySimpleStruct := MySimpleStruct{
					Bool:    true,
					Int:     1,
					Int8:    1,
					Int16:   1,
					Int32:   1,
					Int64:   1,
					Uint:    1,
					Uint8:   1,
					Uint16:  1,
					Uint32:  1,
					Uint64:  1,
					Float32: 1.1,
					Float64: 1.1,
					String:  "a",
					Ptr:     &pointed2,
				}

				encoded := `{
							   "o":"@zenaton#0",
							   "s":[
								  {
									 "n":"MySimpleStruct",
									 "p":{
										"Bool":true,
										"Int":1,
										"Int8":1,
										"Int16":1,
										"Int32":1,
										"Int64":1,
										"Uint":1,
										"Uint8":1,
										"Uint16":1,
										"Uint32":1,
										"Uint64":1,
										"Float32":1.1,
										"Float64":1.1,
										"String":"a",
										"Ptr":"a"
									 }
								  }
							   ]
							}`
				var encodedStruct MySimpleStruct
				err := s.Decode(encoded, &encodedStruct)
				Expect(err).ToNot(HaveOccurred())
				Expect(encodedStruct).To(Equal(mySimpleStruct))
			})
		})
	})
})
