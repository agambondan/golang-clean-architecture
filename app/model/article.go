package model

import (
	"github.com/google/uuid"
	"strings"
	"time"
)

type Article struct {
	BaseInt
	ArticleAPI
	BaseImage
	UserID     *uuid.UUID  `json:"user_id,omitempty"`
	User     *User       `json:"author,omitempty"`
	Category *[]Category `json:"categories,omitempty" gorm:"many2many:article_categories"`
}

type ArticleAPI struct {
	Title       *string `json:"title,omitempty" gorm:"type:varchar(256);not null;index:idx_title_deleted_at,unique,where:deleted_at is null"`
	Description *string `json:"description,omitempty" gorm:"type:text;not null"`
}

type PublicArticle struct {
	Title        *string     `json:"title,omitempty"`
	Description  *string     `json:"description,omitempty"`
	Image        *string     `json:"image,omitempty"`
	ImageURL     *string     `json:"image_url,omitempty"`
	ThumbnailURL *string     `json:"thumbnail_url,omitempty"`
	FirstName    *string     `json:"first_name,omitempty"`
	LastName     *string     `json:"last_name,omitempty"`
	UserImage    *string     `json:"user_image,omitempty"`
	UserImageURL *string     `json:"user_image_url,omitempty"`
	Username     *string     `json:"username,omitempty"`
	Instagram    *string     `json:"instagram,omitempty"`
	Facebook     *string     `json:"facebook,omitempty"`
	Twitter      *string     `json:"twitter,omitempty"`
	LinkedIn     *string     `json:"linked_in,omitempty"`
	CreatedAt    *time.Time  `json:"created_at,omitempty"`
	UpdatedAt    *time.Time  `json:"updated_at,omitempty"`
	Categories   *[]Category `json:"categories,omitempty"`
}

type Articles []Article

// PublicArticles So that we dont expose the user's email address and password to the world
func (posts Articles) PublicArticles() []interface{} {
	result := make([]interface{}, len(posts))
	for index, post := range posts {
		result[index] = post.PublicArticle()
	}
	return result
}

func (p *Article) PublicArticle() interface{} {
	return &PublicArticle{
		Title:        p.Title,
		Description:  p.Description,
		Image:        p.Image,
		ImageURL:     p.ImageURL,
		ThumbnailURL: p.ThumbnailURL,
		FirstName:    p.User.FirstName,
		LastName:     p.User.LastName,
		Username:     p.User.Username,
		UserImage:    p.User.Image,
		UserImageURL: p.User.ImageURL,
		Instagram:    p.User.Instagram,
		Facebook:     p.User.Facebook,
		Twitter:      p.User.Twitter,
		LinkedIn:     p.User.LinkedIn,
		CreatedAt:    p.CreatedAt,
		UpdatedAt:    p.UpdatedAt,
		Categories:   p.Category,
	}
}

func (p *Article) Validate(action string) map[string]string {
	var errorMessages = make(map[string]string)
	switch strings.ToLower(action) {
	case "images":
		if p.ImageURL != nil || p.ThumbnailURL != nil {
			if *p.ImageURL == "" || len(*p.ImageURL) < 45 || *p.ThumbnailURL == "" {
				errorMessages["image"] = "image url is required"
			}
		}
	default:
		if p.Title != nil {
			if *p.Title == "" {
				errorMessages["title_required"] = "title is required"
			}
		} else if p.Description != nil {
			if *p.Description == "" {
				errorMessages["desc_required"] = "description is required"
			}
		} else if p.UserID != nil {
			if p.UserID.String() == "" || p.UserID.String() == "null" {
				errorMessages["user_required"] = "user uuid is required"
			}
		}
	}
	return errorMessages
}
