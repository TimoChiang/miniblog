package models

type Article struct {
	Id int `json:"id" form:"id"`
	Title string `json:"title" form:"title"`
	Description string `json:"description" form:"description"`
}