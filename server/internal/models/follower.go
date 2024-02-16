package models

type Follower struct {
	ID				int64	`json:"id"`
	User      int64	`json:"user"`
	Following int64	`json:"following"`
}
