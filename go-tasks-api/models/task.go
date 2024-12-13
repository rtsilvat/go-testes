package models

type Task struct {
    ID          uint   `json:"id" gorm:"primaryKey"`
    Title       string `json:"title" validate:"required"`
    Description string `json:"description"`
    Completed   bool   `json:"completed"`
}


