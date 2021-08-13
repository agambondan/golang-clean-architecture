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
	FindAll(limit, offset int) ([]model.Post, error)
	FindById(id uint64) (model.Post, error)
	FindByTitle(title string) (model.Post, error)
	FindAllByUserId(uuid uuid.UUID) ([]model.Post, error)
	FindAllByUsername(username string) ([]model.Post, error)
	FindAllByCategoryName(name string, limit, offset int) ([]model.Post, error)
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
	query := fmt.Sprintf("insert into %s (title, description, image, image_url, thumbnail_url, user_uuid, created_at, updated_at, deleted_at) "+
		"VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9) returning id", "posts")
	stmt, err := r.db.Prepare(query)
	if err != nil {
		return post, err
	}
	err = stmt.QueryRow(&post.Title, &post.Description, &post.Image, &post.ImageURL, &post.ThumbnailURL, &post.UserUUID, &post.CreatedAt, &post.UpdatedAt, nil).Scan(&post.ID)
	if err != nil {
		return post, err
	}
	return post, err
}

func (r *postRepo) FindAll(limit, offset int) ([]model.Post, error) {
	var posts []model.Post
	var arrayCategory []uint8
	query := fmt.Sprintf("select p.id, p.title, p.description, p.image, p.image_url, p.thumbnail, u.first_name, "+
		"u.last_name, u.instagram, u.facebook, u.twitter, u.linkedin, array_agg(distinct c.\"name\") as categories, p.created_at, p.updated_at from posts p "+
		"join users u on p.user_uuid = u.uuid join post_categories pc on p.id = pc.post_id "+
		"join categories c on pc.category_id = c.id where p.deleted_at is null "+
		"group by p.id, p.title, p.description, p.image, p.image_url, p.thumbnail, u.first_name, u.last_name, u.instagram, u.facebook, u.twitter, u.linkedin, p.created_at, p.updated_at "+
		"order by p.id limit %d offset %d ;", limit, offset)
	prepare, err := r.db.Prepare(query)
	if err != nil {
		return posts, err
	}
	rows, err := prepare.Query()
	if err != nil {
		return posts, err
	}
	for rows.Next() {
		var post model.Post
		err = rows.Scan(&post.ID, &post.Title, &post.Description, &post.Image, &post.ImageURL, &post.ThumbnailURL, &post.Author.FirstName, &post.Author.LastName,
			&post.Author.Instagram, &post.Author.Facebook, &post.Author.Twitter, &post.Author.LinkedIn, &arrayCategory, &post.CreatedAt, &post.UpdatedAt)
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
	var post model.Post
	var arrayCategory []uint8
	query := fmt.Sprintf("select p.id , p.title , p.description , p.image, p.image_url, p.thumbnail , u.uuid, u.first_name ," +
		" u.last_name , array_agg(distinct c.\"name\") as categories , p.created_at , p.updated_at from posts p " +
		"join users u on p.user_uuid = u.uuid join post_categories pc on p.id = pc.post_id " +
		"join categories c on pc.category_id = c.id where p.id=$1 and p.deleted_at is null " +
		"group by p.id, p.title , p.description , p.image, p.image_url, p.thumbnail , u.uuid, u.first_name , u.last_name , p.created_at , p.updated_at;")
	prepare, err := r.db.Prepare(query)
	if err != nil {
		return post, err
	}
	err = prepare.QueryRow(&id).Scan(&post.ID, &post.Title, &post.Description, &post.Image, &post.ImageURL, &post.ThumbnailURL,
		&post.UserUUID, &post.Author.FirstName, &post.Author.LastName, &arrayCategory, &post.CreatedAt, &post.UpdatedAt)
	if err != nil {
		return post, err
	}
	lenStringArrayCategory := len(string(arrayCategory))
	splitCategories := strings.Split(string(arrayCategory)[1:lenStringArrayCategory-1], ",")
	for i := 0; i < len(splitCategories); i++ {
		post.Categories = append(post.Categories, splitCategories[i])
	}
	return post, nil
}

func (r *postRepo) FindByTitle(title string) (model.Post, error) {
	var post model.Post
	var arrayCategory []uint8
	query := fmt.Sprintf("select p.id , p.title , p.description , p.image, p.image_url, p.thumbnail , u.uuid, u.first_name ,"+
		" u.last_name , array_agg(distinct c.\"name\") as categories , p.created_at , p.updated_at from posts p "+
		"join users u on p.user_uuid = u.uuid join post_categories pc on p.id = pc.post_id "+
		"join categories c on pc.category_id = c.id where p.title like '%s' and p.deleted_at is null "+
		"group by p.id, p.title , p.description , p.image, p.image_url, p.thumbnail , u.uuid, u.first_name , u.last_name , p.created_at , p.updated_at;", title)
	prepare, err := r.db.Prepare(query)
	if err != nil {
		return post, err
	}
	err = prepare.QueryRow().Scan(&post.ID, &post.Title, &post.Description, &post.Image, &post.ImageURL, &post.ThumbnailURL,
		&post.UserUUID, &post.Author.FirstName, &post.Author.LastName, &arrayCategory, &post.CreatedAt, &post.UpdatedAt)
	if err != nil {
		return post, err
	}
	lenStringArrayCategory := len(string(arrayCategory))
	splitCategories := strings.Split(string(arrayCategory)[1:lenStringArrayCategory-1], ",")
	for i := 0; i < len(splitCategories); i++ {
		post.Categories = append(post.Categories, splitCategories[i])
	}
	return post, nil
}

func (r *postRepo) FindAllByUserId(uuid uuid.UUID) ([]model.Post, error) {
	var posts []model.Post
	var arrayCategory []uint8
	query := fmt.Sprintf("select p.id , p.title , p.description , p.image, p.image_url, p.thumbnail , p.user_uuid, u.first_name ," +
		" u.last_name , array_agg(distinct c.\"name\") as categories , p.created_at , p.updated_at from posts p " +
		"join users u on p.user_uuid = u.uuid join post_categories pc on p.id = pc.post_id " +
		"join categories c on pc.category_id = c.id where p.user_uuid = $1 and p.deleted_at is null " +
		"group by p.id, p.title , p.description , p.image, p.image_url, p.thumbnail , u.uuid, u.first_name , u.last_name , p.created_at , p.updated_at;")
	prepare, err := r.db.Prepare(query)
	if err != nil {
		return posts, err
	}
	rows, err := prepare.Query(&uuid)
	if err != nil {
		return posts, err
	}
	for rows.Next() {
		var post model.Post
		err := rows.Scan(&post.ID, &post.Title, &post.Description, &post.Image, &post.ImageURL, &post.ThumbnailURL,
			&post.UserUUID, &post.Author.FirstName, &post.Author.LastName, &arrayCategory, &post.CreatedAt, &post.UpdatedAt)
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

func (r *postRepo) FindAllByUsername(username string) ([]model.Post, error) {
	var posts []model.Post
	var arrayCategory []uint8
	query := fmt.Sprintf("select p.id , p.title , p.description , p.image, p.image_url, p.thumbnail , u.first_name ," +
		" u.last_name , array_agg(distinct c.\"name\") as categories , p.created_at , p.updated_at from posts p " +
		"join users u on p.user_uuid = u.uuid join post_categories pc on p.id = pc.post_id " +
		"join categories c on pc.category_id = c.id where u.username = $1 and p.deleted_at is null " +
		"group by p.id, p.title , p.description , p.image, p.image_url, p.thumbnail , u.first_name , u.last_name , p.created_at , p.updated_at;")
	prepare, err := r.db.Prepare(query)
	if err != nil {
		return posts, err
	}
	rows, err := prepare.Query(&username)
	if err != nil {
		return posts, err
	}
	for rows.Next() {
		var post model.Post
		err := rows.Scan(&post.ID, &post.Title, &post.Description, &post.Image, &post.ImageURL, &post.ThumbnailURL,
			&post.Author.FirstName, &post.Author.LastName, &arrayCategory, &post.CreatedAt, &post.UpdatedAt)
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

func (r *postRepo) FindAllByCategoryName(name string, limit, offset int) ([]model.Post, error) {
	var posts []model.Post
	var arrayCategory []uint8
	query := fmt.Sprintf("select p.id , p.title , p.description , p.image, p.image_url, p.thumbnail , u.first_name ,"+
		" u.last_name , array_agg(distinct c.\"name\") as categories , p.created_at , p.updated_at from posts p "+
		"join users u on p.user_uuid = u.uuid join post_categories pc on p.id = pc.post_id "+
		"join categories c on pc.category_id = c.id where c.name = '%s' and p.deleted_at is null "+
		"group by p.id, p.title , p.description , p.image, p.image_url, p.thumbnail , u.first_name , u.last_name , p.created_at , p.updated_at "+
		"order by p.id limit %d offset %d ;", name, limit, offset)
	prepare, err := r.db.Prepare(query)
	if err != nil {
		return posts, err
	}
	rows, err := prepare.Query()
	if err != nil {
		return posts, err
	}
	for rows.Next() {
		var post model.Post
		err = rows.Scan(&post.ID, &post.Title, &post.Description, &post.Image, &post.ImageURL, &post.ThumbnailURL,
			&post.Author.FirstName, &post.Author.LastName, &arrayCategory, &post.CreatedAt, &post.UpdatedAt)
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

func (r *postRepo) UpdateById(id uint64, post *model.Post) (*model.Post, error) {
	query := fmt.Sprintf("update posts set title = $1, description = $2, image = $3, image_url = $4, thumbnail_url = $5, updated_at = $6 where id = %d", id)
	stmt, err := r.db.Prepare(query)
	if err != nil {
		return post, err
	}
	_, err = stmt.Exec(&post.Title, &post.Description, &post.Image, &post.ImageURL, &post.ThumbnailURL, &post.UpdatedAt)
	if err != nil {
		return post, err
	}
	return post, err
}

func (r *postRepo) DeleteById(id uint64) error {
	queryInsert := fmt.Sprintf("DELETE FROM %s where id = %d", "posts", id)
	_, err := r.db.Prepare(queryInsert)
	if err != nil {
		return err
	}
	return err
}
