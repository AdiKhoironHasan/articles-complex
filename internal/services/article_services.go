package services

import (
	"log"

	integ "github.com/AdiKhoironHasan/articles-complex/internal/integration"
	"github.com/AdiKhoironHasan/articles-complex/internal/repository"
	"github.com/AdiKhoironHasan/articles-complex/pkg/dto"
	"github.com/AdiKhoironHasan/articles-complex/pkg/dto/assembler"
)

type service struct {
	sqlRepo   repository.SqlRepository
	noSqlRepo repository.NoSqlRepository
	IntegServ integ.IntegServices
}

func NewService(sqlRepo repository.SqlRepository, noSqlRepo repository.NoSqlRepository, IntegServ integ.IntegServices) Services {
	return &service{sqlRepo, noSqlRepo, IntegServ}
}

func (s *service) CreateArticles(req *dto.ArticleReqDTO) error {
	dataArticleModel := assembler.ToCreateArticle(req)

	_, err := s.sqlRepo.CreateArticles(dataArticleModel)

	if err != nil {
		return err
	}

	return nil
}

// implementation service
func (s *service) ShowArticles(req *dto.ArticleParamReqDTO) ([]*dto.ArticleResDTO, error) {
	var dataArticles []*dto.ArticleResDTO
	var where string

	// get data from redis
	if req.Query == "" && req.Author == "" {
		dataArticles, err := s.noSqlRepo.ShowArticles()
		if err != nil {
			log.Println("Redis :", err)
		}

		if len(dataArticles) > 0 {
			log.Println("ShowArticles use Redis")
			return dataArticles, nil
		}
	}

	if req.Query != "" && req.Author != "" {
		where = "title LIKE '%" + req.Query + "%' AND body LIKE '%" + req.Query + "%' AND author LIKE '%" + req.Author + "%'"
	} else if req.Query != "" {
		where = "title LIKE '%" + req.Query + "%' AND body LIKE '%" + req.Query + "%'"
	} else if req.Author != "" {
		where = "author LIKE '%" + req.Author + "$'"
	}

	log.Println("ShowArticles use PostgreSQL")
	// show all articles from postgresql
	dataArticlesModels, err := s.sqlRepo.ShowArticles(where)
	if err != nil {
		return nil, err
	}

	// for convert data models to json
	dataArticles = assembler.ToShowArticlesResponse(dataArticlesModels)

	// create data article_all to redis
	err = s.noSqlRepo.CreateAllArticles(dataArticles)
	if err != nil {
		log.Println("Redis :", err)
	}

	return dataArticles, nil
}

func (s *service) ShowArticlesByID(id int) (*dto.ArticleResDTO, error) {
	var dataArticle *dto.ArticleResDTO

	// get data from redis
	dataArticle, err := s.noSqlRepo.ShowArticlesByID(id)
	if err != nil {
		log.Println(err)
	}

	if dataArticle != nil {
		log.Println("ShowArticlesByID use Redis")
		return dataArticle, nil
	}

	// show articles from postgresql by article id
	log.Println("ShowArticlesByID use PostgreSQL")
	dataArticleModels, err := s.sqlRepo.ShowArticlesByID(id)

	if err != nil {
		return nil, err
	}

	if dataArticleModels == nil {
		return nil, nil
	}

	// for convert data models to json
	dataArticle = assembler.ToShowArticlesByIDResponse(dataArticleModels)

	// for create data article_id to redis
	err = s.noSqlRepo.CreateArticles(id, dataArticle)
	if err != nil {
		log.Println("Redis :", err)
	}

	return dataArticle, nil
}

func (s *service) UpdateArticle(req *dto.ArticleReqDTO) error {
	// show articles from postgresql by article id
	dataArticleModels, err := s.sqlRepo.ShowArticlesByID(req.ID)

	if err != nil {
		return err
	}

	// if data nil == no data found, but success
	if dataArticleModels == nil {
		return nil
	}

	// for convert data json to model
	dataArticle := assembler.ToUpdateArticle(req)

	// update articles in postgresql by article id
	_, err = s.sqlRepo.UpdateArticle(dataArticle)
	if err != nil {
		return err
	}

	return nil
}

func (s *service) DeleteArticle(id int) error {
	// show articles from postgresql by article id
	_, err := s.sqlRepo.ShowArticlesByID(id)

	if err != nil {
		return err
	}

	// delete articles in postgresql by article id
	err = s.sqlRepo.DeleteArticle(id)
	if err != nil {
		return err
	}

	return nil
}
