package model

import (
	"strings"
)

type Category struct {
	BaseInt
	Name         *string    `json:"name,omitempty"`
	Image        *string    `json:"image,omitempty"`
	ImageURL     *string    `json:"image_url,omitempty"`
	ThumbnailURL *string    `json:"thumbnail,omitempty"`
	Articles     *[]Article `json:"articles,omitempty" gorm:"many2many:article_categories"`
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
