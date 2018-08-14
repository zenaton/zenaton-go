package zenaton_test

import (
	"encoding/json"

	. "github.com/onsi/ginkgo"

	"fmt"

	"reflect"

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

	Context("with recursive arrays", func() {
		FIt("represents the array as an array", func() {

			//
			//type MyInt int
			//
			//var i = 1
			//
			//encodedInt, err := s.Encode(i)
			//Expect(err).ToNot(HaveOccurred())
			//fmt.Println(encodedInt)
			//
			//fmt.Println(reflect.TypeOf(i))
			//var bob MyInt
			//interfaceInt := intrfc(&bob)
			//err = s.Decode(encodedInt, interfaceInt)
			//Expect(err).ToNot(HaveOccurred())
			//Expect(*interfaceInt.(*MyInt)).To(Equal(MyInt(1)))
			//
			//var bob2 MyInt
			//interfaceInt = intrfc(bob2)
			//interfaceInt = reflect.New(reflect.TypeOf(interfaceInt)).Interface()
			//err = s.Decode(encodedInt, interfaceInt)
			//Expect(err).ToNot(HaveOccurred())
			//Expect(*interfaceInt.(*MyInt)).To(Equal(MyInt(1)))
			//
			//type ID struct {
			//	ID int
			//}
			//
			//var id = ID{1}
			//
			//encodedInt, err = s.Encode(id)
			//Expect(err).ToNot(HaveOccurred())
			//fmt.Println(encodedInt)
			//
			//fmt.Println(reflect.TypeOf(i))
			//var bobID ID
			//interfaceInt = intrfc(&bobID)
			//err = s.Decode(encodedInt, interfaceInt)
			//Expect(err).ToNot(HaveOccurred())
			//Expect(*interfaceInt.(*ID)).To(Equal(ID{1}))
			//
			//var bobID2 ID
			//interfaceInt = intrfc(bobID2)
			//interfaceInt2 := reflect.New(reflect.TypeOf(interfaceInt)).Interface()
			//err = s.Decode(encodedInt, interfaceInt2)
			//Expect(err).ToNot(HaveOccurred())
			//fmt.Println("interfaceInt2: ", interfaceInt2)
			//Expect(*interfaceInt2.(*ID)).To(Equal(ID{1}))
			//
			//type IDmax struct {
			//	ID  int
			//	Max int
			//}
			//
			//idmax := intrfc(IDmax{1, 2})
			//
			//encodedIDmax, err := s.Encode(idmax)
			//Expect(err).ToNot(HaveOccurred())
			//fmt.Println(encodedIDmax)
			//
			//fmt.Println(reflect.TypeOf(idmax))
			//bobmax := IDmax{}
			//newIdmax := intrfc(&bobmax)
			//err = s.Decode(encodedIDmax, newIdmax)
			//Expect(err).ToNot(HaveOccurred())
			//fmt.Printf("newIdmax: %+v\n", newIdmax)
			//
			//bobmax = IDmax{}
			//newIdmax = intrfc(bobmax)
			//newIdmax = reflect.New(reflect.TypeOf(newIdmax)).Interface()
			//err = s.Decode(encodedIDmax, newIdmax)
			//Expect(err).ToNot(HaveOccurred())
			//fmt.Printf("newIdmax: %+v\n", newIdmax)

			expectedOutput := `{"o":"@zenaton#0","s":[{"a":["@zenaton#1"]},{"a":["@zenaton#0"]}]}`
			var arr1 [1]interface{}
			var arr2 [1]interface{}

			arr1 = [1]interface{}{&arr2}
			arr2 = [1]interface{}{&arr1}

			//arr3 := []interface{}{&arr2}
			//
			//spew.Dump(arr1)
			//spew.Dump(arr3)

			intrfcArr1 := intrfc(&arr1)
			//intrfcArr2 := intrfc(arr2)

			fmt.Println("arr1: ", reflect.ValueOf(&arr1).Pointer())
			fmt.Println("arr1: ", reflect.ValueOf(intrfcArr1).Pointer())
			fmt.Println("arr1 in arr2: ", reflect.ValueOf(arr2[0]).Pointer())

			fmt.Println("test arr1: ", reflect.ValueOf(&arr1).Pointer(), "test arr2: ", reflect.ValueOf(&arr2).Pointer())
			encoded, err := s.Encode(&arr1)
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

			var indexOf = func(slice []interface{}, item interface{}) int {
				for i := range slice {
					fmt.Println(slice[i] == item)
				}
				return -1
			}

			parent := Person{}
			child := Person{
				Parent: &parent,
			}
			parent.Child = &child

			parent2 := Person{}
			child2 := Person{
				Parent: &parent2,
			}
			parent.Child = &child2

			parents := []interface{}{&parent, &parent2}

			fmt.Println(indexOf(parents, &parent))

			//expectedSerialized := `{"o":"@zenaton#0","s":[{"n":"Person","p":{"Child":"@zenaton#1","Parent":null}},{"n":"Person","p":{"Child":null,"Parent":"@zenaton#0"}}]}`
			//encoded, err := s.Encode(parent)
			//Expect(err).ToNot(HaveOccurred())
			//fmt.Println("encoded: ", encoded)
			//Expect(encoded).To(Equal(expectedSerialized))
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

func intrfc(i interface{}) interface{} {
	return i
}
