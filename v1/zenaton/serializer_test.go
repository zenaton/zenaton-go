package zenaton_test

import (
	. "github.com/onsi/ginkgo"

	"fmt"

	"reflect"

	"github.com/davecgh/go-spew/spew"
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

func stringPointer(v string) *string                { return &v }
func intPointer(v int) *int                         { return &v }
func floatPointer(v float64) *float64               { return &v }
func boolPointer(v bool) *bool                      { return &v }
func interfacePointer() interface{}                 { var v interface{}; return &v }
func slicePointer(v []interface{}) *[]interface{}   { return &v }
func arrayPointer(v [2]interface{}) *[2]interface{} { return &v }

type entry struct {
	decoded interface{}
	encoded string
	pointer interface{}
	context string
	customDecode  func(entry)
}

var _ = FDescribe("Serializer", func() {
	//this is a slice of an anonymous struct (see https://talks.golang.org/2012/10things.slide#2)

	table := []entry{
		{
			context: "with a string",
			decoded: "a",
			encoded: `{"d":"a","s":[]}`,
			pointer: stringPointer("a"),
		},
		{
			context: "with an int",
			decoded: 1,
			encoded: `{"d":1,"s":[]}`,
			pointer: intPointer(1),
		},
		{
			context: "with a float",
			decoded: float64(1.1),
			encoded: `{"d":1.1,"s":[]}`,
			pointer: floatPointer(float64(1.1)),
		},
		{
			context: "with true",
			decoded: true,
			encoded: `{"d":true,"s":[]}`,
			pointer: boolPointer(true),
		},
		{
			context: "with false",
			decoded: false,
			encoded: `{"d":false,"s":[]}`,
			pointer: boolPointer(false),
		},
		{
			context: "with nil",
			decoded: nil,
			encoded: `{"d":null,"s":[]}`,
			pointer: interfacePointer(),
		},
		{
			context: "with a simple slice",
			decoded: []interface{}{1, "e"},
			encoded: `{"o":"@zenaton#0","s":[{"v":[1,"e"]}]}`,
			pointer: slicePointer([]interface{}{}),
		},
		{
			context: "with a simple array",
			decoded: [2]interface{}{1, "e"},
			encoded: `{"o":"@zenaton#0","s":[{"v":[1,"e"]}]}`,
			pointer: arrayPointer([2]interface{}{}),
			customDecode: decodeSimpleArray
		},
		//{
		//	context: "with circular arrays that only contain each other",
		//	decoded: [2]interface{}{1, "e"},
		//	encoded: `{"o":"@zenaton#0","s":[{"v":["@zenaton#1"]},{"v":["@zenaton#0"]}]}`,
		//	pointer: arrayPointer([2]interface{}{}),
		//},
	}

	s := zenaton.Serializer{}

	for _, entry := range table {
		entry := entry //gotcha!
		Context(entry.context, func() {

			Context("Encode", func() {
				It("should encode to data", func() {
					e, err := s.Encode(entry.decoded)
					Expect(err).ToNot(HaveOccurred())
					Expect(e).To(Equal(entry.encoded))
				})
			})

			Context("Decode", func() {
				It("should decode from data", func() {

					if entry.customDecode != nil {
						entry.customDecode(entry)
						return
					}

					err := s.Decode(entry.encoded, entry.pointer)
					Expect(err).ToNot(HaveOccurred())
					Expect(reflect.ValueOf(entry.pointer).Elem().Interface()).To(BeEquivalentTo(entry.decoded))
				})
			})
		})
	}
}

func decodeSimpleArray(e entry) {
	
}

var _ = Describe("Serializer", func() {
	s := zenaton.Serializer{}

	Context("with circular arrays that only contain each other", func() {
		It("represents the array as an array", func() {

			var arr1 [1]interface{}
			var arr2 [1]interface{}

			arr1 = [1]interface{}{&arr2}
			arr2 = [1]interface{}{&arr1}

			encoded, err := s.Encode(&arr1)
			Expect(err).ToNot(HaveOccurred())
			fmt.Println(encoded)
			Expect(encoded).To(Equal(`{"o":"@zenaton#0","s":[{"v":["@zenaton#1"]},{"v":["@zenaton#0"]}]}`))
		})
	})

	Context("with circular arrays that only contain each other", func() {
		Context("without receiving a pointer to the array", func() {
			It("represents the array as an array", func() {

				var arr1 [1]interface{}
				var arr2 [1]interface{}

				arr1 = [1]interface{}{&arr2}
				arr2 = [1]interface{}{&arr1}

				encoded, err := s.Encode(arr1)
				Expect(err).ToNot(HaveOccurred())
				fmt.Println(encoded)
				Expect(encoded).To(Equal(`{"o":"@zenaton#0","s":[{"v":["@zenaton#1"]},{"v":["@zenaton#2"]},{"v":["@zenaton#1"]}]}`))
			})
		})
	})

	Context("with an array inside an array", func() {
		It("represents the array as an array", func() {

			var arr1 [1]interface{}
			arr2 := [1]interface{}{1}

			arr1 = [1]interface{}{arr2}

			encoded, err := s.Encode(arr1)
			Expect(err).ToNot(HaveOccurred())
			fmt.Println(encoded)
			Expect(encoded).To(Equal(`{"o":"@zenaton#0","s":[{"v":["@zenaton#1"]},{"v":[1]}]}`))
		})
	})

	Context("with circular slices", func() {
		It("represents the slice as an array", func() {

			var arr1 []interface{}
			var arr2 []interface{}

			arr1 = []interface{}{&arr2}
			arr2 = []interface{}{&arr1}

			encoded, err := s.Encode(&arr1)
			Expect(err).ToNot(HaveOccurred())
			fmt.Println(encoded)
			Expect(encoded).To(Equal(`{"o":"@zenaton#0","s":[{"v":["@zenaton#1"]},{"v":["@zenaton#0"]}]}`))
		})
	})
	Context("with circular arrays that contain each other and other elements", func() {
		It("represents the array as an array", func() {

			var arr1 [3]interface{}
			var arr2 [3]interface{}

			arr1 = [3]interface{}{&arr2, 1, 2}
			arr2 = [3]interface{}{&arr1, 3, 4}

			encoded, err := s.Encode(&arr1)
			Expect(err).ToNot(HaveOccurred())
			fmt.Println(encoded)
			Expect(encoded).To(Equal(`{"o":"@zenaton#0","s":[{"v":["@zenaton#1",1,2]},{"v":["@zenaton#0",3,4]}]}`))
		})
	})

	//todo: this needs to be different, as maps can have pointer keys and values
	Context("with a simple map", func() {
		It("represents the map as an map", func() {
			m := map[string]string{
				"k1": "v1",
				"k2": "v2",
			}
			encoded, err := s.Encode(m)
			Expect(err).ToNot(HaveOccurred())
			fmt.Println(encoded)
			Expect(encoded).To(Equal(`{"o":"@zenaton#0","s":[{"k":["k1","k2"],"v":["v1","v2"]}]}`))
		})
	})

	Context("with a circular map", func() {
		It("represents the map as an map", func() {

			m1 := make(map[string]interface{})
			m2 := make(map[*map[string]interface{}]interface{})
			m1["m2"] = &m2
			//here we make sure that the keys can also be pointers
			m2[&m1] = "m1"

			encoded, err := s.Encode(&m1)
			Expect(err).ToNot(HaveOccurred())
			fmt.Println(encoded)
			Expect(encoded).To(Equal(`{"o":"@zenaton#0","s":[{"k":["m2"],"v":["@zenaton#1"]},{"k":["@zenaton#0"],"v":["m1"]}]}`))
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

			parent := Person{}
			child := Person{
				Parent: &parent,
			}

			parent.Child = &child

			expectedSerialized := `{"o":"@zenaton#0","s":[{"n":"Person","p":{"Child":"@zenaton#1","Parent":null}},{"n":"Person","p":{"Child":null,"Parent":"@zenaton#0"}}]}`
			encoded, err := s.Encode(&parent)
			Expect(err).ToNot(HaveOccurred())
			fmt.Println("encoded: ", encoded)
			Expect(encoded).To(Equal(expectedSerialized))
		})
	})

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

	Context("with a simple struct", func() {
		It("returns an instance with the correct instance variables", func() {

			pointed := "v"
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
				String:  "v",
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
										"String":"v",
										"Ptr":"v"
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

	Context("with a simple struct", func() {
		It("returns an instance with the correct instance variables", func() {

			pointed := "v"
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
				String:  "v",
				Ptr:     &pointed2,
			}

			encoded := `{
							   "o":"@zenaton#0",
							   "s":[
								  {
									 "n":"Person",
									 "p":{
										"Child":"@zenaton#1",
										"Parent":null
									 }
								  },
								  {
									 "n":"Person",
									 "p":{
										"Child":null,
										"Parent":"@zenaton#0"
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

	Context("with a simple slice", func() {
		It("decodes into an slice with the correct contents", func() {

			encoded := `{
							   "o":"@zenaton#0",
							   "s":[
								  {
									 "v":[
										"a",
										"b"
									 ]
								  }
							   ]
							}`
			expectedArr := []string{"a", "b"}

			var decodedArr []string
			err := s.Decode(encoded, &decodedArr)
			Expect(err).ToNot(HaveOccurred())
			Expect(decodedArr).To(Equal(expectedArr))
		})
	})

	Context("with a slice inside a slice", func() {
		It("decodes into an slice with the correct contents", func() {

			encoded := `{
							   "o":"@zenaton#0",
							   "s":[
								  {
									 "v":["@zenaton#1"]
								  },
								  {
									 "v":[1]
								  }
							   ]
							}`
			expectedArr := [][]int{{1}}

			var decodedArr [][]int
			err := s.Decode(encoded, &decodedArr)
			Expect(err).ToNot(HaveOccurred())
			Expect(decodedArr).To(Equal(expectedArr))
		})
	})

	Context("with a an array inside an array", func() {
		It("decodes into an array with the correct contents", func() {

			encoded := `{
							   "o":"@zenaton#0",
							   "s":[
								  {
									 "v":["@zenaton#1"]
								  },
								  {
									 "v":[1]
								  }
							   ]
							}`
			expectedArr := [1][1]int{{1}}

			var decodedArr [1][1]int
			err := s.Decode(encoded, &decodedArr)
			Expect(err).ToNot(HaveOccurred())
			Expect(decodedArr).To(Equal(expectedArr))
		})
	})

	Context("with a circular array", func() {
		It("decodes into an array with the correct contents", func() {

			encoded := `{
							   "o":"@zenaton#0",
							   "s":[
								  {
									 "v":["@zenaton#1", 1]
								  },
								  {
									 "v":["@zenaton#0", 2]
								  }
							   ]
							}`
			var arr1 [2]interface{}
			var arr2 [2]interface{}

			arr1 = [2]interface{}{&arr2, 1}
			arr2 = [2]interface{}{&arr1, 2}

			var decodedArr [2]interface{}
			err := s.Decode(encoded, &decodedArr)
			spew.Dump("arr1", &arr1)
			spew.Dump("decodedArr", &decodedArr)
			Expect(err).ToNot(HaveOccurred())
			// the decoder cannot determine the type of the second array, so it is decoded as a slice (as the json unmarshaller would do)
			secondArray := decodedArr[0].(*[]interface{})
			Expect((*secondArray)[0]).To(Equal(&decodedArr))
		})
	})

	Context("with a circular slice", func() {
		It("decodes into an slice with the correct contents", func() {

			encoded := `{
							   "o":"@zenaton#0",
							   "s":[
								  {
									 "v":["@zenaton#1"]
								  },
								  {
									 "v":["@zenaton#0"]
								  }
							   ]
							}`
			var arr1 []interface{}
			var arr2 []interface{}

			arr1 = []interface{}{&arr2}
			arr2 = []interface{}{&arr1}

			var decodedArr []interface{}
			err := s.Decode(encoded, &decodedArr)
			spew.Dump("decodedArr", &decodedArr)
			Expect(err).ToNot(HaveOccurred())
			secondArray := decodedArr[0].(*[]interface{})
			Expect((*secondArray)[0]).To(Equal(&decodedArr))
		})
	})

	Context("an boolean value", func() {
		It("returns the boolean", func() {
			encoded := `{"d":true,"s":[]}`
			var myBool bool
			err := s.Decode(encoded, &myBool)
			Expect(err).ToNot(HaveOccurred())
			Expect(myBool).To(Equal(true))
		})
	})

	Context("an uint value", func() {
		It("returns the uint", func() {
			encoded := `{"d":1,"s":[]}`
			var myUint uint
			err := s.Decode(encoded, &myUint)
			Expect(err).ToNot(HaveOccurred())
			Expect(myUint).To(Equal(uint(1)))
		})
	})
})

//todo: I should test that the types match in decode, right?
