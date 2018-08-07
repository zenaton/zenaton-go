package zenaton

import (
	"encoding/json"
	"fmt"
	"reflect"
	"strconv"
)

const (
	ID_PREFIX             = "@zenaton#"
	KEY_OBJECT            = "o" // JSON key for objects
	KEY_OBJECT_NAME       = "n" // JSON key for class name
	KEY_OBJECT_PROPERTIES = "p" // JSON key for object ivars
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
	fmt.Println("id: ", id)

	if id == -1 {
		dataT := reflect.TypeOf(o)
		id = len(s.decoded)

		s.decoded = append(s.decoded, o)
		newEncoded := Object{
			Name:       dataT.Name(),
			Properties: s.encodeProperties(o),
		}
		s.encoded = append(s.encoded, newEncoded)
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
		fmt.Println("kind: ", reflect.Indirect(dataV.Field(i)).Kind())
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
	//var v map[string]interface{}
	//err := json.Unmarshal([]byte(data), &v)
	//
	//if err != nil {
	//	return err
	//}
	//
	//s.decoded = []interface{}{}
	//s.encoded = v[KEY_STORE].([]interface{})
	//
	//id, ok := v[KEY_OBJECT]
	//if ok {
	//	idInt, err := strconv.Atoi(strings.TrimLeft(id.(string), ID_PREFIX))
	//	if err != nil {
	//		return err
	//	}
	//	s.decodeObject(idInt, s.encoded[idInt], value)
	//	return nil
	//}
	//
	//fmt.Println("v, err", v, err)
	return nil
}

func (s *Serializer) DecodeObject() {
	//dataV := reflect.ValueOf(decodeTo)
	//reflect.Struct
	//obj := encodedObject.(Object)
	//for key, value := range obj.Properties {
	//	switch dataV.Kind() {
	//	case reflect.Bool:
	//		fmt.Println("hi")
	//		dataV.Elem().FieldByName(key).SetBool(value.(bool))
	//	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
	//		dataV.Elem().FieldByName(key).SetInt(value.(int64))
	//	case reflect.Uint:
	//	case reflect.Uint8:
	//	case reflect.Uint16:
	//	case reflect.Uint32:
	//	case reflect.Uint64:
	//	case reflect.Uintptr:
	//	case reflect.Float32:
	//	case reflect.Float64:
	//		//case reflect.Complex64:
	//		//case reflect.Complex128:
	//		//case reflect.Array:
	//		//case reflect.Chan:
	//		//case reflect.Func:
	//		//case reflect.Interface:
	//		//case reflect.Map:
	//		//case reflect.Ptr:
	//		//case reflect.Slice:
	//		//case reflect.String:
	//		//case reflect.Struct:
	//		//case reflect.UnsafePointer:
	//		//default:
	//	}
	//
	//}
	//
	////	value := dataV.Field(i).Interface()
	////
	////	obj := Object{
	////		N: dataT.Name(),
	////		P: map[string]interface{}{
	////			dataT.Field(i).Name: value,
	////		},
	////	}
	////	objs = append(objs, obj)
	////todo: is this right?
	//fmt.Printf("id: %+v\n encodedObject: %+v\n decodeTo: %+v\n", id, encodedObject, decodeTo)
	//
	//if len(s.decoded) > id {
	//	return s.decoded[id]
	//}
	//return nil
}
