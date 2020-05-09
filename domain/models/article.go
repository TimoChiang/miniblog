package models

type Article struct {
	Id int `json:"id" form:"id"`
	Title string `json:"title" form:"title"   validate:"required,min=5,max=50"`
	Slug string `json:"slug" form:"slug"   validate:"max=20,uniqueInDB"`
	Tags []*Tag
	Description string `json:"description" form:"description"  validate:"required"`
}

type Tag struct {
	Id int `json:"id" form:"id"`
	Name string `json:"name" form:"name"   validate:"max=20"`
}
