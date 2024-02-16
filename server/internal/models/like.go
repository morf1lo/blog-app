package models

type Like struct {
	ID	 int64 `json:"id"`
	User int64 `json:"user"`
	Post int64 `json:"post"`
}
