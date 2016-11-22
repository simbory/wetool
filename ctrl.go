package main

import (
	"regexp"
	"strings"
	"os"
	"bytes"
	"io/ioutil"
	"html/template"
)

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
	if IsFile(ctrlDir + fileName) {
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

