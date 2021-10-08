package migrations

import (
	"go-blog-api/app/model"
)

// ModelMigrations models to migrate
var ModelMigrations []interface{} = []interface{}{
	&model.Role{},
	&model.User{},
	&model.Category{},
	&model.Article{},
}
