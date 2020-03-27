package service

import (
	"miniblog/domain/models"
	"miniblog/domain/repository"
)

type ArticleService struct {
	Repo repository.ArticleRepository
}

func (s *ArticleService) GetArticle(id int) (article *models.Article, err error){
	article, err = s.Repo.GetArticle(id)
	if err != nil {
		return nil, err
	}
	return
}


func (s *ArticleService) GetAllArticle() (articles map[int]*models.Article, err error){
	articles, err = s.Repo.GetAllArticle()
	if err != nil {
		return nil, err
	}
	return
}

func (s *ArticleService) CreateArticle(title, description string) (id int64, err error){
	id, err = s.Repo.CreateArticle(title, description)
	if err != nil {
		return 0, err
	}
	return
}
