package main

import (
	"bytes"
	"html/template"
	"io/ioutil"
	"os"
	"regexp"
	"runtime"
	"strings"

	"github.com/Simbory/wemvc"
)

type cmdType uint8

const (
	cmdEmpty cmdType = iota
	cmdInit
	cmdNewNs
	cmdNewCtrl
	cmdRun
	cmdInvalidCmd
)

func getCmdType() (cmdType, []string) {
	args := os.Args
	if len(args) == 1 {
		return cmdEmpty, nil
	}
	if args[1] == "init" {
		return cmdInit, nil
	}
	if args[1] == "ns" {
		return cmdNewNs, args[2:]
	}
	if args[1] == "ctrl" {
		return cmdNewCtrl, args[2:]
	}
	if args[1] == "run" {
		return cmdRun, nil
	}
	return cmdInvalidCmd, []string{args[1]}
}

func showHelp() {
	println("wetool is a tool that can help to create a new wemvc project easily.")
	println("")
	println("Usage:")
	println("    wetool command [arguments]")
	println("")
	println("The commands are:")
	println("    init:           Initialize the project in an empty GO package.")
	println("    ns [name]:      Create a new wemvc namespace with name '[NAME]'.")
	println("    ctrl [name]:    Create a new wemvc controller with name '[NAME]'.")
	println("    run:            Compile and run the wemvc application.")
	return
}

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
	os.MkdirAll(dir+"/views/home", 0777)
	os.MkdirAll(dir+"/views/shared", 0777)
	for fileName, content := range staticFiles {
		ioutil.WriteFile(dir+fileName, []byte(content), 0777)
	}
}

func newNs(name string) {
	dirRegValidate, _ := regexp.Compile("^[a-zA-Z0-9_]+[a-zA-Z0-9_-]*$")
	if !dirRegValidate.Match([]byte(name)) {
		println("wetool: invalid namespace name: " + name)
		return
	}
	nsDir := dir + "/" + name
	if wemvc.IsDir(nsDir) {
		println("The namespace \"" + name + "\" is already exist.")
		return
	}
	os.MkdirAll(nsDir, 0777)
	os.MkdirAll(nsDir+"/views/default", 0777)
	pkgName := strings.Replace(name, "-", "", -1)
	data := map[string]interface{}{
		"pkgName":  pkgName,
		"nsName":   name,
		"startTag": template.HTML("<"),
		"endTag":   template.HTML(">"),
	}
	// write init.go file
	nsInitTpl, _ := template.New("NsInit").Parse(tplNsInit)
	buf1 := &bytes.Buffer{}
	err := nsInitTpl.Execute(buf1, data)
	if err != nil {
		println(err.Error())
		return
	}
	ioutil.WriteFile(nsDir+"/init.go", buf1.Bytes(), 0777)

	// write controller file
	nsCtrlTpl, _ := template.New("NsCtrl").Parse(tplNsCtrlFile)
	buf2 := &bytes.Buffer{}
	nsCtrlTpl.Execute(buf2, data)
	ioutil.WriteFile(nsDir+"/defaultController.go", buf2.Bytes(), 0777)

	// write view file

	// write controller file
	nsViewTpl, _ := template.New("NsView").Parse(tplNsViewFile)
	buf3 := &bytes.Buffer{}
	nsViewTpl.Execute(buf3, data)
	ioutil.WriteFile(nsDir+"/views/default/index.html", buf3.Bytes(), 0777)
}

func newCtrl(ctrlName string) {
	ctrlValidation, _ := regexp.Compile("^[a-zA-Z]+[a-zA-Z0-9_]*$")
	if !ctrlValidation.Match([]byte(ctrlName)) || strings.ToLower(ctrlName) == "controller" {
		println("wetool: invalid ctroller name: " + ctrlName)
		return
	}
	var fileName string
	var data = make(map[string]string)
	ctrlReplaceReg, _ := regexp.Compile("[Cc][Oo][Nn][Tt][Rr][Oo][Ll][Ll][Ee][Rr]$")
	ctrlNameFix := string(ctrlReplaceReg.ReplaceAll([]byte(ctrlName), nil))
	fileName = ctrlNameFix + "Controller.go"
	structName := strings.ToUpper(ctrlNameFix[0:1]) + ctrlNameFix[1:]
	data["structParam"] = strings.ToLower(ctrlNameFix[0:1]) + ctrlNameFix[1:]
	structName = structName + "Controller"
	data["structName"] = structName
	var ctrlDir = dir + "/controllers/"
	if wemvc.IsFile(ctrlDir + fileName) {
		println("wetool: controller file is already exist: " + fileName)
		return
	}
	fInfo, err := os.Stat(ctrlDir)
	if err != nil {
		os.MkdirAll(ctrlDir, 0777)
	} else if fInfo.IsDir() {
		ctrlTpl, _ := template.New("Ctrl").Parse(tplCtrlFile)
		buf := &bytes.Buffer{}
		err = ctrlTpl.Execute(buf, data)
		if err != nil {
			println(err.Error())
			return
		} else {
			println("wetool: create new controller \"" + ctrlNameFix + "\"")
			ioutil.WriteFile(ctrlDir+fileName, buf.Bytes(), 0777)
		}
	}
}

var (
	dir       = wemvc.WorkingDir()
	goPathSrc string
	pkgPath   string
)

func main() {
	goPath := os.Getenv("GOPATH")
	if len(goPath) == 0 || !wemvc.IsDir(goPath) {
		println("Could not find the GOPATH environment variable.")
		return
	}
	if runtime.GOOS == "windows" {
		goPathSrc = goPath + "\\src"
	} else {
		goPathSrc = goPath + "/src"
	}
	if !strings.HasPrefix(dir, goPathSrc) {
		println("Error: It seems that this is not a GO package folder. This command only can be executed under an empty GO pachage folder.\r\n")
		return
	}
	pkgPath = dir[len(goPathSrc)+1:]
	pkgPath = strings.Replace(pkgPath, "\\", "/", -1)
	cmdName, args := getCmdType()
	if cmdName == cmdEmpty {
		showHelp()
		return
	}
	if cmdName == cmdInit {
		initProject()
		return
	}
	if cmdName == cmdNewNs {
		if len(args) != 1 {
			println("wetool: Invalid name of the new namespace.")
			showHelp()
			return
		}
		newNs(args[0])
		return
	}
	if cmdName == cmdNewCtrl {
		if len(args) != 1 {
			println("wetool: Invalid name of the new controller.")
			showHelp()
			return
		}
		newCtrl(args[0])
		return
	}
}
