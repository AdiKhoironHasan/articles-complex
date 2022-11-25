package repository

import (
	"fmt"
	"log"

	"github.com/AdiKhoironHasan/articles-complex/internal/models"
	"github.com/AdiKhoironHasan/articles-complex/internal/repository"
	servErrors "github.com/AdiKhoironHasan/articles-complex/pkg/errors"
	"github.com/jmoiron/sqlx"
)

const (
	CreateArticles   = `INSERT INTO public.articles (author, title, body, created) VALUES ($1, $2, $3, now()) RETURNING id`
	ShowArticles     = `SELECT id, author, title, body, created FROM public.articles WHERE %s ORDER BY created DESC`
	ShowArticlesByID = `SELECT id, author, title, body, created FROM public.articles WHERE id = $1 LIMIT 1`
	UpdateArticle    = `UPDATE public.articles SET author = $1, title = $2, body = $3 WHERE id = $4 RETURNING created`
	DeleteArticle    = `DELETE FROM public.articles WHERE id = $1`
)

var statement PreparedStatement

type PreparedStatement struct {
	createArticles   *sqlx.Stmt
	showArticlesByID *sqlx.Stmt
	updateArticle    *sqlx.Stmt
	deleteArticle    *sqlx.Stmt
}

type PostgreSqlRepo struct {
	Conn *sqlx.DB
}

func NewRepo(Conn *sqlx.DB) repository.SqlRepository {
	repo := &PostgreSqlRepo{Conn}
	InitPreparedStatement(repo)
	return repo
}

func (m *PostgreSqlRepo) Preparex(query string) *sqlx.Stmt {
	statement, err := m.Conn.Preparex(query)
	if err != nil {
		log.Fatalf("Failed to preparex query: %s. Error: %s", query, err.Error())
	}

	return statement
}

func InitPreparedStatement(m *PostgreSqlRepo) {
	statement = PreparedStatement{
		createArticles:   m.Preparex(CreateArticles),
		showArticlesByID: m.Preparex(ShowArticlesByID),
		updateArticle:    m.Preparex(UpdateArticle),
		deleteArticle:    m.Preparex(DeleteArticle),
	}
}

func (m *PostgreSqlRepo) CreateArticles(dataArticle *models.ArticleModels) (int, error) {
	var idArticle int
	err := m.Conn.QueryRow(CreateArticles, dataArticle.Author, dataArticle.Title, dataArticle.Body).Scan(&idArticle)

	if err != nil {
		log.Println("Failed Query CreateArticles: ", err.Error())
		return 0, fmt.Errorf(servErrors.ErrorDB)
	}

	return int(idArticle), nil
}

func (m *PostgreSqlRepo) ShowArticles(where string) ([]*models.ArticleModels, error) {
	var dataArticles []*models.ArticleModels

	var query string

	if where != "" && where != "%s" {
		query = fmt.Sprintf(ShowArticles, where)
	} else {
		query = fmt.Sprintf(ShowArticles, "1=1")
	}

	err := m.Conn.Select(&dataArticles, query)

	if err != nil {
		log.Println("Failed Query ShowArticles : ", err.Error())
		return nil, fmt.Errorf(servErrors.ErrorDB)
	}

	if len(dataArticles) == 0 {
		log.Println("Data Not Found ShowArticles")
		return nil, nil
	}

	return dataArticles, nil
}

func (m *PostgreSqlRepo) ShowArticlesByID(id int) ([]*models.ArticleModels, error) {
	var dataArticle []*models.ArticleModels

	err := m.Conn.Select(&dataArticle, ShowArticlesByID, id)

	if err != nil {
		log.Println("Failed Query ShowArticlesByID : ", err.Error())
		return nil, fmt.Errorf(servErrors.ErrorDB)
	}

	if len(dataArticle) == 0 {
		log.Println("Data Not Found ShowArticlesByID")
		return nil, servErrors.ErrNotFound
	}

	return dataArticle, nil
}

func (m *PostgreSqlRepo) UpdateArticle(dataArticle *models.ArticleModels) (string, error) {
	var created string
	err := m.Conn.QueryRow(UpdateArticle, dataArticle.Author, dataArticle.Title, dataArticle.Body, dataArticle.ID).Scan(&created)

	if err != nil {
		log.Println("Failed Query UpdateArticle : ", err.Error())
		return "", fmt.Errorf(servErrors.ErrorDB)
	}

	return created, nil
}

func (m *PostgreSqlRepo) DeleteArticle(id int) error {
	result, err := statement.deleteArticle.Exec(id)

	if err != nil {
		log.Println("Failed Query DeleteArticle : ", err.Error())
		return fmt.Errorf(servErrors.ErrorDB)
	}

	rows, err := result.RowsAffected()
	if err != nil {
		log.Println("Failed RowAffectd DeleteArticle : ", err.Error())
		return fmt.Errorf(servErrors.ErrorDB)
	}

	if rows < 1 {
		log.Println("DeleteArticle: No Data Deleted")
		return servErrors.ErrNotFound
	}

	return nil
}
