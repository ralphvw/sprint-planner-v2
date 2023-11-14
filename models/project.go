package models

import "time"

type Project struct {
	ID          int       `json:"id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	CreatedAt   time.Time `json:"createdAt"`
	Members     []User    `json:"members"`
	MemberCount int       `json:"memberCount"`
	Owner       User      `json:"owner"`
}
