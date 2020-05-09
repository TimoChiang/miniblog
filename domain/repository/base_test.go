package repository

import (
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql" // use _ to execute package's init function only
	"io/ioutil"
	"log"
	"os"
	"strings"
	"testing"
	"time"
)

const (
	USERNAME = "user"
	PASSWORD = "password"
	NETWORK = "tcp"
	// SERVER = "127.0.0.1"
	PORT = 3307
	DATABASE = "blog"
)

var articleRepo *articleRepository

func setup() {
	//DB initialize
	println("-------setup-------")
	conn := fmt.Sprintf("%s:%s@%s(%s:%d)/%s",USERNAME, PASSWORD, NETWORK, os.Getenv("MYSQL_HOST"), PORT, DATABASE)
	db, err := sql.Open("mysql", conn)
	if err != nil {
		log.Fatal("connection to mysql failed:", err)
	}

	db.SetConnMaxLifetime(100*time.Second)
	db.SetMaxOpenConns(100)
	articleRepo = new(articleRepository)
	articleRepo.Db = db
	initialTables()
	prepareTestForArticle()
}

func teardown() {
	println("-------teardown-------")
	clearTable()
}

func TestMain(m *testing.M) {
	setup()
	code := m.Run()
	teardown()
	os.Exit(code)
}

func initialTables() {
	file, err := ioutil.ReadFile("../../blog.sql")
	if err != nil {
		panic(err)
	}

	requests := strings.Split(string(file), ";\n")
	for _, request := range requests {
		if request != "" {
			_, err := articleRepo.Db.Exec(request)
			if err != nil {
				log.Fatal(err)
			}
		}
	}
}

func clearTable() {
	articleRepo.Db.Exec("DELETE FROM articles")
	articleRepo.Db.Exec("ALTER TABLE articles AUTO_INCREMENT = 1")
	articleRepo.Db.Exec("DELETE FROM tags")
	articleRepo.Db.Exec("ALTER TABLE tags AUTO_INCREMENT = 1")
	articleRepo.Db.Exec("DELETE FROM articles_tags")
}