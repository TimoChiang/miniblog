package controller

import (
	"fmt"
	"log"
	"miniblog/domain/service"
	"net/http"
)

type UserHandler struct {
	Base
	Service *service.UserService
}

type LoginPageData struct {
	Errors string
	User *service.User
}

func (h *UserHandler) SignIn(w http.ResponseWriter, r *http.Request) {
	if user := h.Service.GetUser(r); user == nil {
		t, err := h.LoadPageTemplate("login")
		if err != nil {
			log.Panicln(err)
		}
		data := LoginPageData{"", user}
		if err := t.Execute(w, data); err != nil {
			fmt.Printf("execute template fail: %v\n", err)
		}
	}else{
		fmt.Printf("already sing in finish! \n", )
		http.Redirect(w, r, "/", http.StatusSeeOther)
	}
}

func (h *UserHandler) PostSignIn (w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		log.Fatalf("get post data fail: %v", err)
	}
	fmt.Printf("Post from website! r.PostFrom = %v\n", r.PostForm)
	name := r.FormValue("name")
	password := r.FormValue("password")
	fmt.Printf("Name = %s\n", name)
	fmt.Printf("Password = %s\n", password)

	//TODO: Validation
	if name != "demo" || password != "demo" {
		t, err := h.LoadPageTemplate("login")
		if err != nil {
			log.Panicln(err)
		}
		data := LoginPageData{"Please check your name and password", h.Service.GetUser(r)}
		if err := t.Execute(w, data); err != nil {
			fmt.Printf("execute template fail: %v\n", err)
		}
		return
	}
	h.Service.SetCookie(w)
	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func (h *UserHandler) SignOut (w http.ResponseWriter, r *http.Request){
	if user := h.Service.GetUser(r); user != nil {
		h.Service.DeleteCookie(w)
	}else{
		fmt.Printf("now not login! \n", )
	}
	http.Redirect(w, r, "/login", http.StatusSeeOther)
}



