package service

import (
	"miniblog/domain/models"
	"miniblog/domain/repository"
	"miniblog/domain/validator"
	"net/http"
)

type ArticleService struct {
	Repo repository.ArticleRepository
	V validator.Validator
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

func (s *ArticleService) CreateArticle(articleStruct *models.Article) (id int64, err error){
	id, err = s.Repo.CreateArticle(articleStruct)
	if err != nil {
		return 0, err
	}
	return
}

// convert form data to struct, and then validate value
func (s *ArticleService) LoadArticleStruct(r *http.Request) (articleStruct *models.Article, err error) {
	if err := r.ParseForm(); err != nil {
		return articleStruct, err
	}
	articleStruct = new(models.Article)
	articleStruct.Title = r.FormValue("title")
	articleStruct.Description = r.FormValue("description")
	return articleStruct, nil
}
