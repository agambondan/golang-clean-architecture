package repository

import (
	"database/sql"
	"fmt"
	"github.com/google/uuid"
	"golang-youtube-api/model"
	"strings"
)

type PostRepository interface {
	Save(post *model.Post) (*model.Post, error)
	FindAll() ([]model.Post, error)
	FindById(id uint64) (model.Post, error)
	FindAllByUserId(uuid uuid.UUID) ([]model.Post, error)
	FindAllByCategoryId(id uint64) ([]model.Post, error)
	UpdateById(id uint64, post *model.Post) (*model.Post, error)
	DeleteById(id uint64) error
}

type postRepo struct {
	db *sql.DB
}

func NewPostRepository(db *sql.DB) PostRepository {
	return &postRepo{db}
}

func (r *postRepo) Save(post *model.Post) (*model.Post, error) {
	panic("implement me")
}

func (r *postRepo) FindAll() ([]model.Post, error) {
	var posts []model.Post
	var arrayCategory []uint8
	queryGetUsers := fmt.Sprintf("select p.id , p.title , p.description , p.post_images , u.first_name , u.last_name , array_agg(distinct c.\"name\") as categories , p.created_at , p.updated_at  from posts p join users u on p.user_uuid = u.uuid join post_categories pc on p.id = pc.post_id join categories c on pc.category_id = c.id group by p.id, p.title , p.description , p.post_images , u.first_name , u.last_name , p.created_at , p.updated_at ;")
	prepare, err := r.db.Prepare(queryGetUsers)
	if err != nil {
		return posts, err
	}
	rows, err := prepare.Query()
	if err != nil {
		return posts, err
	}
	for rows.Next() {
		var post model.Post
		err := rows.Scan(&post.ID, &post.Title, &post.Description, &post.PostImage, &post.Author.FirstName, &post.Author.LastName, &arrayCategory, &post.CreatedAt, &post.UpdatedAt)
		if err != nil {
			return posts, err
		}
		lenStringArrayCategory := len(string(arrayCategory))
		splitCategories := strings.Split(string(arrayCategory)[1:lenStringArrayCategory-1], ",")
		for i := 0; i < len(splitCategories); i++ {
			post.Categories = append(post.Categories, splitCategories[i])
		}
		posts = append(posts, post)
	}
	return posts, err
}

func (r *postRepo) FindById(id uint64) (model.Post, error) {
	panic("implement me")
}

func (r *postRepo) FindAllByUserId(uuid uuid.UUID) ([]model.Post, error) {
	panic("implement me")
}

func (r *postRepo) FindAllByCategoryId(id uint64) ([]model.Post, error) {
	panic("implement me")
}


func (r *postRepo) UpdateById(id uint64, post *model.Post) (*model.Post, error) {
	panic("implement me")
}

func (r *postRepo) DeleteById(id uint64) error {
	queryInsert := fmt.Sprintf("DELETE FROM %s where id = %d", "posts", id)
	_, err := r.db.Prepare(queryInsert)
	if err != nil {
		return err
	}
	return err
}
