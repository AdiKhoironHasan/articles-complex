package dto

// for all response json data
type ResponseDTO struct {
	Success bool        `json:"success"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
}

// for ArticleResDTO response json data
type ArticleResDTO struct {
	ID      int    `json:"id"`
	Author  string `json:"author"`
	Title   string `json:"title"`
	Body    string `json:"body"`
	Created string `json:"created"`
}
