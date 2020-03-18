package service

import (
	"fmt"
	"net/http"
)

type UserService struct {
	Name string
}

type User struct {
	Name string
}

func (s *UserService) GetUser(r *http.Request) *User{
	var user *User
	fmt.Printf("NOW user: %v\n", user)
	cookie, err := r.Cookie("user")
	if err != nil {
		fmt.Printf("get cookie failed err:%v\n", err)
		user = nil
		fmt.Printf("user reset: %v\n", user)
	}else{
		fmt.Printf("get cookie: %v\n", cookie.Value)
		if user == nil {
			user = new(User)
			user.Name = cookie.Value
			fmt.Printf("user init: %v\n", user)
		}else if user != nil && user.Name != cookie.Value {
			fmt.Printf("cookie user is wrong with data!! cookie user: %v , local user: %v \n", cookie.Value, user.Name)
			user.Name = cookie.Value
		}
	}
	return user
}
