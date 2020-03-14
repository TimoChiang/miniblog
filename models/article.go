package models

import (
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql" // use _ to execute package's init function only
	"log"
	"os"
	"time"
)

const (
	USERNAME = "user"
	PASSWORD = "password"
	NETWORK = "tcp"
	// SERVER = "127.0.0.1"
	PORT = 3306
	DATABASE = "blog"
)

type Article struct {
	Id int `json:"id" form:"id"`
	Title string `json:"title" form:"title"`
	Description string `json:"description" form:"description"`
}

var db *sql.DB

func init() {
	conn := fmt.Sprintf("%s:%s@%s(%s:%d)/%s",USERNAME, PASSWORD, NETWORK, os.Getenv("MYSQL_HOST"), PORT, DATABASE)
	database, err := sql.Open("mysql", conn)
	if err != nil {
		log.Fatal("connection to mysql failed:", err)
	}
	database.SetConnMaxLifetime(100*time.Second)
	database.SetMaxOpenConns(100)
	db = database
}

func GetSingleArticle (id int) (*Article, error){
	article := new(Article)
	row := db.QueryRow("select id, title, description from articles where id=?", id)
	if err := row.Scan(&article.Id, &article.Title, &article.Description); err != nil {
		fmt.Printf("scan failed, err:%v\n", err)
		return article, err
	}
	fmt.Println("Single row data:", *article)
	return article, nil
}

func GetAllArticle() (map[int]*Article, error) {
	articles := make(map[int]*Article)
	rows, err := db.Query("select id, title, description from articles")
	if err != nil {
		panic(err.Error())
	}
	for rows.Next() {
		article := new(Article)
		if err := rows.Scan(&article.Id, &article.Title, &article.Description); err != nil {
			fmt.Printf("scan failed, err:%v\n", err)
			return articles, err
		}
		articles[article.Id] = article
	}

	fmt.Println("all row data:", articles)
	return articles, nil
}


func CreateArticle(title, description string) (int64, error) {
	result, err := db.Exec("insert INTO articles (title, description) values(?,?)", title, description)
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

