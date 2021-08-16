package repository

import (
	"database/sql"
	"fmt"
	"golang-youtube-api/model"
)

type PostCategoryRepository interface {
	Save(post *model.PostCategory) (*model.PostCategory, error)
	FindAll() ([]model.PostCategory, error)
	UpdateById(id uint64, post *model.PostCategory) (*model.PostCategory, error)
	DeleteById(id uint64) error
	Count() (int, error)
}

type postCategoryRepo struct {
	db *sql.DB
}

func NewPostCategoryRepository(db *sql.DB) PostCategoryRepository {
	return &postCategoryRepo{db}
}

func (p *postCategoryRepo) Save(post *model.PostCategory) (*model.PostCategory, error) {
	queryInsert := fmt.Sprintf("insert into %s (post_id, category_id) "+
		"VALUES ($1, $2)", "post_categories")
	stmt, err := p.db.Prepare(queryInsert)
	if err != nil {
		return post, err
	}
	_, err = stmt.Exec(&post.PostID, &post.CategoryID)
	if err != nil {
		return post, err
	}
	return post, err
}

func (p *postCategoryRepo) FindAll() ([]model.PostCategory, error) {
	var posts []model.PostCategory
	var post model.PostCategory
	query := fmt.Sprintf("select * from post_categories")
	rows, err := p.db.Query(query)
	if err != nil {
		return posts, err
	}
	for rows.Next() {
		err = rows.Scan(&post.PostID, &post.CategoryID)
		if err != nil {
			return posts, err
		}
		posts = append(posts, post)
	}
	return posts, err
}

func (p *postCategoryRepo) UpdateById(id uint64, post *model.PostCategory) (*model.PostCategory, error) {
	query := fmt.Sprintf("update post_categories set category_id = $1 where post_id = %d", id)
	stmt, err := p.db.Prepare(query)
	if err != nil {
		return post, err
	}
	_, err = stmt.Exec(&post.CategoryID)
	if err != nil {
		return post, err
	}
	return post, err
}

func (p *postCategoryRepo) DeleteById(id uint64) error {
	query := fmt.Sprintf("DELETE FROM %s where post_id = %d", "post_categories", id)
	_, err := p.db.Prepare(query)
	if err != nil {
		return err
	}
	return err
}


func (p *postCategoryRepo) Count() (int, error) {
	var count int
	queryInsert := fmt.Sprintf("select count(post_id) from post_categories")
	prepare, err := p.db.Prepare(queryInsert)
	if err != nil {
		return count, err
	}
	err = prepare.QueryRow().Scan(&count)
	if err != nil {
		return count, err
	}
	return count, err
}