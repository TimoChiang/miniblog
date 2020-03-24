package service

import (
	"fmt"
	"net/http"
	"time"
)

type UserService struct {
	Name string
}

type User struct {
	Name string
}
const cookieName = "user"

func (s *UserService) GetUser(r *http.Request) *User{
	var user *User
	cookie, err := r.Cookie(cookieName)
	// TODO: cookie encode
	if err != nil {
		fmt.Printf("get cookie failed err:%v\n", err)
	}else{
		fmt.Printf("get cookie: %v\n", cookie.Value)
		user = new(User)
		user.Name = cookie.Value
		fmt.Printf("initialize user : %v\n", user)
	}
	return user
}

func (s *UserService) SetCookie(w http.ResponseWriter) {
	c := &http.Cookie{
		Name:     cookieName,
		Value:    "demo",
		Path:     "/",
		HttpOnly: true,
		//Secure:   true,
		MaxAge: 86400}
	http.SetCookie(w, c)
	fmt.Printf("cookie is created! => %v\n", c)
}

func (s *UserService) DeleteCookie(w http.ResponseWriter) {
	c := &http.Cookie{
		Name:     cookieName,
		Value:    "",
		Path:     "/",
		HttpOnly: true,
		Expires: time.Unix(0, 0),
		MaxAge: 0,
	}
	http.SetCookie(w, c)
	fmt.Println("cookie is deleted!")
}
