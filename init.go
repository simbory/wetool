package main

import (
	"io/ioutil"
	"strings"
	"bytes"
	"os"
	"text/template"
)

func initProject() {
	files, err := ioutil.ReadDir(dir)
	if err != nil {
		println(err.Error())
	}
	for _, f := range files {
		if f.IsDir() {
			if strings.ToLower(f.Name()) == "controllers" {
				println("wetool: It seems that this project has been initialized. The folder \"controllers\" is already exist.")
				return
			}
			if strings.ToLower(f.Name()) == "content" {
				println("wetool: It seems that this project has been initialized. The folder \"content\" is already exist.")
				return
			}
			if strings.ToLower(f.Name()) == "views" {
				println("wetool: It seems that this project has been initialized. The folder \"views\" is already exist.")
				return
			}
		} else if strings.HasSuffix(strings.ToLower(f.Name()), ".go") {
			println("wetool: You should run this command in an empty folder which is under GOPATH src directory.")
			return
		} else if strings.ToLower(f.Name()) == "config.xml" {
			println("wetool: It seems that this project has been initialize. The file \"config.xml\" is already exist.")
			return
		}
	}
	mainTpl, _ := template.New("main").Parse(tplMain)
	data := map[string]string{"pkgPath": pkgPath}
	var buf = &bytes.Buffer{}
	mainTpl.Execute(buf, data)
	ioutil.WriteFile(dir+"/main.go", buf.Bytes(), 0777)
	os.MkdirAll(dir+"/content/css", 0777)
	os.MkdirAll(dir+"/content/img", 0777)
	os.MkdirAll(dir+"/content/js", 0777)
	os.MkdirAll(dir+"/controllers", 0777)
	os.MkdirAll(dir+"/models", 0777)
	os.MkdirAll(dir+"/views/home", 0777)
	os.MkdirAll(dir+"/views/shared", 0777)
	for fileName, content := range staticFiles {
		ioutil.WriteFile(dir+fileName, []byte(content), 0777)
	}
}

