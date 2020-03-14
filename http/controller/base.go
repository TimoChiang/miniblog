package controller

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
	"text/template"
)

const BLOG_NAME  = "Timo's Blog!!"

type User struct {
	Name string
}

func GetUser(r *http.Request) *User{
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

func loadPageTemplate(name string) *template.Template {
	// must be first in allFiles
	allFiles := []string{
		"./views/pages/" + name + ".tmpl",
	}
	files, err := ioutil.ReadDir("./views/templates")
	if err != nil {
		fmt.Println(err)
	}
	for _, file := range files {
		filename := file.Name()
		if strings.HasSuffix(filename, ".tmpl") {
			allFiles = append(allFiles, "./views/templates/"+filename)
		}
	}
	ts, err := template.ParseFiles(allFiles...)
	if err != nil {
		log.Fatal(err)
	}
	return ts
	//// load base templates
	//t, err := template.ParseGlob("./views/templates/*.tmpl")
	//if err != nil {
	//	log.Fatal(err)
	//}
	//
	//t, err = t.ParseFiles("./views/pages/" + name + ".tmpl")
	//if err != nil {
	//	log.Fatal(err)
	//}
	//return t
}