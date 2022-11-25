package http

import (
	"net/http"
	"os"
	"strconv"

	"github.com/AdiKhoironHasan/articles-complex/internal/services"
	servConst "github.com/AdiKhoironHasan/articles-complex/pkg/common/const"
	"github.com/AdiKhoironHasan/articles-complex/pkg/dto"
	servErrors "github.com/AdiKhoironHasan/articles-complex/pkg/errors"
	"github.com/spf13/viper"

	"github.com/apex/log"
	"github.com/labstack/echo"
	middleware "github.com/labstack/echo/middleware"
)

type HttpHandler struct {
	service services.Services
}

func NewHttpHandler(e *echo.Echo, srv services.Services) {
	handler := &HttpHandler{
		srv,
	}

	// secret_key jwt from config-dev.yaml
	secret_key := viper.GetString("api.secret_key")

	// use prefix "api/v1" for all endpoints
	eJWT := e.Group("api/v1")
	// use jwt auth for all endpoints
	eJWT.Use(middleware.JWT([]byte(secret_key)))

	eJWT.GET("/ping", handler.Ping)
	eJWT.POST("/articles", handler.CreateArticles)
	eJWT.GET("/articles", handler.ShowArticles)
	eJWT.GET("/articles/:id", handler.ShowArticlesByID)
	eJWT.PUT("/articles/:id", handler.UpdateArticle)
	eJWT.DELETE("/articles/:id", handler.DeleteArticle)
}

func (h *HttpHandler) Ping(c echo.Context) error {
	// get version app from environment
	version := os.Getenv("VERSION")

	if version == "" {
		version = "pong"
	}

	data := version

	return c.JSON(http.StatusOK, data)

}

func (h *HttpHandler) CreateArticles(c echo.Context) error {
	postDTO := dto.ArticleReqDTO{}

	// fill variable from request data
	if err := c.Bind(&postDTO); err != nil {
		log.Error(err.Error())
		return c.NoContent(http.StatusBadRequest)
	}
	// data validation from validate in dto.ArticleReqDTO{}
	err := postDTO.Validate()
	if err != nil {
		log.Error(err.Error())
		return c.JSON(getStatusCode(err), dto.ResponseDTO{
			Success: false,
			Message: err.Error(),
			Data:    nil,
		})
	}

	// to running service
	err = h.service.CreateArticles(&postDTO)
	if err != nil {
		log.Error(err.Error())
		return c.JSON(getStatusCode(err), dto.ResponseDTO{
			Success: false,
			Message: err.Error(),
			Data:    nil,
		})
	}

	var resp = dto.ResponseDTO{
		Success: true,
		Message: servConst.SaveSuccess,
		Data:    nil,
	}

	return c.JSON(http.StatusOK, resp)
}

func (h *HttpHandler) ShowArticles(c echo.Context) error {
	getDTO := dto.ArticleParamReqDTO{}

	// fill variable from request data
	if err := c.Bind(&getDTO); err != nil {
		log.Error(err.Error())
		return c.NoContent(http.StatusBadRequest)
	}

	// data validation from validate in dto.ArticleParamReqDTO{}
	err := getDTO.Validate()
	if err != nil {
		log.Error(err.Error())
		return c.JSON(getStatusCode(err), dto.ResponseDTO{
			Success: false,
			Message: err.Error(),
			Data:    nil,
		})
	}

	// to running service
	result, err := h.service.ShowArticles(&getDTO)
	if err != nil {
		log.Error(err.Error())
		return c.JSON(getStatusCode(err), dto.ResponseDTO{
			Success: false,
			Message: err.Error(),
			Data:    nil,
		})
	}

	var resp = dto.ResponseDTO{
		Success: true,
		Message: servConst.GetDataSuccess,
		Data:    result,
	}

	return c.JSON(http.StatusOK, resp)
}

func (h *HttpHandler) ShowArticlesByID(c echo.Context) error {
	var articleID int

	// get article id from parameter
	articleID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		log.Error(err.Error())
		return c.NoContent(http.StatusBadRequest)
	}

	// to running service
	result, err := h.service.ShowArticlesByID(articleID)
	if err != nil {
		log.Error(err.Error())
		return c.JSON(getStatusCode(err), dto.ResponseDTO{
			Success: false,
			Message: err.Error(),
			Data:    nil,
		})
	}

	var resp = dto.ResponseDTO{
		Success: true,
		Message: servConst.GetDataSuccess,
		Data:    result,
	}

	return c.JSON(http.StatusOK, resp)
}

func (h *HttpHandler) UpdateArticle(c echo.Context) error {
	var articleID int
	putDTO := dto.ArticleReqDTO{}

	// get article id from parameter
	articleID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		log.Error(err.Error())
		return c.NoContent(http.StatusBadRequest)
	}

	// fill variable from request data
	if err := c.Bind(&putDTO); err != nil {
		log.Error(err.Error())
		return c.NoContent(http.StatusBadRequest)
	}

	putDTO.ID = articleID
	// data validation from validate in dto.ArticleReqDTO{}
	err = putDTO.Validate()
	if err != nil {
		log.Error(err.Error())
		return c.JSON(getStatusCode(err), dto.ResponseDTO{
			Success: false,
			Message: err.Error(),
			Data:    nil,
		})
	}

	// to running service
	err = h.service.UpdateArticle(&putDTO)
	if err != nil {
		log.Error(err.Error())
		return c.JSON(getStatusCode(err), dto.ResponseDTO{
			Success: false,
			Message: err.Error(),
			Data:    nil,
		})
	}

	var resp = dto.ResponseDTO{
		Success: true,
		Message: servConst.UpdateSuccess,
		Data:    nil,
	}

	return c.JSON(http.StatusOK, resp)

}

func (h *HttpHandler) DeleteArticle(c echo.Context) error {
	var articleID int

	// get article id from parameter
	articleID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		log.Error(err.Error())
		return c.NoContent(http.StatusBadRequest)
	}

	// to running service
	err = h.service.DeleteArticle(articleID)
	if err != nil {
		log.Error(err.Error())
		return c.JSON(getStatusCode(err), dto.ResponseDTO{
			Success: false,
			Message: err.Error(),
			Data:    nil,
		})
	}

	var resp = dto.ResponseDTO{
		Success: true,
		Message: servConst.DeleteSuccess,
		Data:    nil,
	}

	return c.JSON(http.StatusOK, resp)

}

func getStatusCode(err error) int {
	if err == nil {
		return http.StatusOK
	}

	switch err {
	case servErrors.ErrInternalServerError:
		return http.StatusInternalServerError
	case servErrors.ErrNotFound:
		return http.StatusNotFound
	case servErrors.ErrConflict:
		return http.StatusConflict
	case servErrors.ErrInvalidRequest:
		return http.StatusBadRequest
	case servErrors.ErrFailAuth:
		return http.StatusForbidden
	default:
		return http.StatusInternalServerError
	}
}
