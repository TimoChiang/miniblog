package controller

import (
	"fmt"
	"io/ioutil"
	"log"
	"strings"
	"text/template"
)

type Base struct {
}

const BLOG_NAME  = "Timo's Blog!!"

func (b *Base) LoadPageTemplate(name string) *template.Template {
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
}