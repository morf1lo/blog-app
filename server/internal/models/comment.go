package models

import "github.com/go-playground/validator/v10"

type Comment struct {
	ID       int64   		 `json:"id"`
	Post     CommentPost `json:"post"`
	AuthorID int64  		 `json:"author_id" validate:"required"`
	Text     string 		 `json:"text" validate:"required"`
}

type CommentPost struct {
	ID		   int64 `json:"id"`
	AuthorID int64 `json:"author"`
}

func (c *Comment) Validate() error {
	validate := validator.New()
	return validate.Struct(c)
}
