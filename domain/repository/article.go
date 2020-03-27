package repository

import (
	"database/sql"
	"fmt"
	"miniblog/domain/models"
)

type ArticleRepository struct {
	Db *sql.DB
}

func (r *ArticleRepository) GetArticle(id int) (*models.Article, error){
	article := new(models.Article)
	row := r.Db.QueryRow("select id, title, description from articles where id=?", id)
	if err := row.Scan(&article.Id, &article.Title, &article.Description); err != nil {
		if err == sql.ErrNoRows {
			fmt.Println("no row match conditions")
			return nil, nil
		}else{
			fmt.Printf("scan failed, err: %v\n", err)
			return nil, err
		}
	}
	fmt.Println("Single row data:", *article)
	return article, nil
}

func (r *ArticleRepository) GetAllArticle() (map[int]*models.Article, error) {
	articles := make(map[int]*models.Article)
	rows, err := r.Db.Query("select id, title, description from articles")
	if err != nil {
		panic(err.Error())
	}
	for rows.Next() {
		article := new(models.Article)
		if err := rows.Scan(&article.Id, &article.Title, &article.Description); err != nil {
			fmt.Printf("scan failed, err:%v\n", err)
			return articles, err
		}
		articles[article.Id] = article
	}

	fmt.Println("all row data:", articles)
	return articles, nil
}


func (r *ArticleRepository) CreateArticle(title, description string) (int64, error) {
	result, err := r.Db.Exec("insert INTO articles (title, description) values(?,?)", title, description)
	if err != nil {
		fmt.Printf("Insert failed,err:%v\n", err)
		return 0, err
	}
	lastInsertID, err := result.LastInsertId()
	if err != nil {
		fmt.Printf("insertData Get lastInsertID failed,err:%v\n", err)
		return 0, err
	}
	return lastInsertID, nil
}
