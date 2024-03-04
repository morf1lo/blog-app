package models

import "github.com/go-playground/validator/v10"

type Post struct {
	ID		         int64	`json:"id"`
	AuthorID       int64	`json:"author_id" validate:"required"`
	AuthorUsername string `json:"author_username"`
	Title	         string	`json:"title" validate:"min=1,max=50,required"`
	Text	         string	`json:"text" validate:"min=1,max=120,required"`
	Likes	         uint64	`json:"likes"`
}

func (p *Post) Validate() error {
	validate := validator.New()
	return validate.Struct(p)
}
