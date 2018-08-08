package zenaton

import (
	"encoding/json"
	"fmt"
	"reflect"
	"strconv"
	"strings"
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
	decoded []interface{}
}

type Object struct {
	Name       string                 `json:"n"`
	Properties map[string]interface{} `json:"p"`
}

//todo: is this enough? what about the recursive case?
//todo: make sure anonomous structs work
//todo: make sure maps with other things as keys work

func (s *Serializer) Encode(data interface{}) (string, error) {

	s.encoded = []interface{}{}
	s.decoded = []interface{}{}

	value := make(map[string]interface{})

	if reflect.TypeOf(data) == nil {
		value[KEY_DATA] = data
	} else {
		//todo: handle pointers to values
		switch reflect.Indirect(reflect.ValueOf(data)).Kind() {
		case reflect.Struct:
			value[KEY_OBJECT] = s.encodeObject(data)
		case reflect.Array:
			value[KEY_ARRAY] = s.encodeArray(data.([]interface{}))
		case reflect.Slice:
			value[KEY_ARRAY] = s.encodeArray(data.([]interface{}))
		case reflect.Map:
			value[KEY_ARRAY] = data
		default:
			value[KEY_DATA] = data
		}
	}

	value[KEY_STORE] = s.encoded
	//spew.Dump("value: %+v", value)
	encoded, err := json.Marshal(value)
	if err != nil {
		return "", err
	}

	return string(encoded), nil
}

func (s *Serializer) encodeArray(a []interface{}) []interface{} {

	var array []interface{}
	for _, v := range a {
		vValue := reflect.Indirect(reflect.ValueOf(v))
		switch vValue.Kind() {
		case reflect.Struct:
			array = append(array, s.encodeObject(v))
		case reflect.Slice:
			//todo: I hope this is ok
			array = append(array, s.encodeArray(vValue.Interface().([]interface{})))
		case reflect.Array:
			//todo: I hope this is ok
			array = append(array, s.encodeArray(v.([]interface{})))
		default:
			array = append(array, v)
		}
	}
	return array
}

func (s *Serializer) encodeObject(o interface{}) string {

	id := indexOf(o, s.decoded)

	fmt.Println(id)
	if id == -1 {
		dataT := reflect.TypeOf(o)
		id = len(s.decoded)

		s.decoded = insert(s.decoded, o, id)
		newEncoded := Object{
			Name:       dataT.Name(),
			Properties: s.encodeProperties(o),
		}

		s.encoded = insert(s.encoded, newEncoded, id)
	}

	return ID_PREFIX + strconv.Itoa(id)
}

func (s *Serializer) encodeProperties(o interface{}) map[string]interface{} {
	dataV := reflect.Indirect(reflect.ValueOf(o))
	dataT := reflect.TypeOf(o)
	propMap := make(map[string]interface{})
	for i := 0; i < dataV.NumField(); i++ {
		key := dataT.Field(i).Name
		//todo: handle the other cases
		//fmt.Println("kind: ", reflect.Indirect(dataV.Field(i)).Kind())
		switch reflect.Indirect(dataV.Field(i)).Kind() {
		case reflect.Struct:
			value := reflect.Indirect(dataV.Field(i)).Interface()
			propMap[key] = s.encodeObject(value)
		default:
			propMap[key] = dataV.Field(i).Interface()
		}
	}
	return propMap
}

func indexOf(item interface{}, slice []interface{}) int {
	for i, v := range slice {
		if reflect.DeepEqual(v, item) {
			return i
		}
	}
	return -1
}

func (s *Serializer) Decode(data string, value interface{}) error {
	////todo: handle race conditions
	//
	var v map[string]interface{}
	err := json.Unmarshal([]byte(data), &v)

	if err != nil {
		return err
	}

	s.decoded = []interface{}{}
	s.encoded = v[KEY_STORE].([]interface{})

	id, ok := v[KEY_OBJECT]
	if ok {
		idInt, err := strconv.Atoi(strings.TrimLeft(id.(string), ID_PREFIX))
		if err != nil {
			return err
		}
		s.decodeObject(idInt, s.encoded[idInt], value)
		return nil
	}

	fmt.Println("v, err", v, err)
	return nil
}

func (s *Serializer) decodeObject(id int, encodedObject interface{}, decodeTo interface{}) {

	fmt.Printf("id: %+v\n encodedObject: %+v\n decodeTo: %+v\n", id, encodedObject, decodeTo)
	if len(s.decoded) > id {
		decodeTo = s.decoded[id]
		return
	}

	//todo: this only works with ptr value, make sure that is always going to be the case
	o := encodedObject.(map[string]interface{})
	s.decoded = append(s.decoded, o)

	dataV := reflect.ValueOf(decodeTo)
	for key, value := range o[KEY_OBJECT_PROPERTIES].(map[string]interface{}) {

		fmt.Println("reflect.Indirect(reflect.Indirect(dataV).FieldByName(key)).Kind()", reflect.Indirect(reflect.Indirect(dataV).FieldByName(key)).Kind(), value)
		if reflect.Indirect(reflect.Indirect(dataV).FieldByName(key)).Kind() == reflect.Invalid {
			fmt.Println("my pointerzzzzzzz: ", dataV.Elem().FieldByName(key).Kind())
		}
		dataVField := reflect.Indirect(dataV.Elem().FieldByName(key))
		switch dataVField.Kind() {
		case reflect.Bool:
			dataVField.SetBool(value.(bool))
		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
			dataVField.SetInt(int64(value.(float64)))
		case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
			dataVField.SetUint(uint64(value.(float64)))
		case reflect.Float32, reflect.Float64:
			dataVField.SetFloat(value.(float64))
		//todo: this is not supported by json i guess?
		//case reflect.Complex64, reflect.Complex128:
		//	var c complex128
		//	str := fmt.Sprintf(`"%s"`, value.(string))
		//	err := json.Unmarshal([]byte(str), &c)
		//	if err != nil {
		//		//todo: panic?
		//		panic(err)
		//	}
		//	dataVField.SetComplex(c)
		case reflect.String:
			dataVField.SetString(value.(string))
			//todo:
		//case reflect.Uintptr:
		//case reflect.Array:
		//case reflect.Chan:
		//case reflect.Func:
		//case reflect.Interface:
		//case reflect.Map:
		//case reflect.Ptr:
		//case reflect.Slice:
		case reflect.Struct:
			dataVField.Set(reflect.ValueOf(value))
			//case reflect.UnsafePointer:
			//default:
		}
	}
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
