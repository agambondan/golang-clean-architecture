package model

import (
	"github.com/badoux/checkmail"
	"strings"
)

type User struct {
	BaseUUID
	UserAPI
	BaseImage
	Articles *[]Article `json:"articles,omitempty" gorm:"foreignKey:UserID"`
	Role     *Role      `json:"role,omitempty"`
}

type UserAPI struct {
	FirstName   *string `json:"first_name,omitempty" gorm:"type:varchar(24);not null;"`
	LastName    *string `json:"last_name,omitempty" gorm:"type:varchar(24);not null;"`
	Gender      *string `json:"gender,omitempty" gorm:"type:varchar(24);not null;"`
	Email       *string `json:"email,omitempty" gorm:"type:varchar(64);not null;index:idx_email_deleted_at,unique,where:deleted_at is null"`
	PhoneNumber *string `json:"phone_number,omitempty" gorm:"type:varchar(14);not null;index:idx_phone_number_deleted_at,unique,where:deleted_at is null"`
	Username    *string `json:"username,omitempty" gorm:"type:varchar(36);not null;index:idx_username_deleted_at,unique,where:deleted_at is null"`
	Password    *string `json:"password,omitempty" gorm:"type:varchar(256);not null;"`
	Instagram   *string `json:"instagram,omitempty" gorm:"type:varchar(24)"`
	Facebook    *string `json:"facebook,omitempty" gorm:"type:varchar(24)"`
	Twitter     *string `json:"twitter,omitempty" gorm:"type:varchar(24)"`
	LinkedIn    *string `json:"linked_in,omitempty" gorm:"type:varchar(24)"`
	RoleId      *int64  `json:"role_id,omitempty" gorm:"type:smallint;not null;"`
}

type PublicUser struct {
	FirstName    *string `json:"first_name,omitempty"`
	LastName     *string `json:"last_name,omitempty"`
	Gender       *string `json:"gender,omitempty"`
	Username     *string `json:"username,omitempty"`
	Instagram    *string `json:"instagram,omitempty"`
	Facebook     *string `json:"facebook,omitempty"`
	Twitter      *string `json:"twitter,omitempty"`
	LinkedIn     *string `json:"linked_in,omitempty"`
	Image        *string `json:"image,omitempty"`
	ImageURL     *string `json:"image_url,omitempty"`
	ThumbnailURL *string `json:"thumbnail_url,omitempty"`
}

type Users []User

// PublicUsers So that we don't expose the user's email address and password to the world
func (users Users) PublicUsers() []interface{} {
	result := make([]interface{}, len(users))
	for index, user := range users {
		result[index] = user.PublicUser()
	}
	return result
}

// PrivateUsers So that we don't expose the user's email address and password to the world
func (users Users) PrivateUsers() []interface{} {
	result := make([]interface{}, len(users))
	for index, user := range users {
		result[index] = user.PrivateUser()
	}
	return result
}

// PrivateUser So that we don't expose the user's email address and password to the world
func (u *User) PrivateUser() interface{} {
	return &User{
		BaseUUID: BaseUUID{
			ID: u.ID,
			BaseDate: BaseDate{
				CreatedAt: u.CreatedAt,
				UpdatedAt: u.UpdatedAt,
			},
		},
		UserAPI: UserAPI{
			FirstName:   u.FirstName,
			LastName:    u.LastName,
			Gender:      u.Gender,
			Email:       u.Email,
			PhoneNumber: u.PhoneNumber,
			Username:    u.Username,
			Password:    u.Password,
			Instagram:   u.Instagram,
			Facebook:    u.Facebook,
			Twitter:     u.Twitter,
			LinkedIn:    u.LinkedIn,
			RoleId:      u.RoleId,
		},
		BaseImage: BaseImage{
			Image:        u.Image,
			ImageURL:     u.ImageURL,
			ThumbnailURL: u.ThumbnailURL,
		},
	}
}

// PublicUser So that we don't expose the user's email address and password to the world
func (u *User) PublicUser() interface{} {
	return &PublicUser{
		FirstName:    u.FirstName,
		LastName:     u.LastName,
		Gender:       u.Gender,
		Username:     u.Username,
		Instagram:    u.Instagram,
		Facebook:     u.Facebook,
		Twitter:      u.Twitter,
		LinkedIn:     u.LinkedIn,
		Image:        u.Image,
		ImageURL:     u.ImageURL,
		ThumbnailURL: u.ThumbnailURL,
	}
}

func (u *User) Validate(action string) map[string]string {
	var errorMessages = make(map[string]string)
	var err error
	switch strings.ToLower(action) {
	case "images":
		if u.ImageURL != nil || u.ThumbnailURL != nil {
			if *u.ImageURL == "" || len(*u.ImageURL) < 45 || *u.ThumbnailURL == "" {
				errorMessages["image"] = "image url is required"
			}
		}
	case "update":
		if u.Email != nil {
			if *u.Email == "" {
				errorMessages["email_required"] = "email required"
				return errorMessages
			}
			if *u.Email != "" {
				if err = checkmail.ValidateFormat(*u.Email); err != nil {
					errorMessages["invalid_email"] = "email email"
					return errorMessages
				}
			}
		}
	case "login":
		if u.Password != nil {
			if *u.Password == "" {
				errorMessages["password_required"] = "password is required"
				return errorMessages
			}
		}
		if u.Email != nil {
			if u.Username != nil {
				if *u.Email == "" && *u.Username == "" {
					errorMessages["invalid_login"] = "email or username is required"
					return errorMessages
				}
			}
			if *u.Email != "" {
				if err = checkmail.ValidateFormat(*u.Email); err != nil {
					errorMessages["invalid_email"] = "please provide a valid email"
					return errorMessages
				}
			}
		}
	case "forgot_password":
		if u.Email != nil {
			if *u.Email == "" {
				errorMessages["email_required"] = "email required"
				return errorMessages
			}
		}
		if u.Email != nil {
			if *u.Email != "" {
				if err = checkmail.ValidateFormat(*u.Email); err != nil {
					errorMessages["invalid_email"] = "please provide a valid email"
				}
			}
		}
	default:
		if u.FirstName != nil {
			if *u.FirstName == "" {
				errorMessages["firstname_required"] = "first name is required"
				return errorMessages
			}
		}
		if u.LastName != nil {
			if *u.LastName == "" {
				errorMessages["lastname_required"] = "last name is required"
				return errorMessages
			}
		}
		if u.Password != nil {
			if *u.Password == "" {
				errorMessages["password_required"] = "password is required"
				return errorMessages
			}
			if *u.Password != "" && len(*u.Password) < 6 {
				errorMessages["invalid_password"] = "password should be at least 6 characters"
				return errorMessages
			}
		}
		if u.PhoneNumber != nil {
			if *u.PhoneNumber == "" {
				errorMessages["phone_number_required"] = "phone number is required"
				return errorMessages
			}
		}
		if u.RoleId != nil {
			if *u.RoleId == 0 {
				errorMessages["role_id_required"] = "role id is required"
				return errorMessages
			}
		}
		if u.Email != nil {
			if *u.Email == "" {
				errorMessages["email_required"] = "email is required"
				return errorMessages
			}
			if *u.Email != "" {
				if err = checkmail.ValidateFormat(*u.Email); err != nil {
					errorMessages["invalid_email"] = "please provide a valid email"
					return errorMessages
				}
			}
		}
	}
	return nil
}
