package assembler

import (
	"github.com/AdiKhoironHasan/articles-complex/internal/models"
	"github.com/AdiKhoironHasan/articles-complex/pkg/dto"
)

func ToCreateArticle(d *dto.ArticleReqDTO) *models.ArticleModels {
	return &models.ArticleModels{
		Author: d.Author,
		Title:  d.Title,
		Body:   d.Body,
	}
}

func ToUpdateArticle(d *dto.ArticleReqDTO) *models.ArticleModels {
	return &models.ArticleModels{
		ID:     d.ID,
		Author: d.Author,
		Title:  d.Title,
		Body:   d.Body,
	}
}

func ToCreateArticleResponse(d []*models.ArticleModels) []*dto.ArticleResDTO {
	var dataArticles []*dto.ArticleResDTO
	for _, val := range d {
		dataArticles = append(dataArticles, &dto.ArticleResDTO{
			ID:      val.ID,
			Title:   val.Title,
			Author:  val.Author,
			Body:    val.Body,
			Created: val.Created,
		})
	}

	return dataArticles
}

func ToShowArticlesResponse(d []*models.ArticleModels) []*dto.ArticleResDTO {
	var dataArticles []*dto.ArticleResDTO

	for _, val := range d {
		dataArticles = append(dataArticles, &dto.ArticleResDTO{
			ID:      val.ID,
			Title:   val.Title,
			Author:  val.Author,
			Body:    val.Body,
			Created: val.Created,
		})
	}

	return dataArticles
}

func ToShowArticlesByIDResponse(d []*models.ArticleModels) *dto.ArticleResDTO {
	var dataArticles *dto.ArticleResDTO
	dataArticles = &dto.ArticleResDTO{
		ID:      d[0].ID,
		Title:   d[0].Title,
		Author:  d[0].Author,
		Body:    d[0].Body,
		Created: d[0].Created,
	}

	return dataArticles
}
