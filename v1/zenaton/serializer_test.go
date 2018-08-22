package zenaton

import (
	"fmt"

	"strings"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	//"github.com/zenaton/zenaton-go/v1/zenaton"
	"reflect"

	"github.com/davecgh/go-spew/spew"
)

var s = Serializer{}
var pointed = "v"
var pointed2 = &pointed

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

func stringPointer() *string         { var v string; return &v }
func intPointer() *int               { var v int; return &v }
func uintPointer() *uint             { var v uint; return &v }
func floatPointer() *float64         { var v float64; return &v }
func boolPointer() *bool             { var v bool; return &v }
func interfacePointer() interface{}  { var v interface{}; return &v }
func slicePointer() *[]string        { var v []string; return &v }
func arrayPointer() *[2]string       { var v [2]string; return &v }
func nestedArrayPointer() *[1][1]int { var v [1][1]int; return &v }
func nestedSlicePointer() *[][]int   { var v [][]int; return &v }
func mapPointer() *map[string]string { v := make(map[string]string); return &v }
func structPointer() *MySimpleStruct { var v MySimpleStruct; return &v }
func circularStructPointer() *Person { var v Person; return &v }

var testString string

type entry struct {
	decoded    interface{}
	encoded    string
	pointer    interface{}
	context    string
	testDecode bool
}

var _ = Describe("Serializer", func() {
	s := Serializer{}

	table := []entry{
		{
			context: "with a string",
			decoded: "a",
			encoded: `{"d":"a","s":[]}`,
			pointer: stringPointer(),
		},
		{
			context: "with pointer to a string",
			decoded: "a",
			encoded: `{"d":"a","s":[]}`,
			pointer: stringPointer(),
		},
		{
			context: "with an int",
			decoded: 1,
			encoded: `{"d":1,"s":[]}`,
			pointer: intPointer(),
		},
		{
			context: "with an uint",
			decoded: 1,
			encoded: `{"d":1,"s":[]}`,
			pointer: uintPointer(),
		},
		{
			context: "with a float",
			decoded: float64(1.1),
			encoded: `{"d":1.1,"s":[]}`,
			pointer: floatPointer(),
		},
		{
			context: "with true",
			decoded: true,
			encoded: `{"d":true,"s":[]}`,
			pointer: boolPointer(),
		},
		{
			context: "with false",
			decoded: false,
			encoded: `{"d":false,"s":[]}`,
			pointer: boolPointer(),
		},
		{
			context: "with a simple array",
			decoded: [2]string{"a", "b"},
			encoded: `{"o":"@zenaton#0","s":[{"v":["a","b"]}]}`,
			pointer: arrayPointer(),
		},
		{
			context: "with a simple slice",
			decoded: []string{"a", "b"},
			encoded: `{"o":"@zenaton#0","s":[{"v":["a","b"]}]}`,
			pointer: slicePointer(),
		},
		{
			context: "with an array inside an array",
			decoded: [1][1]int{{1}},
			encoded: `{"o":"@zenaton#0","s":[{"v":["@zenaton#1"]},{"v":[1]}]}`,
			pointer: nestedArrayPointer(),
		},
		{
			context: "with a slice inside a slice",
			decoded: [][]int{{1}},
			encoded: `{"o":"@zenaton#0","s":[{"v":["@zenaton#1"]},{"v":[1]}]}`,
			pointer: nestedSlicePointer(),
		},
		{
			context: "with a simple map",
			decoded: map[string]string{"k1": "v1"},
			encoded: `{"o":"@zenaton#0","s":[{"k":["k1"],"v":["v1"]}]}`,
			pointer: mapPointer(),
		},
		{
			context: "with a simple struct",
			decoded: MySimpleStruct{
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
			},
			//fields must be in alphabetical order
			encoded: `{
					   "o":"@zenaton#0",
					   "s":[
						  {
							 "n":"MySimpleStruct",
							 "p":{
								"Bool":true,
								"Float32":1.1,
								"Float64":1.1,
								"Int":1,
								"Int16":1,
								"Int32":1,
								"Int64":1,
								"Int8":1,
								"Ptr":"v",
								"String":"v",
								"Uint":1,
								"Uint16":1,
								"Uint32":1,
								"Uint64":1,
								"Uint8":1
							 }
						  }
					   ]
					}`,
			pointer: structPointer(),
		},
		//todo:
		//{
		//	context: "with nil",
		//	decoded: nil,
		//	encoded: `{"d":null,"s":[]}`,
		//	pointer: interfacePointer(),
		//},
	}

	for _, entry := range table {
		entry := entry //gotcha!
		Context(entry.context, func() {

			Context("Encode", func() {
				It("should encode to data", func() {
					e, err := s.Encode(entry.decoded)
					Expect(err).ToNot(HaveOccurred())
					Expect(e).To(Equal(strings.Join(strings.Fields(entry.encoded), ""))) //removes whitespace
				})
			})

			Context("Decode", func() {
				It("should decode from data", func() {
					err := s.Decode(entry.encoded, entry.pointer)
					Expect(err).ToNot(HaveOccurred())
					Expect(reflect.ValueOf(entry.pointer).Elem().Interface()).To(BeEquivalentTo(entry.decoded))
				})
			})
		})
	}

	Context("with circular arrays that only contain each other", func() {
		var decoded [1]interface{}
		var arr2 [1]interface{}

		decoded = [1]interface{}{&arr2}
		arr2 = [1]interface{}{&decoded}

		encoded := `{"o":"@zenaton#0","s":[{"v":["@zenaton#1"]},{"v":["@zenaton#0"]}]}`

		It("should encode to data", func() {
			encoded, err := s.Encode(&decoded)
			Expect(err).ToNot(HaveOccurred())
			fmt.Println(encoded)
			Expect(encoded).To(Equal(encoded))
		})

		Context("with circular arrays that only contain each other", func() {
			It("should encode to data", func() {
				encoded, err := s.Encode(decoded)
				Expect(err).ToNot(HaveOccurred())
				fmt.Println(encoded)
				Expect(encoded).To(Equal(`{"o":"@zenaton#0","s":[{"v":["@zenaton#1"]},{"v":["@zenaton#2"]},{"v":["@zenaton#1"]}]}`))
			})
		})

		It("should decode from data", func() {

			var arr [1]interface{}
			err := s.Decode(encoded, &arr)
			spew.Dump("**********************1: ", &decoded)
			spew.Dump("**********************2: ", &arr)
			Expect(err).ToNot(HaveOccurred())
			secondArray := arr[0].(*[]interface{})
			Expect((*secondArray)[0]).To(Equal(&arr))
		})
	})

	Context("with circular slices", func() {
		var decoded []interface{}
		var arr2 []interface{}

		decoded = []interface{}{&arr2, 1}
		arr2 = []interface{}{&decoded, 2}
		encoded := `{"o":"@zenaton#0","s":[{"v":["@zenaton#1"]},{"v":["@zenaton#0"]}]}`

		It("should encode to data", func() {

			encoded, err := s.Encode(&decoded)
			Expect(err).ToNot(HaveOccurred())
			fmt.Println(encoded)
			Expect(encoded).To(Equal(encoded))
		})

		It("should decode from data", func() {

			var decodedSlice []interface{}
			err := s.Decode(encoded, &decodedSlice)
			spew.Dump("decoded", &decoded)
			spew.Dump("decodedSlice", &decodedSlice)
			Expect(err).ToNot(HaveOccurred())
			secondArray := decodedSlice[0].(*[]interface{})
			Expect((*secondArray)[0]).To(Equal(&decodedSlice))
		})
	})

	Context("with a circular struct", func() {

		var Parent Person
		var Child = Person{Parent: &Parent}
		Parent.Child = &Child

		decoded := &Parent
		encoded := `{"o":"@zenaton#0","s":[{"n":"Person","p":{"Child":"@zenaton#1","Parent":null}},{"n":"Person","p":{"Child":null,"Parent":"@zenaton#0"}}]}`

		It("should encode to data", func() {

			encoded, err := s.Encode(decoded)
			Expect(err).ToNot(HaveOccurred())
			fmt.Println(encoded)
			Expect(encoded).To(Equal(encoded))
		})

		It("should decode from data", func() {

			var toDecode Person
			err := s.Decode(encoded, &toDecode)
			Expect(err).ToNot(HaveOccurred())
			Expect(toDecode.Child.Parent).To(Equal(&toDecode))
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

	Context("with a circular map", func() {
		m1 := make(map[string]interface{})
		m2 := make(map[*map[string]interface{}]interface{})
		m1["m2"] = &m2
		//here we make sure that the keys can also be pointers
		m2[&m1] = "m1"

		decoded := &m1
		encoded := `{"o":"@zenaton#0","s":[{"k":["m2"],"v":["@zenaton#1"]},{"k":["@zenaton#0"],"v":["m1"]}]}`

		It("should encode to data", func() {

			encoded, err := s.Encode(decoded)
			Expect(err).ToNot(HaveOccurred())
			fmt.Println(encoded)
			Expect(encoded).To(Equal(encoded))
		})

		It("should decode from data", func() {

			var toDecode map[string]interface{}
			err := s.Decode(encoded, &toDecode)
			Expect(err).ToNot(HaveOccurred())
			m2 := toDecode["m2"].(map[interface{}]interface{})
			var keys []interface{}
			for k := range m2 {
				keys = append(keys, k)
			}
			Expect(keys[0]).To(Equal(&toDecode))
		})
	})
})

//todo: I should test that the types match in decode, right?
