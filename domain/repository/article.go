package repository

import (
	"database/sql"
	"fmt"
	"log"
	"miniblog/domain/models"
	"strings"
)

type ArticleRepository interface {
	GetArticle(id int) (*models.Article, error)
	GetAllArticle() (map[int]*models.Article, error)
	CreateArticle(articleStruct *models.Article) (int64, error)
	SlugExists(slug string) bool
}

type articleRepository struct {
	Db *sql.DB
}

func NewArticleRepository(db *sql.DB) ArticleRepository {
	return &articleRepository{db}
}

func (r *articleRepository) GetArticle(id int) (*models.Article, error){
	article := new(models.Article)
	row := r.Db.QueryRow("select id, title, description, slug from articles where id=?", id)
	if err := row.Scan(&article.Id, &article.Title, &article.Description, &article.Slug); err != nil {
		if err == sql.ErrNoRows {
			fmt.Println("no row match conditions")
			return nil, nil
		}else{
			fmt.Printf("scan failed, err: %v\n", err)
			return nil, err
		}
	}

	rowTags, err := r.Db.Query("select a.tag_id, t.name from articles_tags AS a join tags AS t on a.tag_id = t.id where article_id=?", id)
	if err != nil {
		panic(err.Error())
	}
	for rowTags.Next() {
		tag := new(models.Tag)
		if err := rowTags.Scan(&tag.Id, &tag.Name); err != nil {
			fmt.Printf("scan failed, err:%v\n", err)
			return article, err
		}
		article.Tags = append(article.Tags, tag)
	}

	fmt.Println("Single row data:", *article)
	return article, nil
}

func (r *articleRepository) GetAllArticle() (map[int]*models.Article, error) {
	articles := make(map[int]*models.Article)
	rows, err := r.Db.Query("select id, title, description, slug from articles")
	if err != nil {
		panic(err.Error())
	}

	for rows.Next() {
		article := new(models.Article)
		if err := rows.Scan(&article.Id, &article.Title, &article.Description, &article.Slug); err != nil {
			fmt.Printf("scan failed, err:%v\n", err)
			return articles, err
		}

		rowTags, err := r.Db.Query("select a.tag_id, t.name from articles_tags AS a join tags AS t on a.tag_id = t.id where article_id=?", article.Id)
		if err != nil {
			panic(err.Error())
		}
		for rowTags.Next() {
			tag := new(models.Tag)
			if err := rowTags.Scan(&tag.Id, &tag.Name); err != nil {
				fmt.Printf("scan failed, err:%v\n", err)
				return articles, err
			}
			article.Tags = append(article.Tags, tag)
		}

		articles[article.Id] = article
	}

	fmt.Println("all row data:", articles)
	return articles, nil
}


func (r *articleRepository) CreateArticle(articleStruct *models.Article) (int64, error) {
	result, err := r.Db.Exec("insert INTO articles (title, description, slug) values(?,?,?)", articleStruct.Title, articleStruct.Description, articleStruct.Slug)
	if err != nil {
		fmt.Printf("Insert failed,err:%v\n", err)
		return 0, err
	}
	lastInsertID, err := result.LastInsertId()
	if err != nil {
		fmt.Printf("insertData Get lastInsertID failed,err:%v\n", err)
		return 0, err
	}

	// bulk insert
	if len(articleStruct.Tags) > 0 {
		valueStrings := make([]string, 0, len(articleStruct.Tags))
		valueArgs := make([]interface{}, 0, len(articleStruct.Tags) * 1)
		for _, tag := range articleStruct.Tags {
			valueStrings = append(valueStrings, "(?)")
			valueArgs = append(valueArgs, tag.Name)
		}
		query := fmt.Sprintf("INSERT IGNORE INTO tags (name) VALUES %s",
			strings.Join(valueStrings, ","))
		_, err = r.Db.Exec(query, valueArgs...)
		if err != nil {
			fmt.Printf("Insert failed,err:%v\n", err)
			return 0, err
		}
		// get tags id
		selectQuery := fmt.Sprintf("select id from tags where tags.name in (?%s)",
			strings.Repeat(",?", len(valueStrings)-1))
		rowTags, err := r.Db.Query(selectQuery, valueArgs...)
		if err != nil {
			panic(err.Error())
		}

		for rowTags.Next() {
			tag := new(models.Tag)
			if err := rowTags.Scan(&tag.Id); err != nil {
				return 0, err
			}
			// insert to the relation table
			_, err := r.Db.Exec("insert ignore INTO articles_tags (article_id, tag_id) values(?,?)", lastInsertID, tag.Id)
			if err != nil {
				fmt.Printf("Insert failed,err:%v\n", err)
				return 0, err
			}
		}
	}

	return lastInsertID, nil
}

func (r *articleRepository) SlugExists(slug string) bool {
	if err := r.Db.QueryRow("select slug from articles where slug=?", slug).Scan(&slug); err != nil {
		if err != sql.ErrNoRows {
			log.Print(err)
		}
		return false
	}
	return true
}