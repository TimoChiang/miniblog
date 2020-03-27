package repository

import (
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql" // use _ to execute package's init function only
	"log"
	"os"
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

var articleRepo *ArticleRepository

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
	articleRepo = new(ArticleRepository)
	articleRepo.Db = db
	ensureTableExists()
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

func ensureTableExists() {
	if _, err := articleRepo.Db.Exec(createArticleQuery); err != nil {
		log.Fatal(err)
	}
}

func clearTable() {
	articleRepo.Db.Exec("DELETE FROM articles")
	articleRepo.Db.Exec("ALTER TABLE articles AUTO_INCREMENT = 1")
}


const createArticleQuery = `
	CREATE TABLE IF NOT EXISTS articles (
	id int(11) unsigned NOT NULL AUTO_INCREMENT,
	title varchar(255) CHARACTER SET utf8mb4 DEFAULT NULL,
	description text CHARACTER SET utf8mb4,
	created_at datetime DEFAULT CURRENT_TIMESTAMP,
	updated_at datetime DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
	PRIMARY KEY (id)
	) ENGINE=InnoDB AUTO_INCREMENT=1 DEFAULT CHARSET=latin1;
`
