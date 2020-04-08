package controller

import (
	"fmt"
	"github.com/gorilla/mux"
	"log"
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
	Errors map[string][]string
	User *service.User
	Inputs *models.Article
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
		log.Panicln(err)
	}
	if article != nil {
		t, err := h.LoadPageTemplate("detail")
		if err != nil {
			log.Panicln(err)
		}
		data := struct {
			Article *models.Article // is it good?
			User *service.User
		}{
			article,
			h.UserService.GetUser(r),
		}

		if err := t.Execute(w, data); err != nil {
			fmt.Printf("execute template fail: %v\n", err)
		}
	}else{
		// TODO 404 function
		w.WriteHeader(http.StatusNotFound)
		t, err := h.LoadPageTemplate("404")
		if err != nil {
			log.Panicln(err)
		}

		if err := t.Execute(w, nil); err != nil {
			fmt.Printf("execute template fail: %v\n", err)
		}
	}
}

func (h *ArticleHandler) GetArticles (w http.ResponseWriter, r *http.Request) {
	articles, err := h.Service.GetAllArticle()
	if err != nil {
		fmt.Printf("get articles fail: %v\n", err)
		w.WriteHeader(http.StatusNotFound)
		t, err := h.LoadPageTemplate("404")
		if err != nil {
			log.Panicln(err)
		}
		t.Execute(w, nil)
	}

	t, err := h.LoadPageTemplate("list")
	if err != nil {
		log.Panicln(err)
	}

	data := struct {
		Articles map[int]*models.Article
		User *service.User
	}{
		articles,
		h.UserService.GetUser(r),
	}

	if err := t.Execute(w, data); err != nil {
		fmt.Printf("execute template fail: %v\n", err)
	}

}

func (h *ArticleHandler) NewArticle (w http.ResponseWriter, r *http.Request) {
	if user := h.UserService.GetUser(r); user != nil {
		t, err := h.LoadPageTemplate("create")
		if err != nil {
			log.Panicln(err)
		}
		data := NewArticlePageData{nil, user, nil}
		if err := t.Execute(w, data); err != nil {
			fmt.Printf("execute template fail: %v\n", err)
		}
	}else{
		fmt.Printf("please sign in before create article! \n", )
		http.Redirect(w, r, "/login", http.StatusFound)
	}
	return
}

func (h *ArticleHandler) CreateArticle (w http.ResponseWriter, r *http.Request){
	if user := h.UserService.GetUser(r); user != nil {
		articleStruct, err := h.Service.LoadArticleStruct(r)
		if err != nil {
			fmt.Printf("get post data fail: %v\n", err)
		}

		// validation failed
		if errorMessage := h.Service.V.Exec(articleStruct); errorMessage != nil {
			fmt.Printf("get error message: %v\n", errorMessage)
			data := NewArticlePageData{errorMessage, user, articleStruct}
			w.WriteHeader(http.StatusUnprocessableEntity)
			t, err := h.LoadPageTemplate("create")
			if err != nil {
				log.Panicln(err)
			}

			if err := t.Execute(w, data); err != nil {
				fmt.Printf("execute template fail: %v\n", err)
			}
			return
		}

		// validation success, create article
		lastInsertID, err := h.Service.CreateArticle(articleStruct)
		if err != nil {
			fmt.Printf("create article data fail: %v\n", err)
			return
		}
		fmt.Printf("get last article id: %v\n", lastInsertID)
		http.Redirect(w, r, "/article/" + strconv.Itoa(int(lastInsertID)), http.StatusSeeOther)
	}else{
		fmt.Printf("please sign in before create article! \n", )
		http.Redirect(w, r, "/login", http.StatusFound)
	}
}
