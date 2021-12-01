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
	FindAllArticleByWord(search string, limit, offset int) (*[]model.Article, error)
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
	return article, nil
}

func (a *articleRepo) FindAll(limit, offset int) (*[]model.Article, error) {
	var articles *[]model.Article
	if tx := a.db.Preload("User").Preload("Categories").Offset(offset).Limit(limit).Find(&articles); tx.Error != nil {
		return articles, tx.Error
	}
	return articles, nil
}

func (a *articleRepo) FindById(id int64) (*model.Article, error) {
	var article *model.Article
	if tx := a.db.Preload("User").Preload("Categories").First(&article, id); tx.Error != nil {
		return article, tx.Error
	}
	return article, nil
}

func (a *articleRepo) FindByTitle(title string) (*model.Article, error) {
	var article *model.Article
	if tx := a.db.Preload("User").Preload("Categories").First(&article, "article.title = ?", title); tx.Error != nil {
		return article, tx.Error
	}
	return article, nil
}

func (a *articleRepo) FindAllArticleByWord(search string, limit, offset int) (*[]model.Article, error) {
	var articles *[]model.Article
	if tx := a.db.Preload("User").Preload("Categories").Joins("join article_categories a_c on article.id = a_c.article_id").
		Joins("join category c on a_c.category_id = c.id").Where("article.description LIKE ? or article.title LIKE ? or c.name LIKE ?", "%"+search+"%", "%"+search+"%", "%"+search+"%").Offset(offset).Limit(limit).Group("article.id").Find(&articles); tx.Error != nil {
		return articles, tx.Error
	}
	return articles, nil
}

func (a *articleRepo) FindAllByUserId(uuid uuid.UUID, limit, offset int) (*[]model.Article, error) {
	var articles *[]model.Article
	if tx := a.db.Preload("User").Preload("Categories").Offset(offset).Limit(limit).Find(&articles, "article.user_id = ?", uuid); tx.Error != nil {
		return articles, tx.Error
	}
	return articles, nil
}

func (a *articleRepo) FindAllByUsername(username string, limit, offset int) (*[]model.Article, error) {
	var articles *[]model.Article
	if tx := a.db.Preload("User").Preload("Categories").Offset(offset).Limit(limit).Find(&articles, "u.username = ?", username); tx.Error != nil {
		return articles, tx.Error
	}
	return articles, nil
}

func (a *articleRepo) FindAllByCategoryName(name string, limit, offset int) (*[]model.Article, error) {
	var articles *[]model.Article
	if tx := a.db.Preload("Categories").Preload("User").Joins("join article_categories a_c on article.id = a_c.article_id").
		Joins("join category c on a_c.category_id = c.id").Limit(limit).Offset(offset).Find(&articles, "c.name = ?", &name); tx.Error != nil {
		return articles, tx.Error
	}
	return articles, nil
}

func (a *articleRepo) CountByCategoryName(name string) (int64, error) {
	var count int64
	if tx := a.db.Model(&[]model.Article{}).Joins("join article_categories a_c on article.id = a_c.article_id").
		Joins("join category c on a_c.category_id = c.id").Where("c.name = ?", name).Select("count(article.id)").Count(&count); tx.Error != nil {
		return count, tx.Error
	}
	return count, nil
}

func (a *articleRepo) UpdateById(id int64, article *model.Article) (*model.Article, error) {
	article.ID = &id
	a.db.Exec("DELETE FROM article_categories WHERE article_categories.article_id = ?", id)
	a.db.Updates(&article).Association("Categories")
	return article, nil
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
	a.db.Table("article").Select("id").Where("deleted_at is null").Count(&count)
	return count, nil
}
