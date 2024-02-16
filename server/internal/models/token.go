package models

type Token struct {
	ID				int64		`json:"id"`
	Username	string	`json:"username"`
	Avatar		string	`json:"avatar"`
}
