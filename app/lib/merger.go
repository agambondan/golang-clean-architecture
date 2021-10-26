package lib

import (
	"encoding/json"
	"fmt"
)

// Merge a struct to another struct
func Merge(from interface{}, to interface{}) error {
	var err error
	var j []byte
	if fmt.Sprintf("%T", from) == "[]byte" {
		j = from.([]byte)
	} else {
		j, err = json.Marshal(from)
		return err
	}
	err = json.Unmarshal(j, to)

	return err
}
