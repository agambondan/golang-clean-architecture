package repository

import (
	"github.com/google/uuid"
	"go-blog-api/app/model"
	"gorm.io/gorm"
)

type ArticleRepository interface {
	Save(article *model.Article) (*model.Article, error)
	FindAll(limit, offset int) (*[]model.Article, error)
	FindById(id int64) (*model.Article, error)
	FindByTitle(title string) (*model.Article, error)
	FindAllByUserId(uuid uuid.UUID, limit, offset int) (*[]model.Article, error)
	FindAllByUsername(username string, limit, offset int) (*[]model.Article, error)
	FindAllByCategoryName(name string, limit, offset int) (*[]model.Article, error)
	CountByCategoryName(name string) (int64, error)
	UpdateById(id int64, article *model.Article) (*model.Article, error)
	DeleteById(id int64) error
	Count() (int64, error)
}

type articleRepo struct {
	db *gorm.DB
}

func NewArticleRepository(db *gorm.DB) ArticleRepository {
	return &articleRepo{db}
}

func (a *articleRepo) Save(article *model.Article) (*model.Article, error) {
	if tx := a.db.Create(&article); tx.Error != nil {
		return article, tx.Error
	}
	if err := a.db.Model(&article).Association("Category").Append(&article.Categories); err != nil {
		return article, err
	}
	return article, nil
}

func (a *articleRepo) FindAll(limit, offset int) (*[]model.Article, error) {
	var articles *[]model.Article
	if tx := a.db.
		Find(&articles).Select("*").Offset(offset).Limit(limit); tx.Error != nil {
		return articles, tx.Error
	}
	for i, article := range *articles {
		var categories *[]model.Category
		err := a.db.Model(&article).Select("category.name").Association("Categories").Find(&categories)
		if err != nil {
			return articles, err
		}
		(*articles)[i].Categories = categories
	}
	return articles, nil
}

func (a *articleRepo) FindById(id int64) (*model.Article, error) {
	var article *model.Article
	if tx := a.db.Model(&model.Article{}).Joins("join \"user\" u on article.user_id = u.id").
		Joins("join article_categories a_c on article.id = a_c.article_id").
		Joins("join category c on a_c.category_id = c.id").
		Preload("Categories").
		First(&article, id); tx.Error != nil {
		return article, tx.Error
	}
	return article, nil
}

func (a *articleRepo) FindByTitle(title string) (*model.Article, error) {
	var article *model.Article
	if tx := a.db.
		//Model(&model.Article{}).
		//Joins("join user u on article.user_id = u.id").
		//Joins("join article_categories a_c on article.id = a_c.article_id").
		//Joins("join category c on a_c.category_id = c.id").
		Preload("Categories").
		First(&article, "article.title = ?", title); tx.Error != nil {
		return article, tx.Error
	}
	return article, nil
}

func (a *articleRepo) FindAllByUserId(uuid uuid.UUID, limit, offset int) (*[]model.Article, error) {
	var articles *[]model.Article
	if tx := a.db.
		//Model(&model.Article{}).
		//Joins("join user u on article.user_id = u.id").
		//Joins("join article_categories a_c on article.id = a_c.article_id").
		//Joins("join category c on a_c.category_id = c.id").
		Preload("Categories").
		Find(&articles, "article.user_id = ?", uuid).Offset(offset).Limit(limit); tx.Error != nil {
		return articles, tx.Error
	}
	return articles, nil
}

func (a *articleRepo) FindAllByUsername(username string, limit, offset int) (*[]model.Article, error) {
	var articles *[]model.Article
	if tx := a.db.
		//Model(&model.Article{}).
		//Joins("join user u on article.user_id = u.id").
		//Joins("join article_categories a_c on article.id = a_c.article_id").
		//Joins("join category c on a_c.category_id = c.id").
		Preload("Categories").
		Find(&articles, "u.username = ?", username).Offset(offset).Limit(limit); tx.Error != nil {
		return articles, tx.Error
	}
	return articles, nil
}

func (a *articleRepo) FindAllByCategoryName(name string, limit, offset int) (*[]model.Article, error) {
	var articles *[]model.Article
	if tx := a.db.
		//Model(&model.Article{}).
		//Joins("join user u on article.user_id = u.id").
		//Joins("join article_categories a_c on article.id = a_c.article_id").
		//Joins("join category c on a_c.category_id = c.id").
		Preload("Categories").
		Preload("Category").
		Find(&articles, "c.name = ?", name).Offset(offset).Limit(limit); tx.Error != nil {
		return articles, tx.Error
	}
	return articles, nil
}

func (a *articleRepo) CountByCategoryName(name string) (int64, error) {
	var count int64
	if tx := a.db.Model(&[]model.Article{}).Joins("join article_categories a_c on article.id = a_c.article_id").
		Joins("join category c on a_c.category_id = c.id").Count(&count); tx.Error != nil {
		return count, tx.Error
	}
	return count, nil
}

func (a *articleRepo) UpdateById(id int64, article *model.Article) (*model.Article, error) {
	findById, err := a.FindById(id)
	if err != nil {
		return findById, err
	}
	if tx := a.db.Updates(&article); tx.Error != nil {
		return article, tx.Error
	}
	return article, err
}

func (a *articleRepo) DeleteById(id int64) error {
	_, err := a.FindById(id)
	if err != nil {
		return err
	}
	a.db.Delete(&model.Article{}, id)
	return nil
}

func (a *articleRepo) Count() (int64, error) {
	var count int64
	a.db.Table("article").Select("id").Count(&count)
	return count, nil
}
