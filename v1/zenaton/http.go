package zenaton

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
)

//todo: figure out what to do with errors

func Get(url string) (*http.Response, error) {
	return http.Get(url)
}

func Post(url string, body interface{}) (*http.Response, error) {
	jsonBody, err := json.Marshal(body)
	if err != nil {
		return nil, err
	}

	fmt.Println(string(jsonBody), url)
	return http.Post(url, "application/json", bytes.NewBuffer(jsonBody))
}

//
//func Put(url string, body interface{})(*http.Response, error){
//	jsonBody, err := json.Marshal(body)
//	if err != nil {
//		return nil, err
//	}
//	return http.(url, "application/json", bytes.NewBuffer(jsonBody))
//}
