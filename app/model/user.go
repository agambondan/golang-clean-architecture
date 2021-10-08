package model

import (
	"github.com/badoux/checkmail"
	"strings"
)

type User struct {
	BaseUUID
	UserAPI
	RoleId   *int64     `json:"role_id,omitempty"`
	Articles *[]Article `json:"articles,omitempty" gorm:"foreignKey:UserID"`
	Role     *Role      `json:"role,omitempty"`
}

type UserAPI struct {
	FirstName    *string `json:"first_name,omitempty"`
	LastName     *string `json:"last_name,omitempty"`
	Email        *string `json:"email,omitempty"`
	PhoneNumber  *string `json:"phone_number,omitempty"`
	Username     *string `json:"username,omitempty"`
	Password     *string `json:"password,omitempty"`
	Instagram    *string `json:"instagram,omitempty"`
	Facebook     *string `json:"facebook,omitempty"`
	Twitter      *string `json:"twitter,omitempty"`
	LinkedIn     *string `json:"linked_in,omitempty"`
	Image        *string `json:"image,omitempty"`
	ImageURL     *string `json:"image_url,omitempty"`
	ThumbnailURL *string `json:"thumbnail_url,omitempty"`
}

type PublicUser struct {
	FirstName    *string `json:"first_name,omitempty"`
	LastName     *string `json:"last_name,omitempty"`
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

// PublicUsers So that we dont expose the user's email address and password to the world
func (users Users) PublicUsers() []interface{} {
	result := make([]interface{}, len(users))
	for index, user := range users {
		result[index] = user.PublicUser()
	}
	return result
}

// PrivateUser So that we dont expose the user's email address and password to the world
func (u *User) PrivateUser() interface{} {
	return &PublicUser{
		FirstName:    u.FirstName,
		LastName:     u.LastName,
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

// PublicUser So that we dont expose the user's email address and password to the world
func (u *User) PublicUser() interface{} {
	return &PublicUser{
		FirstName:    u.FirstName,
		LastName:     u.LastName,
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
