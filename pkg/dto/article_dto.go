package dto

import (
	"github.com/AdiKhoironHasan/articles-complex/pkg/common/validator"
)

// for ArticleReqDTO request json data
type ArticleReqDTO struct {
	ID      int
	Author  string `json:"author" valid:"required" validname:"author"`
	Title   string `json:"title" valid:"required" validname:"title"`
	Body    string `json:"body" valid:"required" validname:"body"`
	Created string
}

// for run validation in ArticleReqDTO
func (dto *ArticleReqDTO) Validate() error {
	v := validator.NewValidate(dto)

	return v.Validate()
}

// for ArticleParamReqDTO request json data
type ArticleParamReqDTO struct {
	Query  string `json:"query" validname:"query" query:"query"`
	Author string `json:"author" validname:"author" query:"author"`
}

// for run validation in ArticleParamReqDTO
func (dto *ArticleParamReqDTO) Validate() error {
	v := validator.NewValidate(dto)

	return v.Validate()
}
