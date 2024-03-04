package models

import "strings"

type PostUpdateOptions struct {
	Title string `json:"title"`
	Text  string `json:"text"`
}

func (u *PostUpdateOptions) FilterUpdateOptions() (string, []interface{}) {
	query := "UPDATE posts SET"
	var values []interface{}

	if strings.TrimSpace(u.Title) == "" && strings.TrimSpace(u.Text) == "" {
		return "", nil
	}

	if strings.TrimSpace(u.Title) != "" {
		query += " title = ?,"
		values = append(values, strings.TrimSpace(u.Title))
	}

	if strings.TrimSpace(u.Text) != "" {
		query += " text = ?,"
		values = append(values, strings.TrimSpace(u.Text))
	}

	query = strings.TrimSuffix(query, ",")

	return query, values
}
