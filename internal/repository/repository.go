package repository

import (
	"github.com/AdiKhoironHasan/articles-complex/internal/models"
	"github.com/AdiKhoironHasan/articles-complex/pkg/dto"
)

// repository contract postgresql
type SqlRepository interface {
	CreateArticles(dataArticle *models.ArticleModels) (int, error)
	ShowArticles(where string) ([]*models.ArticleModels, error)
	ShowArticlesByID(id int) ([]*models.ArticleModels, error)
	UpdateArticle(dataArticle *models.ArticleModels) (string, error)
	DeleteArticle(id int) error
}

// repository contract redis
type NoSqlRepository interface {
	CreateArticles(id int, req *dto.ArticleResDTO) error
	CreateAllArticles(articles []*dto.ArticleResDTO) error
	ShowArticles() ([]*dto.ArticleResDTO, error)
	ShowArticlesByID(id int) (*dto.ArticleResDTO, error)
	DeleteArticle(id int)
}
