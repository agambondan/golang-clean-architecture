package lib

import (
	"encoding/json"
	"fmt"
)

// Merge a struct to another struct
func Merge(from interface{}, to interface{}) (err error) {
	typeFrom := fmt.Sprintf("%T", from)
	if typeFrom == "string" {
		err = json.Unmarshal([]byte(from.(string)), &to)
		return err
	}
	var bytes []byte
	if typeFrom == "[]uint8" {
		bytes = from.([]byte)
	} else {
		bytes, err = json.Marshal(from)
		if err != nil {
			return err
		}
	}
	err = json.Unmarshal(bytes, &to)
	return err
}
