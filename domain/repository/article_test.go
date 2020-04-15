package repository

import (
	"miniblog/domain/models"
	"testing"
)

var expectedTags = []*models.Tag{
	&models.Tag{Id: 1, Name: "test tag"},
	&models.Tag{Id: 2, Name: "test tag##2"},
}

var expectedArticles = map[int]*models.Article{
		1: &models.Article{Id: 1, Title: "test title", Slug: "test Slug", Tags: []*models.Tag{expectedTags[0]}, Description: "test description"},
		2: &models.Article{Id: 2, Title: "This this article2", Description: "Hello world"},
}

func prepareTestForArticle() {
	for _, tag := range expectedTags {
		articleRepo.Db.Exec("INSERT INTO tags (name) VALUES(?)", tag.Name)
	}
	for _, article := range expectedArticles {
		articleRepo.Db.Exec("INSERT INTO articles (title, description, slug) VALUES(?,?,?)", article.Title, article.Description, article.Slug)
		for _, tag := range article.Tags {
			articleRepo.Db.Exec("INSERT INTO articles_tags (article_id, tag_id) VALUES(?,?)", article.Id, tag.Id)
		}
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
		expected := expectedArticles[1]
		if  article == nil {
			t.Fatalf("expected %v, got nil", *expected)
		}
		isArticleEqualExpected(t, article, expected)
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
	for index, expected := range expectedArticles {
		isArticleEqualExpected(t, articles[index], expected)
	}
}

func TestArticleRepository_CreateArticle(t *testing.T) {
	var tests = []struct {
		in  string
		expected *models.Article
	}{
		{"expected", &models.Article{Id: 3, Title: "test create title", Description: "test create description", Slug: "mySlug", Tags: expectedTags}},
		{"expectedWithRequiredFields", &models.Article{Id: 4, Title: "test create title", Description: "test create description", Slug: "", Tags: []*models.Tag{}}},
	}
	for _, tt := range tests {
		t.Run(tt.in, func(t *testing.T) {
			lastInsertID, err := articleRepo.CreateArticle(tt.expected)
			if err != nil {
				t.Errorf("got error: %v", err)
			}
			article, err := articleRepo.GetArticle(int(lastInsertID))
			if err != nil {
				t.Errorf("got error: %v", err)
			}
			if  article == nil {
				t.Fatalf("expected %v, got nil", *tt.expected)
			}
			isArticleEqualExpected(t, article, tt.expected)
		})
	}
}

func isArticleEqualExpected (t *testing.T, article, expected *models.Article) {
	if article.Title != expected.Title {
		t.Errorf("expected %v, got: %v", expected.Title, article.Title)
	}

	if article.Description != expected.Description {
		t.Errorf("expected %v, got: %v", expected.Description, article.Description)
	}

	if article.Slug != expected.Slug {
		t.Errorf("expected %v, got: %v", expected.Slug, article.Slug)
	}

	if len(article.Tags) != len(expected.Tags) {
		t.Errorf("expected length %v, got: %v", len(expected.Tags), len(article.Tags))
	}

	for index, tag := range expected.Tags {
		if *article.Tags[index] != *tag {
				t.Errorf("expected %v, got: %v", *tag, *article.Tags[index])
		}
	}
}
