package models

type User struct {
	ID    string `json:"id,omitempty" bson:"_id,omitempty"`
	Name  string `json:"name"`
	Email string `json:"email"`
}

type CreateUser struct {
	Name     string `validate:"required" json:"name"`
	Email    string `validate:"required" json:"email"`
	Password string `validate:"required" json:"password"`
}
