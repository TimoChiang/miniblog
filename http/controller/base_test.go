package controller

import (
	"fmt"
	"io/ioutil"
	"os"
	"testing"
)

// move dir to root to load templates successfully
func init() {
	if err := os.Chdir("../.."); err != nil {
		panic(err)
	}
}

func TestLoadPageTemplate(t *testing.T) {
	base := new(Base)
	mainFileName := "login"
	template := base.LoadPageTemplate(mainFileName)
	allFiles := map[string]bool{}
	for _, file := range template.Templates() {
		allFiles[file.ParseName] = true
	}

	//check if load all files in templates folder
	files, err := ioutil.ReadDir("./views/templates")
	if err != nil {
		fmt.Println(err)
	}

	for _, file := range files {
		filename := file.Name()
		if _, ok := allFiles[filename]; ok == false  {
			t.Errorf("%v should be load but not found in templates", filename)
		}
	}

	//check if load expected file
	if _, ok := allFiles[mainFileName + ".tmpl"]; ok == false  {
		t.Errorf("file '%v' should be load but not found in templates", mainFileName)
	}
}