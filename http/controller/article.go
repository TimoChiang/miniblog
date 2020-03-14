package controller

import (
	"fmt"
	"github.com/gorilla/mux"
	m "miniblog/models"
	"net/http"
	"strconv"
)


func GetArticle (w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		fmt.Printf("convert id fail: %v\n", err)
		return
	}
	article, err := m.GetSingleArticle(id)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		t := loadPageTemplate("404")
		t.Execute(w, nil)
		return
	}

	t := loadPageTemplate("detail")

	data := struct {
		BlogName string
		Article *m.Article
		User *User
	}{
		BLOG_NAME,
		article,
		GetUser(r),

	}

	if err := t.Execute(w, data); err != nil {
		fmt.Printf("execute template fail: %v\n", err)
	}
}

func GetArticles (w http.ResponseWriter, r *http.Request) {
	articles, err := m.GetAllArticle()
	if err != nil {
		fmt.Printf("get articles fail: %v\n", err)
		w.WriteHeader(http.StatusNotFound)
		t := loadPageTemplate("404")
		t.Execute(w, nil)
	}

	t := loadPageTemplate("list")

	data := struct {
		BlogName string
		Articles map[int]*m.Article
		User *User
	}{
		BLOG_NAME,
		articles,
		GetUser(r),
	}

	if err := t.Execute(w, data); err != nil {
		fmt.Printf("execute template fail: %v\n", err)
	}

}

func NewArticle (w http.ResponseWriter, r *http.Request) {
	if user := GetUser(r); user != nil {
		t := loadPageTemplate("create")
		data := LoginPageData{BLOG_NAME, "", user}
		if err := t.Execute(w, data); err != nil {
			fmt.Printf("execute template fail: %v\n", err)
		}
	}else{
		fmt.Printf("please sign in before create article! \n", )
		http.Redirect(w, r, "/login/", http.StatusSeeOther)
	}
	return
}


func CreateArticle (w http.ResponseWriter, r *http.Request){
	if err := r.ParseForm(); err != nil {
		fmt.Printf("get post data fail: %v\n", err)
		return
	}
	fmt.Printf("Post from website! r.PostFrom = %v\n", r.PostForm)
	title := r.FormValue("title")
	description := r.FormValue("description")
	fmt.Printf("Name = %s\n", title)
	fmt.Printf("Password = %s\n", description)

	//TODO: Validation
	lastInsertID, err := m.CreateArticle(title, description)
	if err != nil {
		fmt.Printf("create article data fail: %v\n", err)
		return
	}
	fmt.Printf("get last article id: %v\n", lastInsertID)
	http.Redirect(w, r, "/article/" + strconv.Itoa(int(lastInsertID)), http.StatusSeeOther)
	return
}
