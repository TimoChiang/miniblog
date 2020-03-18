package controller

import (
	"fmt"
	"miniblog/domain/service"
	"net/http"
	"time"
)

type UserHandler struct {
	Base
	Service *service.UserService
}

type LoginPageData struct {
	BlogName string
	Errors string
	User *service.User
}

func (h *UserHandler) SignIn(w http.ResponseWriter, r *http.Request) {
	if user := h.Service.GetUser(r); user == nil {
		t := h.LoadPageTemplate("login")
		data := LoginPageData{BLOG_NAME, "", user}
		if err := t.Execute(w, data); err != nil {
			fmt.Printf("execute template fail: %v\n", err)
		}
	}else{
		fmt.Printf("already sing in finish! \n", )
		http.Redirect(w, r, "/", http.StatusSeeOther)
	}
	return
}

func (h *UserHandler) PostSignIn (w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		fmt.Printf("get post data fail: %v\n", err)
		return
	}
	fmt.Printf("Post from website! r.PostFrom = %v\n", r.PostForm)
	name := r.FormValue("name")
	password := r.FormValue("password")
	fmt.Printf("Name = %s\n", name)
	fmt.Printf("Password = %s\n", password)

	//TODO: Validation
	if name != "demo" || password != "demo" {
		t := h.LoadPageTemplate("login")
		data := LoginPageData{BLOG_NAME, "Please check your name and password", h.Service.GetUser(r)}
		if err := t.Execute(w, data); err != nil {
			fmt.Printf("execute template fail: %v\n", err)
		}
		return
	}
	setCookie(w)
	http.Redirect(w, r, "/", http.StatusSeeOther)
	return
}

func (h *UserHandler) SignOut (w http.ResponseWriter, r *http.Request){
	if user := h.Service.GetUser(r); user != nil {
		deleteCookie(w)
	}else{
		fmt.Printf("now not login! \n", )
	}
	http.Redirect(w, r, "/login", http.StatusSeeOther)
	return
}

func setCookie(w http.ResponseWriter) {
	c := http.Cookie{
		Name:     "user",
		Value:    "demo",
		Path:     "/",
		HttpOnly: true,
		//Secure:   true,
		MaxAge: 86400}
	http.SetCookie(w, &c)
	fmt.Printf("cookie is created! => %v\n", c)
}

func deleteCookie(w http.ResponseWriter) {
	c := &http.Cookie{
		Name:     "user",
		Value:    "",
		Path:     "/",
		Expires: time.Unix(0, 0),
		HttpOnly: true,
	}
	http.SetCookie(w, c)
}