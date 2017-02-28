package main

import (
	"regexp"
	"os"
	"strings"
	"bytes"
	"io/ioutil"
	"html/template"
)

func newArea(name string) {
	dirRegValidate, _ := regexp.Compile("^[a-zA-Z0-9_]+[a-zA-Z0-9_-]*$")
	if !dirRegValidate.Match([]byte(name)) {
		println("wetool: invalid area name: " + name)
		return
	}
	areaDir := dir + "/" + name
	if IsDir(areaDir) {
		println("The area \"" + name + "\" is already exist.")
		return
	}
	os.MkdirAll(areaDir, 0777)
	os.MkdirAll(areaDir+"/controllers", 0777)
	os.MkdirAll(areaDir+"/views/default", 0777)
	pkgName := strings.Replace(name, "-", "", -1)
	data := map[string]interface{}{
		"pkgName":  pkgName,
		"nsCtrlPkg": pkgPath + "/" + name + "/controllers",
		"areaName":   name,
		"startTag": template.HTML("<"),
		"endTag":   template.HTML(">"),
	}
	// write init.go file
	areaInitTpl, _ := template.New("AreaInit").Parse(tplAreaInitFile)
	bufInit := &bytes.Buffer{}
	err := areaInitTpl.Execute(bufInit, data)
	if err != nil {
		println(err.Error())
		return
	}
	ioutil.WriteFile(areaDir+"/init.go", bufInit.Bytes(), 0777)

	// write settings.xml file
	areaSettingTpl, _ := template.New("AreaSetting").Parse(tplAreaSettingFile)
	bufSetting := &bytes.Buffer{}
	err = areaSettingTpl.Execute(bufSetting, data)
	if err != nil {
		println(err.Error())
		return
	}
	ioutil.WriteFile(areaDir+"/settings.xml", bufSetting.Bytes(), 0777)

	// write controller file
	areaCtrlTpl, _ := template.New("AreaCtrl").Parse(tplAreaCtrlFile)
	bufCtrl := &bytes.Buffer{}
	err = areaCtrlTpl.Execute(bufCtrl, data)
	if err != nil {
		println(err.Error())
		return
	}
	ioutil.WriteFile(areaDir+"/controllers/defaultController.go", bufCtrl.Bytes(), 0777)

	// write view file

	// write controller file
	t := template.New("NsView")
	areaViewTpl, _ := t.Parse(tplAreaViewFile)
	bufView := &bytes.Buffer{}
	err = areaViewTpl.Execute(bufView, data)
	if err != nil {
		println(err.Error())
		return
	}
	ioutil.WriteFile(areaDir+"/views/default/index.html", bufView.Bytes(), 0777)
}

