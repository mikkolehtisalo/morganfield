package morganfield

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"reflect"
)

// Optional validation functionality
type Validator interface {
	Validate()
}

// Marshals any object into string
func Marshal(in interface{}) (string, error) {
	arr, err := json.Marshal(in)
	return string(arr), err
}

// Unmarshals any string into object
func UnMarshal(in string, out interface{}) (interface{}, error) {
	err := json.Unmarshal([]byte(in), out)
	return out, err
}

// filters input json
func filter_input_json(r *http.Request, s Service_Definition) (string, error) {
	defer r.Body.Close()
	return filter_object_json(&r.Body, s.In_Object)
}

// filters output json
func filter_output_json(r *http.Response, s Service_Definition) (string, error) {
	defer r.Body.Close()
	return filter_object_json(&r.Body, s.Out_Object)
}

// filters JSON by unmarshaling & marshaling it
func filter_object_json(in *io.ReadCloser, o interface{}) (string, error) {
	var filterd string

	if o == nil {
		return filterd, fmt.Errorf("filter_object_json: no Object set!")
	}

	// Build new o object
	tgtType := reflect.ValueOf(o).Type()
	target := reflect.New(tgtType).Interface()

	// Read the string representation
	bs, err := ioutil.ReadAll(*in)
	if err != nil {
		panic(fmt.Sprintf("filter_object_json: %v", err))
	}
	tr := string(bs)

	ino, err := UnMarshal(tr, target)
	if err != nil {
		panic(fmt.Sprintf("filter_object_json: %v", err))
	}

	// Optional extra validation
	if validd, ok := ino.(Validator); ok {
		validd.Validate()
	}

	// Marshal the result back to string
	cleand, err := Marshal(target)
	if err != nil {
		panic(fmt.Sprintf("filter_object_json: %v", err))
	}

	filterd = string(cleand)

	return filterd, nil
}
