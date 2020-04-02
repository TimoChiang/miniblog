package service

import (
	"miniblog/domain/mocks"
	"miniblog/domain/models"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"
)

func TestArticleService_GetArticle(t *testing.T) {
	articleRepository := new(mocks.ArticleRepository)
	articleService := &ArticleService{Repo: articleRepository}

	t.Run("no result", func(t *testing.T) {
		article, err := articleService.GetArticle(0)
		if err != nil {
			t.Errorf("got error: %v", err)
		}
		if article != nil {
			t.Errorf("expected nil, got: %v", article)
		}
	})

	// test result
	t.Run("got article", func(t *testing.T) {
		article, err := articleService.GetArticle(1)
		if  err != nil {
			t.Errorf("got error: %v", err)
		}
		if  article == nil {
			t.Errorf("expected return article, got: %v", article)
		}
	})
}

func TestArticleService_GetAllArticle(t *testing.T) {
	articleRepository := new(mocks.ArticleRepository)
	articleService := &ArticleService{Repo: articleRepository}
	articles, err := articleService.GetAllArticle()
	if err != nil {
		t.Errorf("got error: %v", err)
	}
	if len(articles) != 2 {
		t.Errorf("expected total article number is %q, got : %q", 2, len(articles))
	}
}

func TestArticleService_CreateArticle(t *testing.T) {
	articleRepository := new(mocks.ArticleRepository)
	articleService := &ArticleService{Repo: articleRepository}
	expected := &models.Article{Id: 3, Title: "test create title", Description: "test create description"}
	lastInsertID, err := articleService.CreateArticle(expected)
	if err != nil {
		t.Errorf("got error: %v", err)
	}
	article, err := articleService.GetArticle(int(lastInsertID))
	if err != nil {
		t.Errorf("got error: %v", err)
	}
	if article == nil {
		t.Fatalf("expected %v, got nil", *expected)
	}
	if *article != *expected {
		t.Errorf("expected %v, got: %v", *expected, *article)
	}
}

func TestArticleService_LoadArticleStruct(t *testing.T) {
	articleRepository := new(mocks.ArticleRepository)
	articleService := &ArticleService{Repo: articleRepository}

	expected := &models.Article{Title: "test load title", Description: "test load description"}
	body := url.Values{"title":{expected.Title}, "description":{expected.Description}}
	r := httptest.NewRequest("POST", "/articles", strings.NewReader(body.Encode()))
	r.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	articleStruct, err := articleService.LoadArticleStruct(r)
	if err != nil {
		t.Errorf("got error: %v", err)
	}
	if *articleStruct != *expected {
		t.Errorf("expected %v, got: %v", *expected, *articleStruct)
	}

}
