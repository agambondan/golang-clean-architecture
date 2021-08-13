package model

import (
	"html"
	"strings"
	"time"
)

type Category struct {
	ID           uint64    `sql:"primary_key" json:"id,omitempty"`
	Name         string    `json:"name,omitempty"`
	Image        string    `json:"image,omitempty"`
	ImageURL     string    `json:"image_url"`
	ThumbnailURL string    `json:"thumbnail,omitempty"`
	Posts        []Post    `json:"categories,omitempty"`
	CreatedAt    time.Time `json:"created_at,string,omitempty"`
	UpdatedAt    time.Time `json:"updated_at,string,omitempty"`
	DeletedAt    time.Time `sql:"index" json:"deleted_at,string,omitempty"`
}

type Categories []Category

// Categories So that we dont expose the user's email address and password to the world
func (categories Categories) Categories() []interface{} {
	result := make([]interface{}, len(categories))
	for index, category := range categories {
		result[index] = category.Category()
	}
	return result
}

func (c *Category) Category() interface{} {
	return &Category{
		ID:           c.ID,
		Name:         c.Name,
		Image:        c.Image,
		ImageURL:     c.ImageURL,
		ThumbnailURL: c.ThumbnailURL,
		CreatedAt:    c.CreatedAt,
		UpdatedAt:    c.UpdatedAt,
	}
}

func (c *Category) Prepare() {
	c.Name = html.EscapeString(strings.TrimSpace(c.Name))
	c.Image = html.EscapeString(strings.TrimSpace(c.Image))
	c.ImageURL = html.EscapeString(strings.TrimSpace(c.ImageURL))
	c.ThumbnailURL = html.EscapeString(strings.TrimSpace(c.ThumbnailURL))
	c.CreatedAt = time.Now()
	c.UpdatedAt = time.Now()
}

func (c *Category) Validate(action string) map[string]string {
	var errorMessages = make(map[string]string)

	switch strings.ToLower(action) {
	case "update":
		if c.Name == "" || c.Name == "null" {
			errorMessages["name_required"] = "name is required"
		}
	default:
		if c.Name == "" || c.Name == "null" {
			errorMessages["name_required"] = "name is required"
		}
	}
	return errorMessages
}
