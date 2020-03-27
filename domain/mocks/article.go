package mocks

import (
	"fmt"
	"miniblog/domain/models"
)

type ArticleRepository struct {
}

var mockArticles = map[int]*models.Article{
	1: &models.Article{Id: 1, Title: "test title", Description: "test description"},
	2: &models.Article{Id: 2, Title: "This this article2", Description: "Hello world"},
}

func (r *ArticleRepository) GetArticle(id int) (*models.Article, error){
	fmt.Println(mockArticles)
	if id > 0 && id <= len(mockArticles)  {
		return mockArticles[id], nil
	}else{
		return nil, nil
	}
}

func (r *ArticleRepository) GetAllArticle() (map[int]*models.Article, error) {
	return mockArticles, nil
}

func (r *ArticleRepository) CreateArticle(title, description string) (int64, error) {
	count := len(mockArticles)
	count++
	mockArticles[count] = &models.Article{Id: count, Title: title, Description: description}
	return int64(count), nil
}
