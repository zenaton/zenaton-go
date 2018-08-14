package zenaton

import (
	"encoding/json"
	"fmt"
	"reflect"
	"strconv"

	"github.com/davecgh/go-spew/spew"
)

const (
	ID_PREFIX             = "@zenaton#"
	KEY_OBJECT            = "o" // JSON key for objects
	KEY_OBJECT_NAME       = "n" // JSON key for class name
	KEY_OBJECT_PROPERTIES = "p" // JSON key for object vars
	KEY_ARRAY             = "a" // JSON key for array and hashes
	KEY_DATA              = "d" // JSON key for json compatibles types
	KEY_STORE             = "s" // JSON key for deserialized complex object
)

type Serializer struct {
	encoded []interface{}
	decoded []uintptr
}

type Object struct {
	Name       string                 `json:"n"`
	Properties map[string]interface{} `json:"p"`
}

//todo: is this enough? what about the recursive case?
//todo: make sure anonomous structs work
//todo: handle pointers to values
//todo: make sure maps with other things besides strings as keys work
//todo: handle interfaces

func (s *Serializer) Encode(data interface{}) (string, error) {
	//rv := indirect(reflect.ValueOf(data))
	//arr2Pointer := reflect.ValueOf(rv.Interface().([]interface{})[0]).Pointer()
	//
	//fmt.Println("kind: ", reflect.ValueOf(rv.Interface().([]interface{})[0]).Kind(), "arr1: ", reflect.ValueOf(data).Pointer(), "arr2: ", arr2Pointer)
	//
	//rv = reflect.ValueOf(data)
	//kind := rv.Kind()
	//isValid := validType(kind)
	//if !isValid {
	//	return "", errors.New(fmt.Sprintf("cannot encode data of kind: %s", kind.String()))
	//}
	//
	//s.encoded = []interface{}{}
	//s.decoded = []uintptr{}
	//
	//value := make(map[string]interface{})
	//
	////todo: handle pointers to values
	////todo: handle interfaces
	//if basicType(kind) || reflect.TypeOf(data) == nil {
	//	value[KEY_DATA] = data
	//} else {
	//	value[KEY_OBJECT] = s.encodeToStore(rv)
	//}
	//
	//value[KEY_STORE] = s.encoded
	//encoded, err := json.Marshal(value)
	//if err != nil {
	//	return "", err
	//}
	//
	//return string(encoded), nil
	v, err := json.Marshal(data)
	return string(v), err
}

func basicType(kind reflect.Kind) bool {
	return kind == reflect.Bool ||
		kind == reflect.Int ||
		kind == reflect.Int8 ||
		kind == reflect.Int16 ||
		kind == reflect.Int32 ||
		kind == reflect.Int64 ||
		kind == reflect.Uint ||
		kind == reflect.Uint8 ||
		kind == reflect.Uint16 ||
		kind == reflect.Uint32 ||
		kind == reflect.Uint64 ||
		kind == reflect.Uintptr ||
		kind == reflect.Float32 ||
		kind == reflect.Float64 ||
		kind == reflect.String
}

func (s *Serializer) encodeToStore(object reflect.Value) string {
	spew.Dump(object)
	fmt.Println("kind in encodeToStore: ", object.Kind())
	if object.Kind() == reflect.Ptr {
		id := indexOf(s.decoded, object.Pointer())
		fmt.Println("id: ", id)
		if id != -1 {
			return storeID(id)
		}
	}
	return s.storeAndEncode(object)
}

func (s *Serializer) storeAndEncode(object reflect.Value) string {
	//fmt.Printf("2: %+v\n", object)
	id := len(s.decoded)
	if object.Kind() == reflect.Ptr {
		s.decoded = insertPtr(s.decoded, object.Pointer(), id)
	}
	fmt.Println("decode: ", s.decoded)
	if len(s.decoded) > 2 {
		panic("nope")
	}
	s.encoded = insert(s.encoded, s.encodedObjectByType(object), id)
	return storeID(id)
}

func (s *Serializer) encodedObjectByType(object reflect.Value) map[string]interface{} {
	object = reflect.Indirect(object)
	kind := object.Kind()
	fmt.Println("kind: ", kind)
	switch kind {
	case reflect.Struct:
		return s.encodeObject(object)
	case reflect.Slice:
		//fmt.Printf("4: %+v\n", object)
		return s.encodeArray(object)
		//case reflect.Array:
		//case reflect.Map:
	}

	//todo??
	return nil
}

func storeID(id int) string {
	return ID_PREFIX + strconv.Itoa(id)
}

func validType(kind reflect.Kind) bool {
	//todo: can we really not serialize complex?
	return !(kind == reflect.Complex64 ||
		kind == reflect.Complex128 ||
		kind == reflect.Chan ||
		kind == reflect.Func ||
		kind == reflect.UnsafePointer)
}

func (s *Serializer) encodeArray(a reflect.Value) map[string]interface{} {

	arrayToEncode := a.Interface().([]interface{})
	var array []interface{}
	for _, v := range arrayToEncode {
		rv := reflect.ValueOf(v)
		kind := rv.Kind()
		fmt.Println("encode Array element kind: ", kind, rv.Pointer())
		switch kind {
		case reflect.Struct, reflect.Array, reflect.Slice, reflect.Map:
			if reflect.ValueOf(v).Kind() == reflect.Ptr {
				fmt.Printf("))))))))))))))))) %+v\n", rv.Pointer())
			}
			array = append(array, s.encodeToStore(rv))
		default:
			array = append(array, v)
		}
	}

	return map[string]interface{}{
		KEY_ARRAY: array,
	}
}

func (s *Serializer) encodeObject(object interface{}) map[string]interface{} {
	return map[string]interface{}{
		KEY_OBJECT_NAME:       reflect.TypeOf(object).Name(),
		KEY_OBJECT_PROPERTIES: s.encodeProperties(object),
	}
}

func (s *Serializer) encodeProperties(o interface{}) map[string]interface{} {
	dataV := indirect(reflect.ValueOf(o))
	dataT := reflect.TypeOf(o)
	propMap := make(map[string]interface{})
	for i := 0; i < dataV.NumField(); i++ {
		key := dataT.Field(i).Name
		switch indirect(dataV.Field(i)).Kind() {
		case reflect.Struct, reflect.Array, reflect.Slice, reflect.Map:
			value := dataV.Field(i)
			propMap[key] = s.encodeToStore(value)
		default:
			propMap[key] = dataV.Field(i).Interface()
		}
	}
	return propMap
}

var count = 0

func indexOf(slice []uintptr, item uintptr) int {
	count++
	if count > 4 {
		panic(count)
	}

	fmt.Println("decoded: ", slice, "item: ", item)
	for i := range slice {
		if slice[i] == item {
			return i
		}
	}
	return -1
}

func (s *Serializer) Decode(data string, value interface{}) error {
	////todo: handle race conditions
	//var v map[string]interface{}
	//err := json.Unmarshal([]byte(data), &v)
	//
	//if err != nil {
	//	return err
	//}
	//
	//s.decoded = []uintptr{}
	//s.encoded = v[KEY_STORE].([]interface{})
	//
	//id, ok := v[KEY_OBJECT]
	//if ok {
	//	idInt, err := strconv.Atoi(strings.TrimLeft(id.(string), ID_PREFIX))
	//	if err != nil {
	//		return err
	//	}
	//	s.decodeStruct(idInt, s.encoded[idInt], value)
	//	return nil
	//}
	//
	//fmt.Println("v, err", v, err)
	//return nil
	return json.Unmarshal([]byte(data), value)
}

//todo: should I handle the case in which the KEY_OBJECT_PROPERTIES don't match the struct passed in? seems yes
//func (s *Serializer) decodeStruct(id int, encodedObject interface{}, decodeTo interface{}) {
//
//	fmt.Printf("id: %+v\n encodedObject: %+v\n decodeTo: %+v\n", id, encodedObject, decodeTo)
//	if len(s.decoded) > id {
//		decodeTo = s.decoded[id]
//		return
//	}
//
//	//todo: this only works with ptr value, make sure that is always going to be the case
//	o := encodedObject.(map[string]interface{})
//	s.decoded = append(s.decoded, o)
//
//	v := reflect.ValueOf(decodeTo)
//	for key, value := range o[KEY_OBJECT_PROPERTIES].(map[string]interface{}) {
//
//		field := v.Elem().FieldByName(key)
//
//		field = indirect(field)
//
//		switch field.Kind() {
//		case reflect.Bool:
//			field.SetBool(value.(bool))
//		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
//			field.SetInt(int64(value.(float64)))
//		case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
//			field.SetUint(uint64(value.(float64)))
//		case reflect.Float32, reflect.Float64:
//			field.SetFloat(value.(float64))
//		case reflect.String:
//			field.SetString(value.(string))
//			//todo:
//		//case reflect.Uintptr:
//		//case reflect.Array:
//		//case reflect.Chan:
//		//case reflect.Func:
//		//case reflect.Interface:
//		//case reflect.Map:
//		//case reflect.Slice:
//		case reflect.Struct:
//			field.Set(reflect.ValueOf(value))
//		//todo: this is not supported by json i guess?
//		//case reflect.Complex64, reflect.Complex128:
//		//	var c complex128
//		//	str := fmt.Sprintf(`"%s"`, value.(string))
//		//	err := json.Unmarshal([]byte(str), &c)
//		//	if err != nil {
//		//		//todo: panic?
//		//		panic(err)
//		//	}
//		//	fld.SetComplex(c)
//		//case reflect.UnsafePointer:
//		default:
//			panic(fmt.Sprintf("unknown kind: %s", field.Kind()))
//		}
//	}
//}

func setPtr(field reflect.Value, value interface{}) {
	switch field.Type().Elem().Kind() {
	case reflect.Ptr:
		//todo: recurse
	case reflect.Bool:
		v := value.(bool)
		field.Set(reflect.ValueOf(&v))
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		v := int64(value.(float64))
		field.Set(reflect.ValueOf(&v))
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		v := uint64(value.(float64))
		field.Set(reflect.ValueOf(&v))
	case reflect.Float32, reflect.Float64:
		v := value.(float64)
		field.Set(reflect.ValueOf(&v))
	case reflect.String:
		v := value.(string)
		field.Set(reflect.ValueOf(&v))
	}
	//todo: other possible kinds
}

func insert(arr []interface{}, value interface{}, i int) []interface{} {
	if len(arr) > i {
		arr[i] = value
		return arr
	}
	newArr := make([]interface{}, i+1)
	copy(newArr, arr)
	newArr[i] = value
	return newArr
}

func insertPtr(arr []uintptr, value uintptr, i int) []uintptr {
	if len(arr) > i {
		arr[i] = value
		return arr
	}
	newArr := make([]uintptr, i+1)
	copy(newArr, arr)
	newArr[i] = value
	return newArr
}

func indirect(v reflect.Value) reflect.Value {

	// makes indirect more safe to call on values that are not Ptr or Interface
	if v.Kind() != reflect.Ptr && v.Kind() != reflect.Interface {
		return v
	}

	// If v is a named type and is addressable,
	// start with its address, so that if the type has pointer methods,
	// we find them.
	if v.Kind() != reflect.Ptr && v.Type().Name() != "" && v.CanAddr() {
		v = v.Addr()
	}
	for {
		// Load value from interface, but only if the result will be
		// usefully addressable.
		if v.Kind() == reflect.Interface && !v.IsNil() {
			fmt.Println("hi")
			e := v.Elem()
			//if e.Kind() == reflect.Ptr && !e.IsNil() {
			fmt.Println("hi2")
			v = e
			continue
			//}
		}

		if v.Kind() != reflect.Ptr {
			break
		}

		if v.IsNil() {
			v.Set(reflect.New(v.Type().Elem()))
		}

		v = v.Elem()
	}
	return v
}

//todo: could be useful
//// unpackValue returns values inside of non-nil interfaces when possible.
//// This is useful for data types like structs, arrays, slices, and maps which
//// can contain varying types packed inside an interface.
//func (d *dumpState) unpackValue(v reflect.Value) reflect.Value {
//	if v.Kind() == reflect.Interface && !v.IsNil() {
//		v = v.Elem()
//	}
//	return v
//}
