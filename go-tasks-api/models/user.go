package models

type User struct {
    ID          uint   `json:"id"`
    login       string `json:"login" validate:"required"`
}


