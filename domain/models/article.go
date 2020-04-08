package models

type Article struct {
	Id int `json:"id" form:"id"`
	Title string `json:"title" form:"title"   validate:"required,min=5,max=50"`
	Description string `json:"description" form:"description"  validate:"required"`
}