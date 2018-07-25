package zenaton

import "encoding/json"

type Serializer struct{}

//todo: is this enough? what about the recursive case?
func (s Serializer) Encode(data interface{}) (string, error) {
	jsonData, err := json.Marshal(data)
	return string(jsonData), err
}

func (s Serializer) Decode(data string, value *interface{}) error {
	err := json.Unmarshal([]byte(data), value)
	return err
}
