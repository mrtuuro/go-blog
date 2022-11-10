package models

import "time"

type Comment struct {
	ID      string    `json:"id,omitempty" bson:"_id,omitempty"`
	Owner   string    `json:"owner,omitempty" bson:"owner,omitempty"`
	Content string    `json:"content,omitempty" bson:"content,omitempty"`
	Date    time.Time `json:"date" bson:"date"`
}

type CreateComment struct {
	Owner   string    `json:"owner,omitempty" bson:"owner,omitempty"`
	Content string    `json:"content,omitempty" bson:"content,omitempty"`
	Date    time.Time `json:"date" bson:"date"`
}
