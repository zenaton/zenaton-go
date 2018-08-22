package zenaton

import (
	"encoding/json"
	"errors"
	"fmt"
	"reflect"
	"strconv"

	"strings"

	"github.com/davecgh/go-spew/spew"
)

type Person struct {
	Parent *Person
	Child  *Person
}

const (
	ID_PREFIX             = "@zenaton#"
	KEY_OBJECT            = "o" // JSON key for objects
	KEY_OBJECT_NAME       = "n" // JSON key for class name
	KEY_OBJECT_PROPERTIES = "p" // JSON key for object vars
	KEY_ARRAY             = "v" // JSON key for array and hashes
	KEY_KEYS              = "k" // JSON key the array of keys (in a map)
	KEY_DATA              = "d" // JSON key for json compatibles types
	KEY_STORE             = "s" // JSON key for deserialized complex object
)

type Serializer struct {
	encoded  []interface{}
	decoded  []reflect.Value
	pointers []uintptr
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
	s.encoded = []interface{}{}
	s.pointers = []uintptr{}

	rv := reflect.ValueOf(data)
	kind := rv.Kind()
	isValid := validType(kind)
	if !isValid {
		return "", errors.New(fmt.Sprintf("cannot encode data of kind: %s", kind.String()))
	}

	value := make(map[string]interface{})

	//todo: handle pointers to values
	//todo: handle interfaces
	if basicType(rv) || reflect.TypeOf(data) == nil {
		value[KEY_DATA] = data
	} else {
		value[KEY_OBJECT] = s.encodeToStore(rv)
	}

	value[KEY_STORE] = s.encoded
	encoded, err := json.Marshal(value)
	if err != nil {
		return "", err
	}

	return string(encoded), nil

}

func basicType(rv reflect.Value) bool {
	v := rv
	for v.Kind() == reflect.Ptr || v.Kind() == reflect.Interface {
		v = v.Elem()
	}
	kind := v.Kind()
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

func (s *Serializer) encodeToStore(object reflect.Value) interface{} {
	spew.Dump(object)
	for object.Kind() == reflect.Interface {
		object = object.Elem()
	}
	fmt.Println("kind in encodeToStore: ", object.Kind())
	if object.Kind() == reflect.Ptr {
		if object.IsNil() {
			return nil
		}
		id := indexOf(s.pointers, object.Pointer())
		fmt.Println("id: ", id)
		if id != -1 {
			return storeID(id)
		}
	}
	return s.storeAndEncode(object)
}

func (s *Serializer) storeAndEncode(object reflect.Value) string {
	//fmt.Printf("2: %+v\n", object)
	id := len(s.pointers)
	if object.Kind() == reflect.Ptr {
		s.pointers = insertPtr(s.pointers, object.Pointer(), id)
	} else {
		// this pointer is never actually used. It is only added so that the length of pointers is correct
		s.pointers = insertPtr(s.pointers, reflect.ValueOf(&object).Pointer(), id)
	}
	if len(s.pointers) > 6 {
		//todo: do something better here
		panic("nope")
	}
	s.encoded = insert(s.encoded, s.encodedObjectByType(object), id)
	return storeID(id)
}

func (s *Serializer) encodedObjectByType(object reflect.Value) map[string]interface{} {
	object = reflect.Indirect(object)
	kind := object.Kind()
	fmt.Println("kind in encodedObjectByType: ", kind)
	switch kind {
	case reflect.Struct:
		return s.encodeStruct(object)
	case reflect.Array, reflect.Slice:
		return s.encodeArray(object)
	case reflect.Map:
		return s.encodeMap(object)
	}

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

	var array []interface{}
	for i := 0; i < a.Len(); i++ {
		rv := a.Index(i)
		fmt.Println("rv: ", rv)
		for rv.Kind() == reflect.Interface && rv.Elem().Kind() != reflect.Invalid {
			rv = rv.Elem()
		}
		kind := rv.Kind()
		fmt.Println("kind in encodeArray2::::: ", kind)
		if basicType(rv) || kind == reflect.Interface {
			fmt.Println("basic kind or interface: ", rv.Interface(), rv)
			array = append(array, rv.Interface())
			continue
		}
		array = append(array, s.encodeToStore(rv))
	}

	return map[string]interface{}{
		KEY_ARRAY: array,
	}
}

//todo: test with NaN keys :\
func (s *Serializer) encodeMap(m reflect.Value) map[string]interface{} {

	var keys []interface{}
	var values []interface{}

	keyValues := m.MapKeys()
	for _, kv := range keyValues {
		for kv.Kind() == reflect.Interface {
			kv = kv.Elem()
		}
		kind := kv.Kind()
		fmt.Println("key::::: ", kind)
		if basicType(kv) {
			keys = append(keys, kv.Interface())
		} else {
			keys = append(keys, s.encodeToStore(kv))
		}

		//todo: abstract this out into another function, as I keep doing this
		valueValue := m.MapIndex(kv)
		for valueValue.Kind() == reflect.Interface {
			valueValue = valueValue.Elem()
		}
		fmt.Println("value::::: ", valueValue.Kind())
		if basicType(valueValue) {
			values = append(values, valueValue.Interface())
		} else {
			values = append(values, s.encodeToStore(valueValue))
		}
	}

	return map[string]interface{}{
		KEY_KEYS:  keys,
		KEY_ARRAY: values,
	}
}

func (s *Serializer) encodeStruct(object reflect.Value) map[string]interface{} {

	return map[string]interface{}{
		KEY_OBJECT_NAME:       object.Type().Name(),
		KEY_OBJECT_PROPERTIES: s.encodeProperties(object),
	}
}

func (s *Serializer) encodeProperties(o reflect.Value) map[string]interface{} {
	dataT := o.Type()
	propMap := make(map[string]interface{})
	for i := 0; i < o.NumField(); i++ {
		key := dataT.Field(i).Name
		fmt.Println("1key: ", key, "kind: ", o.Field(i).Kind())
		if basicType(o.Field(i)) {
			propMap[key] = o.Field(i).Interface()
			continue
		}
		fmt.Println("")
		propMap[key] = s.encodeToStore(o.Field(i))
	}
	return propMap
}

func indexOf(slice []uintptr, item uintptr) int {

	fmt.Println("pointers: ", slice, "item: ", item)
	for i := range slice {
		if slice[i] == item {
			return i
		}
	}
	return -1
}

func (s *Serializer) Decode(data string, value interface{}) error {
	//todo: handle race conditions
	//todo: check that value is a pointer! think the json part at the biginning does it

	//todo: I think I can build a better object to unmarshal into
	var parsedJSON map[string]interface{}
	err := json.Unmarshal([]byte(data), &parsedJSON)
	if err != nil {
		return err
	}

	s.decoded = []reflect.Value{}
	s.encoded = parsedJSON[KEY_STORE].([]interface{})
	fmt.Println("encoded: ", s.encoded)

	rv := reflect.ValueOf(value)

	simpleValue, ok := parsedJSON[KEY_DATA]
	if ok {
		//json.Unmarshal()
		switch rv.Elem().Kind() {
		case reflect.Bool:
			rv.Elem().SetBool(simpleValue.(bool))
		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
			rv.Elem().SetInt(int64(simpleValue.(float64)))
		case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
			rv.Elem().SetUint(uint64(simpleValue.(float64)))
		case reflect.Float32, reflect.Float64:
			rv.Elem().SetFloat(simpleValue.(float64))
		case reflect.String:
			rv.Elem().SetString(simpleValue.(string))
		default:
			panic(fmt.Sprintf("unknown kind: %s", rv.Elem().Kind()))
		}
		return nil
	}

	id, ok := parsedJSON[KEY_OBJECT]
	if ok {
		idInt, err := strconv.Atoi(strings.TrimLeft(id.(string), ID_PREFIX))
		if err != nil {
			return err
		}
		s.decodeFromStore(idInt, s.encoded[idInt].(map[string]interface{}), rv)
		return nil
	}

	fmt.Println("v, err", parsedJSON, err)
	return nil
	//return json.Unmarshal([]byte(data), value)
}

//todo: i'm sure this function needs to return an error
func (s *Serializer) decodeFromStore(id int, encoded map[string]interface{}, rv reflect.Value) {

	fmt.Println("id, decoded: ", id, s.decoded)

	if len(s.decoded) > id {
		decoded := s.decoded[id]
		fmt.Println("********", decoded.Kind(), rv.Kind())
		//rv.Set(indirect(decoded))
		indirect(rv).Set(decoded)
		return
	}

	_, ok := encoded[KEY_OBJECT_NAME]
	if ok {
		fmt.Println("in the thing2", rv, rv.Kind())
		s.decodeStruct(id, encoded[KEY_OBJECT_PROPERTIES], rv)
		return
	}

	_, ok = encoded[KEY_KEYS]
	if ok {
		s.decodeMap(id, encoded[KEY_KEYS], encoded[KEY_ARRAY], rv)
		return
	} else {
		s.decodeArray(id, encoded[KEY_ARRAY], rv)
		return
	}
}

func (s *Serializer) decodeArray(id int, array interface{}, rv reflect.Value) {

	arr := array.([]interface{})

	var newRV reflect.Value
	switch rv.Kind() {
	case reflect.Interface:
		var newSlice []interface{}
		newRV = reflect.ValueOf(&newSlice)
		fmt.Println("interface(((((((((((((((((((((")
	default:
		fmt.Println("default(((((((((((((((((((((")
		newRV = rv
	}

	fmt.Println("in the thing", rv, rv.Kind())
	s.decoded = insertRV(s.decoded, newRV, id)
	newRV = indirect(newRV)

	for i, arrV := range arr {

		// Get element of array, growing if necessary.
		if newRV.Kind() == reflect.Slice {

			// Grow slice if necessary
			if i >= newRV.Cap() {
				newcap := newRV.Cap() + newRV.Cap()/2
				if newcap < 4 {
					newcap = 4
				}
				newv := reflect.MakeSlice(newRV.Type(), newRV.Len(), newcap)
				reflect.Copy(newv, newRV)
				newRV.Set(newv)
			}
			if i >= newRV.Len() {
				newRV.SetLen(i + 1)
			}
		}

		if i < newRV.Len() {
			// Decode into element.
			s.decodeElement(newRV.Index(i), arrV)
		} else {
			panic("shouldn't get here")
		}
	}

	if rv.CanAddr() {
		rv.Set(s.decoded[id])
	} else {
		rv.Elem().Set(indirect(s.decoded[id]))
	}
}

//todo: should I handle the case in which the KEY_OBJECT_PROPERTIES don't match the struct passed in? seems yes/. actually just do like the json package does
func (s *Serializer) decodeStruct(id int, encodedObject interface{}, v reflect.Value) {

	//todo: this only works with ptr value, make sure that is always going to be the case
	object := encodedObject.(map[string]interface{})

	newV := v
	s.decoded = insertRV(s.decoded, newV, id)
	newV = indirect(newV)

	for key, value := range object {
		field := indirect(newV).FieldByName(key)
		fmt.Println("type of field: ", field)
		fmt.Println("in the thing", field, field.Kind())
		s.decodeElement(field, value)
	}

	if v.CanAddr() {
		v.Set(s.decoded[id])
	} else {
		v.Elem().Set(indirect(s.decoded[id]))
	}
}

func (s *Serializer) decodeMap(id int, keys interface{}, values interface{}, v reflect.Value) {

	//todo: this only works with ptr value, make sure that is always going to be the case
	ks := keys.([]interface{})
	vs := values.([]interface{})

	var newV reflect.Value
	switch indirect(v).Kind() {
	case reflect.Interface:
		newMap := make(map[interface{}]interface{})
		newV = reflect.ValueOf(&newMap)
		fmt.Println("interface(((((((((((((((((((((")
	default:
		fmt.Println("default(((((((((((((((((((((")
		newV = v
	}

	fmt.Println("the thing:::::::::: ", newV, newV.Kind())
	s.decoded = append(s.decoded, newV)
	newV = indirect(newV)

	for i, k := range ks {
		v := vs[i]

		fmt.Println("newV.Type()", newV.Type())

		newKey := reflect.New(newV.Type().Key())
		newValue := reflect.New(newV.Type().Elem())
		s.decodeElement(newKey, k)
		s.decodeElement(newValue, v)
		newV.SetMapIndex(indirect(newKey), indirect(newValue))
	}

	if v.CanAddr() {
		v.Set(indirect(newV))
	} else {
		v.Elem().Set(indirect(newV))
	}
}

func (s *Serializer) decodeElement(rv reflect.Value, value interface{}) {
	fmt.Println("decodeElement value: ", value)

	potentialID, ok := value.(string)
	fmt.Println("str: ", value)
	if ok {
		id, isStoreID := s.storeID(potentialID)
		if isStoreID {
			encoded := s.encoded[id]
			fmt.Println("decodeElement id, encoded", id, encoded)
			s.decodeFromStore(id, encoded.(map[string]interface{}), rv)
			fmt.Printf("returned %+v, %+v\n", rv.Interface(), rv.Kind())
			return
		}
	}

	fmt.Println("rv kind:::::::::::::::: ", rv, rv.Kind())
	//for rv.Kind() == reflect.Ptr || rv.Kind() == reflect.Interface {
	//	rv = rv.Elem()
	//}
	rv = indirect(rv)
	switch rv.Kind() {
	case reflect.Bool:
		rv.SetBool(value.(bool))
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		rv.SetInt(int64(value.(float64)))
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		rv.SetUint(uint64(value.(float64)))
	case reflect.Float32, reflect.Float64:
		//panic("bob")
		//todo: change this
		rv.Set(reflect.ValueOf(value).Convert(rv.Type()))
	case reflect.String:
		rv.SetString(value.(string))
		//todo:
		//case reflect.Uintptr:
		//case reflect.Chan:
		//case reflect.Func:
	case reflect.Interface:
		//todo: shouldn't have tehse here?
		rv.Set(reflect.ValueOf(value))
	case reflect.Ptr:
		rv.Set(reflect.ValueOf(value))
	case reflect.Invalid:
		panic("this should never be invalid")
	case reflect.Array, reflect.Slice:
		//todo? s.decodeLegacyArray(value, rv)
	case reflect.Struct:
		//todo? why do I have to do nothing here?
	//case reflect.Complex64, reflect.Complex128: 	//todo: this is not supported by json i guess?
	//	var c complex128
	//	str := fmt.Sprintf(`"%s"`, value.(string))
	//	err := json.Unmarshal([]byte(str), &c)
	//	if err != nil {
	//		//todo: panic?
	//		panic(err)
	//	}
	//	fld.SetComplex(c)
	//case reflect.UnsafePointer:

	default:
		panic(fmt.Sprintf("unknown kind: %s", rv.Kind()))
	}
}

func (s *Serializer) storeID(str string) (int, bool) {
	if !strings.HasPrefix(str, ID_PREFIX) {
		return 0, false
	}
	id := strings.TrimLeft(str, ID_PREFIX)
	idInt, err := strconv.Atoi(id)
	if err != nil {
		return idInt, false
	}
	return idInt, idInt <= len(s.encoded)
}

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

func insertRV(arr []reflect.Value, value reflect.Value, i int) []reflect.Value {
	if len(arr) > i {
		arr[i] = value
		return arr
	}
	newArr := make([]reflect.Value, i+1)
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
			e := v.Elem()
			//if e.Kind() == reflect.Ptr && !e.IsNil() {
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
