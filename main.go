package main

import (
	"database/sql"
	"fmt"
	//"github.com/go-playground/validator/v10"
	_ "github.com/go-sql-driver/mysql" // use _ to execute package's init function only
	"log"
	"miniblog/domain/repository"
	"miniblog/domain/service"
	"miniblog/domain/validator"
	c "miniblog/http/controller"
	"miniblog/routes"
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

func main() {
	//DB initialize
	conn := fmt.Sprintf("%s:%s@%s(%s:%d)/%s",USERNAME, PASSWORD, NETWORK, os.Getenv("MYSQL_HOST"), PORT, DATABASE)
	db, err := sql.Open("mysql", conn)
	if err != nil {
		log.Fatal("connection to mysql failed:", err)
	}
	db.SetConnMaxLifetime(100*time.Second)
	db.SetMaxOpenConns(100)

	// Route initialize
	router := routes.InitialRouters()

	// User
	userService := &service.UserService{}
	userHandler := &c.UserHandler{Service: userService}
	routes.SetUserRouters(router, userHandler)

	//Article
	articleRepository := repository.NewArticleRepository(db)
	articleService := &service.ArticleService{Repo: articleRepository, V: validator.NewValidator(articleRepository)}
	articleHandler := &c.ArticleHandler{Service: articleService, UserService:userService}
	routes.SetArticleRouters(router, articleHandler)

	routes.Serve("8888", router)
}
