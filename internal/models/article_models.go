package models

// for model in articles table
type ArticleModels struct {
	ID      int    `db:"id"`
	Author  string `db:"author"`
	Title   string `db:"title"`
	Body    string `db:"body"`
	Created string `db:"created"`
}
