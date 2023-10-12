package services

import (
	"errors"

	"github.com/maxik12233/blog/types"
	"go.uber.org/zap"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type ResourceService interface {
	CreateArticle(art *types.Article) (uint, error)
	DeleteArticle(id uint) error
	UpdateArticle(art *types.Article) error // TODO
	GetArticles(amount uint, page uint) ([]*types.Article, int, error)
	GetOneArticle(id uint) (*types.Article, error)

	CreateComment(comm *types.Comment) (uint, error)
	DeleteComment(id uint) error
	UpdateComment(comm *types.Comment) error // TODO
	GetArticleComments(artid uint) ([]*types.Comment, error)

	ToggleLike(like *types.Like, flag bool) error
}

type ResourceServiceImpl struct {
	db     *gorm.DB
	logger *zap.Logger
}

func NewResourceService(db *gorm.DB, logger *zap.Logger) ResourceService {
	return &ResourceServiceImpl{
		db:     db,
		logger: logger,
	}
}

func (r *ResourceServiceImpl) CreateArticle(art *types.Article) (uint, error) {

	result := r.db.Create(&art)
	if result.Error != nil {
		return 0, result.Error
	}

	return art.ID, nil

}

func (r *ResourceServiceImpl) DeleteArticle(id uint) error {

	var art *types.Article

	result := r.db.Find(&art, id)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return errors.New("Nothing was deleted, probably wrong id")
	}

	result = r.db.Unscoped().Delete(&art)
	if result.Error != nil {
		return result.Error
	}

	return nil

}

func (r *ResourceServiceImpl) UpdateArticle(art *types.Article) error {
	return nil
}

func (r *ResourceServiceImpl) GetOneArticle(id uint) (*types.Article, error) {

	var art *types.Article
	result := r.db.Where("ID = ?", id).Preload(clause.Associations).Find(&art)
	if result.Error != nil {
		return nil, result.Error
	}
	if result.RowsAffected == 0 {
		return nil, errors.New("Article not found")
	}

	return art, nil
}

func (r *ResourceServiceImpl) GetArticles(amount uint, page uint) ([]*types.Article, int, error) {

	var arts []*types.Article

	result := r.db.Preload(clause.Associations).Limit(int(amount)).Offset(int(amount * page)).Find(&arts)
	if result.Error != nil {
		return nil, 0, result.Error
	}
	if result.RowsAffected == 0 {
		return nil, 0, errors.New("Articles not found")
	}

	var count int64
	r.db.Model(&types.Article{}).Count(&count)

	return arts, int(count), nil
}

func (r *ResourceServiceImpl) CreateComment(comm *types.Comment) (uint, error) {

	result := r.db.Create(&comm)
	if result.Error != nil {
		return 0, result.Error
	}

	return comm.ID, nil

}

func (r *ResourceServiceImpl) DeleteComment(id uint) error {

	var comm *types.Comment

	result := r.db.Find(&comm, id)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return errors.New("Nothing was deleted, probably wrong id")
	}

	result = r.db.Unscoped().Delete(&comm)
	if result.Error != nil {
		return result.Error
	}

	return nil

}

func (r *ResourceServiceImpl) UpdateComment(comm *types.Comment) error {
	return nil
}

func (r *ResourceServiceImpl) GetArticleComments(artid uint) ([]*types.Comment, error) {

	var comms []*types.Comment
	result := r.db.Where("article_id = ?", artid).Preload("Like").Find(&comms)
	if result.Error != nil {
		return nil, result.Error
	}

	return comms, nil
}

func (r *ResourceServiceImpl) ToggleLike(like *types.Like, flag bool) error {

	if flag {
		var count int64
		r.db.Model(&types.Like{}).Where("user_id = ? AND (article_id = ? OR comment_id = ?)", like.UserID, like.ArticleID, like.CommentID).Count(&count)
		if count > 0 {
			return errors.New("This user already like that resource")
		}
		result := r.db.Create(&like)
		if result.Error != nil {
			return result.Error
		}

		return nil
	}

	result := r.db.Delete(&types.Like{}, "user_id = ? AND (article_id = ? OR article_id IS NULL) AND (comment_id = ? OR comment_id IS NULL)", like.UserID, like.ArticleID, like.CommentID)
	if result.Error != nil {
		return result.Error
	}

	return nil
}
