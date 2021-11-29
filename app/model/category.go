package model

import (
	"strings"
)

type Category struct {
	BaseInt
	CategoryAPI
	BaseImage
	Articles *[]Article `json:"articles,omitempty" gorm:"many2many:article_categories"`
}

type CategoryAPI struct {
	Name *string `json:"name,omitempty" gorm:"type:varchar(16);index:idx_category_name_deleted_at,unique,where:deleted_at is null"`
}

func (c *Category) Validate(action string) map[string]string {
	var errorMessages = make(map[string]string)
	switch strings.ToLower(action) {
	case "update":
		if c.Name != nil {
			if *c.Name == "" || *c.Name == "null" {
				errorMessages["name_required"] = "name is required"
				return errorMessages
			}
		}
	default:
		if c.Name != nil {
			if *c.Name == "" || *c.Name == "null" {
				errorMessages["name_required"] = "name is required"
				return errorMessages
			}
		}
	}
	return errorMessages
}
