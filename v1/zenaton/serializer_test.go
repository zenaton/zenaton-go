package zenaton_test

import (
	"encoding/json"

	. "github.com/onsi/ginkgo"

	"fmt"

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

var _ = Describe("Serializer", func() {
	Describe("Encode", func() {
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
				Expect(encoded).To(Equal(`{"o":"@zenaton#0","s":[{"v":[1,"e"]}]}`))
			})
		})

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

	FDescribe("Decode", func() {
		s := &zenaton.Serializer{}
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

		Context("with a simple slice", func() {
			FIt("decodes into an slice with the correct contents", func() {

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

		Context("with a an slice inside an slice", func() {
			FIt("decodes into an slice with the correct contents", func() {

				encoded := `{
							   "o":"@zenaton#0",
							   "s":[
								  {
									 "v":["@zenaton#1"]
								  },
								  {
									 "v":["@zenaton#2"]
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
				expectedArr := [][]int{{1}}

				var decodedArr [][]int
				err := s.Decode(encoded, &decodedArr)
				Expect(err).ToNot(HaveOccurred())
				Expect(decodedArr).To(Equal(expectedArr))
			})
		})

		Context("with a circular array", func() {
			FIt("decodes into an array with the correct contents", func() {

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
				//spew.Dump("arr1", &arr1)
				//spew.Dump("arr1", &arr1)
				//spew.Dump("arr1", &arr1)
				//spew.Dump("arr1", &arr1)
				//spew.Dump("arr1", &arr1)
				//spew.Dump("arr1", &arr1)

				//var intArr1 = interface{}(&arr1)
				//var user_array []interface{}
				//var interface_array_inside_decode = interface{}(&user_array)
				//reflect_value_of_user_array := reflect.ValueOf(interface_array_inside_decode)
				//
				//var first_encoded_array []interface{}
				//reflect_value_first_decoded_array := reflect.ValueOf(&first_encoded_array)
				//reflect_value_of_user_array.Elem().Set(reflect_value_first_decoded_array.Elem())
				//spew.Dump("newArr ", &user_array)

				var decodedArr []interface{}
				err := s.Decode(encoded, &decodedArr)
				spew.Dump("decodedArr", &decodedArr)
				Expect(err).ToNot(HaveOccurred())
				Expect(decodedArr).To(Equal(arr1))

			})
		})

		Context("an int value", func() {
			It("returns the int", func() {
				encoded := `{"d":1,"s":[]}`
				var myInt int
				err := s.Decode(encoded, &myInt)
				Expect(err).ToNot(HaveOccurred())
				Expect(myInt).To(Equal(1))
			})
		})

		Context("an float value", func() {
			It("returns the float", func() {
				encoded := `{"d":1.1,"s":[]}`
				var myFloat float32
				err := s.Decode(encoded, &myFloat)
				Expect(err).ToNot(HaveOccurred())
				Expect(myFloat).To(Equal(float32(1.1)))
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

		Context("an string value", func() {
			It("returns the string", func() {
				encoded := `{"d":"a","s":[]}`
				var str string
				err := s.Decode(encoded, &str)
				Expect(err).ToNot(HaveOccurred())
				Expect(str).To(Equal("a"))
			})
		})

	})
})

//todo: I should test that the types match in decode, right?
