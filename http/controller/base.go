package controller

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
	"text/template"
)

type Base struct {}

func (b *Base) LoadPageTemplate(name string) (*template.Template, error) {
	// must be first in allFiles
	allFiles := []string{
		"./views/pages/" + name + ".tmpl",
	}
	files, err := ioutil.ReadDir("./views/templates")
	if err != nil {
		return nil, err
	}
	for _, file := range files {
		filename := file.Name()
		if strings.HasSuffix(filename, ".tmpl") {
			allFiles = append(allFiles, "./views/templates/"+filename)
		}
	}
	ts, err := template.ParseFiles(allFiles...)
	if err != nil {
		return nil, err
	}
	return ts, nil
}

func CheckHealth (w http.ResponseWriter, r *http.Request){
	fmt.Fprintf(w, "I'm fine!")
}