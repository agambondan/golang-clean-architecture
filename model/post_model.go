package model

import (
	"github.com/google/uuid"
	"html"
	"strings"
	"time"
)

type Post struct {
	ID           uint64    `json:"id,omitempty"`
	UserUUID     uuid.UUID `json:"user_id,omitempty"`
	Title        string    `json:"title,omitempty"`
	Description  string    `json:"description,omitempty"`
	Image        string    `json:"image,omitempty"`
	ImageURL     string    `json:"image_url,omitempty"`
	ThumbnailURL string    `json:"thumbnail_url,omitempty"`
	Author       User      `json:"author,omitempty"`
	Categories   []string  `json:"categories,omitempty"`
	CreatedAt    time.Time `json:"created_at,string,omitempty"`
	UpdatedAt    time.Time `json:"updated_at,string,omitempty"`
	DeletedAt    time.Time `sql:"index" json:"deleted_at,string,omitempty"`
}

type PublicPost struct {
	ID           uint64    `json:"id,omitempty"`
	Title        string    `json:"title,omitempty"`
	Description  string    `json:"description,omitempty"`
	Image        string    `json:"image,omitempty"`
	ImageURL     string    `json:"image_url,omitempty"`
	ThumbnailURL string    `json:"thumbnail_url,omitempty"`
	Categories   []string  `json:"categories,omitempty"`
	FirstName    string    `json:"first_name,omitempty"`
	LastName     string    `json:"last_name,omitempty"`
	Username     string    `json:"username,omitempty"`
	Instagram    string    `json:"instagram,omitempty"`
	Facebook     string    `json:"facebook,omitempty"`
	Twitter      string    `json:"twitter,omitempty"`
	LinkedIn     string    `json:"linked_in,omitempty"`
	CreatedAt    time.Time `json:"created_at,omitempty"`
	UpdatedAt    time.Time `json:"updated_at,omitempty"`
}

type Posts []Post

// PublicPosts So that we dont expose the user's email address and password to the world
func (posts Posts) PublicPosts() []interface{} {
	result := make([]interface{}, len(posts))
	for index, post := range posts {
		result[index] = post.PublicPost()
	}
	return result
}

func (p *Post) PublicPost() interface{} {
	return &PublicPost{
		ID:           p.ID,
		Title:        p.Title,
		Description:  p.Description,
		Image:        p.Image,
		ImageURL:     p.ImageURL,
		ThumbnailURL: p.ThumbnailURL,
		Categories:   p.Categories,
		FirstName:    p.Author.FirstName,
		LastName:     p.Author.LastName,
		Username:     p.Author.Username,
		Instagram:    p.Author.Instagram,
		Facebook:     p.Author.Facebook,
		Twitter:      p.Author.Twitter,
		LinkedIn:     p.Author.LinkedIn,
		CreatedAt:    p.CreatedAt,
		UpdatedAt:    p.UpdatedAt,
	}
}

func (p *Post) Prepare() {
	p.Title = html.EscapeString(strings.TrimSpace(p.Title))
	p.Description = html.EscapeString(strings.TrimSpace(p.Description))
	p.Image = html.EscapeString(strings.TrimSpace(p.Image))
	p.ImageURL = html.EscapeString(strings.TrimSpace(p.ImageURL))
	p.ThumbnailURL = html.EscapeString(strings.TrimSpace(p.ThumbnailURL))
	p.CreatedAt = time.Now()
	p.UpdatedAt = time.Now()
}

func (p *Post) Validate(action string) map[string]string {
	var errorMessages = make(map[string]string)
	switch strings.ToLower(action) {
	case "update":
		if p.Title == "" || p.Title == "null" {
			errorMessages["title_required"] = "title is required"
		}
		if p.Description == "" || p.Description == "null" {
			errorMessages["desc_required"] = "description is required"
		}
		if p.ImageURL == "" || p.ThumbnailURL == "" {
			errorMessages["image"] = "image url is required"
		}
	default:
		if p.Title == "" || p.Title == "null" {
			errorMessages["title_required"] = "title is required"
		}
		if p.Description == "" || p.Description == "null" {
			errorMessages["desc_required"] = "description is required"
		}
		if p.UserUUID.String() == "" || p.UserUUID.String() == "null" {
			errorMessages["user_required"] = "user uuid is required"
		}
		if p.ImageURL == "" || p.ThumbnailURL == "" {
			errorMessages["image"] = "image url is required"
		}
	}
	return errorMessages
}
