package repository

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/AdiKhoironHasan/articles-complex/internal/repository"
	"github.com/AdiKhoironHasan/articles-complex/pkg/dto"
	servErrors "github.com/AdiKhoironHasan/articles-complex/pkg/errors"
	"github.com/go-redis/redis/v8"
)

const ()

type RedisNoSqlRepo struct {
	Conn *redis.Client
}

func NewRepo(Conn *redis.Client) repository.NoSqlRepository {

	repo := &RedisNoSqlRepo{Conn}
	return repo
}

// implementation repository redis
func (p *RedisNoSqlRepo) CreateArticles(id int, req *dto.ArticleResDTO) error {
	key := fmt.Sprintf("article_%d", id)
	data := dto.ArticleResDTO{
		ID:      id,
		Author:  req.Author,
		Title:   req.Title,
		Body:    req.Body,
		Created: req.Created,
	}

	// encode data string to json
	value, err := json.Marshal(data)
	if err != nil {
		log.Println("Failed Marshal CreateArticles:", err.Error())
		return err
	}

	// expired time data in redis (10 minute)
	ttl := 10 * time.Minute

	// create data in redis
	redis := p.Conn.Set(context.Background(), key, value, ttl)

	result, err := redis.Result()
	if err != nil {
		log.Println("Failed Redis Get Result CreateArticles:", err.Error())
		return fmt.Errorf(servErrors.ErrorDB)
	}

	if result != "OK" {
		log.Println("Failed Redis Not OK CreateArticles: ", result)
		return fmt.Errorf(result)
	}

	return nil
}

func (p *RedisNoSqlRepo) CreateAllArticles(articles []*dto.ArticleResDTO) error {
	key := "article_all"

	// encode data string to json
	value, err := json.Marshal(articles)
	if err != nil {
		log.Println("Failed Marshal CreateAllArticles:", err.Error())
		return err
	}

	// expired time data in redis (10 minute)
	ttl := 10 * time.Minute

	// create data in redis
	redis := p.Conn.Set(context.Background(), key, value, ttl)

	result, err := redis.Result()
	if err != nil {
		log.Println("Failed Redis Get Result CreateAllArticles:", err.Error())
		return fmt.Errorf(servErrors.ErrorDB)
	}

	if result != "OK" {
		log.Println("Failed Redis Not OK CreateAllArticles: ", result)
		return fmt.Errorf(result)
	}

	return nil
}

func (p *RedisNoSqlRepo) ShowArticles() ([]*dto.ArticleResDTO, error) {
	var dataArticles []*dto.ArticleResDTO
	key := "article_all"

	// get data in redis
	result := p.Conn.Get(context.Background(), key).Val()
	if result == "" {
		log.Println("Redis: No Result for key 'article_all' ShowArticles")
		return nil, nil
	}

	jsonData := []byte(result)

	// decode data json to string
	err := json.Unmarshal(jsonData, &dataArticles)
	if err != nil {
		log.Println("Failed UnMarshal ShowArticles:", err.Error())
		return nil, err
	}

	return dataArticles, nil
}

func (p *RedisNoSqlRepo) ShowArticlesByID(id int) (*dto.ArticleResDTO, error) {
	var dataArticles *dto.ArticleResDTO
	key := fmt.Sprintf("article_%d", id)

	// get data in redis
	result := p.Conn.Get(context.Background(), key).Val()
	if result == "" {
		log.Println(fmt.Sprintf("Redis: No Result for key '%s' ShowArticlesByID", key))
		return nil, nil
	}

	jsonData := []byte(result)

	// decode data json to string
	err := json.Unmarshal(jsonData, &dataArticles)
	if err != nil {
		log.Println("Failed UnMarshal ShowArticlesByID:", err.Error())
		return nil, err
	}

	return dataArticles, nil
}

func (p *RedisNoSqlRepo) DeleteArticle(id int) {
	key := fmt.Sprintf("article_%d", id)

	// delete data in redis
	result := p.Conn.Del(context.Background(), key).Val()
	// result < 1 == fail
	if result < 1 {
		log.Println(fmt.Sprintf("Redis: Failed detele for key '%s' DeleteArticle", key))
	}

}
