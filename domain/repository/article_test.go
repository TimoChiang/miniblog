package repository

import (
	"miniblog/domain/models"
	"testing"
)

var expectedArticles = []*models.Article{
		&models.Article{Id: 1, Title: "test title", Description: "test description"},
		&models.Article{Id: 2, Title: "This this article2", Description: "Hello world"},
}

func prepareTestForArticle() {
	for _, article := range expectedArticles {
		articleRepo.Db.Exec("INSERT INTO articles (title, description) VALUES(?,?)", article.Title, article.Description)
	}
}

func TestArticleRepository_GetArticle(t *testing.T) {
	// test no result
	t.Run("no result", func(t *testing.T) {
		article, err := articleRepo.GetArticle(0)
		if  err != nil {
			t.Errorf("got error: %v", err)
		}
		if  article != nil {
			t.Errorf("expected nil, got: %v", article)
		}
	})

	// test result
	t.Run("got article", func(t *testing.T) {
		article, err := articleRepo.GetArticle(1)
		if  err != nil {
			t.Errorf("got error: %v", err)
		}
		expected := expectedArticles[0]
		if  article == nil {
			t.Fatalf("expected %v, got nil", *expected)
		}
		if  *article != *expected {
			t.Errorf("expected %v, got: %v", *expected, *article)
		}
	})
}

func TestArticleRepository_GetAllArticle(t *testing.T) {
	articles, err := articleRepo.GetAllArticle()
	if err != nil {
		t.Errorf("got error: %v", err)
	}
	if len(articles) != len(expectedArticles) {
		t.Errorf("expected total article number is %q, got : %q", len(expectedArticles), len(articles))
	}
}

func TestArticleRepository_CreateArticle(t *testing.T) {
	expected := &models.Article{Id: 3, Title: "test create title", Description: "test create description"}
	lastInsertID, err := articleRepo.CreateArticle(expected)
	if err != nil {
		t.Errorf("got error: %v", err)
	}
	article, err := articleRepo.GetArticle(int(lastInsertID))
	if err != nil {
		t.Errorf("got error: %v", err)
	}
	if  article == nil {
		t.Fatalf("expected %v, got nil", *expected)
	}
	if  *article != *expected {
		t.Errorf("expected %v, got: %v", *expected, *article)
	}
}
