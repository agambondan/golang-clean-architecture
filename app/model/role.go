package model

import (
	"strings"
)

type Role struct {
	BaseInt
	Name *string `json:"name"`
}

func (r *Role) Validate(action string) map[string]string {
	var errorMessages = make(map[string]string)
	switch strings.ToLower(action) {
	case "update":
		if r.Name != nil {
			if *r.Name == "" || *r.Name == "null" {
				errorMessages["title_required"] = "title is required"
			}
		}
	default:
		if r.Name != nil {
			if *r.Name == "" || *r.Name == "null" {
				errorMessages["title_required"] = "title is required"
			}
		}
	}
	return errorMessages
}
