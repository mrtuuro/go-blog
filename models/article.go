package models

type Article struct {
	ID       string          `json:"id,omitempty" bson:"_id,omitempty"`
	Name     string          `json:"name"`
	Rating   float64         `json:"rating"`
	Author   string          `json:"author"`
	Content  string          `json:"content"`
	Comments []CreateComment `json:"comments,omitempty" bson:"comments,omitempty"`
}
