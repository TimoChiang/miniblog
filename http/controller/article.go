package controller

import (
	"fmt"
	"github.com/gorilla/mux"
	"miniblog/domain/models"
	"miniblog/domain/service"
	"net/http"
	"strconv"
)

type ArticleHandler struct {
	Base
	Service *service.ArticleService
	UserService *service.UserService
}

type NewArticlePageData struct {
	BlogName string
	Errors string
	User *service.User
}

func (h *ArticleHandler) GetArticle (w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		fmt.Printf("convert id fail: %v\n", err)
		return
	}
	article, err := h.Service.GetArticle(id)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		t := h.LoadPageTemplate("404")
		t.Execute(w, nil)
		return
	}

	t := h.LoadPageTemplate("detail")

	data := struct {
		BlogName string
		Article *models.Article // is it good?
		User *service.User
	}{
		BLOG_NAME,
		article,
		h.UserService.GetUser(r),

	}

	if err := t.Execute(w, data); err != nil {
		fmt.Printf("execute template fail: %v\n", err)
	}
}

func (h *ArticleHandler) GetArticles (w http.ResponseWriter, r *http.Request) {
	articles, err := h.Service.GetAllArticle()
	if err != nil {
		fmt.Printf("get articles fail: %v\n", err)
		w.WriteHeader(http.StatusNotFound)
		t := h.LoadPageTemplate("404")
		t.Execute(w, nil)
	}

	t := h.LoadPageTemplate("list")

	data := struct {
		BlogName string
		Articles map[int]*models.Article
		User *service.User
	}{
		BLOG_NAME,
		articles,
		h.UserService.GetUser(r),
	}

	if err := t.Execute(w, data); err != nil {
		fmt.Printf("execute template fail: %v\n", err)
	}

}

func (h *ArticleHandler) NewArticle (w http.ResponseWriter, r *http.Request) {
	if user := h.UserService.GetUser(r); user != nil {
		t := h.LoadPageTemplate("create")
		data := NewArticlePageData{BLOG_NAME, "", user}
		if err := t.Execute(w, data); err != nil {
			fmt.Printf("execute template fail: %v\n", err)
		}
	}else{
		fmt.Printf("please sign in before create article! \n", )
		http.Redirect(w, r, "/login/", http.StatusSeeOther)
	}
	return
}


func (h *ArticleHandler) CreateArticle (w http.ResponseWriter, r *http.Request){
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
	lastInsertID, err := h.Service.CreateArticle(title, description)
	if err != nil {
		fmt.Printf("create article data fail: %v\n", err)
		return
	}
	fmt.Printf("get last article id: %v\n", lastInsertID)
	http.Redirect(w, r, "/article/" + strconv.Itoa(int(lastInsertID)), http.StatusSeeOther)
	return
}
