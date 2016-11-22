package main

import (
	"regexp"
	"os"
	"strings"
	"bytes"
	"io/ioutil"
	"html/template"
)

func newNs(name string) {
	dirRegValidate, _ := regexp.Compile("^[a-zA-Z0-9_]+[a-zA-Z0-9_-]*$")
	if !dirRegValidate.Match([]byte(name)) {
		println("wetool: invalid namespace name: " + name)
		return
	}
	nsDir := dir + "/" + name
	if IsDir(nsDir) {
		println("The namespace \"" + name + "\" is already exist.")
		return
	}
	os.MkdirAll(nsDir, 0777)
	os.MkdirAll(nsDir+"/controllers", 0777)
	os.MkdirAll(nsDir+"/views/default", 0777)
	pkgName := strings.Replace(name, "-", "", -1)
	data := map[string]interface{}{
		"pkgName":  pkgName,
		"nsCtrlPkg": pkgPath + "/" + name + "/controllers",
		"nsName":   name,
		"startTag": template.HTML("<"),
		"endTag":   template.HTML(">"),
	}
	// write init.go file
	nsInitTpl, _ := template.New("NsInit").Parse(tplNsInitFile)
	bufInit := &bytes.Buffer{}
	err := nsInitTpl.Execute(bufInit, data)
	if err != nil {
		println(err.Error())
		return
	}
	ioutil.WriteFile(nsDir+"/init.go", bufInit.Bytes(), 0777)

	// write settings.xml file
	nsSettingTpl, _ := template.New("NsSetting").Parse(tplNsSettingFile)
	bufSetting := &bytes.Buffer{}
	err = nsSettingTpl.Execute(bufSetting, data)
	if err != nil {
		println(err.Error())
		return
	}
	ioutil.WriteFile(nsDir+"/settings.xml", bufSetting.Bytes(), 0777)

	// write controller file
	nsCtrlTpl, _ := template.New("NsCtrl").Parse(tplNsCtrlFile)
	bufCtrl := &bytes.Buffer{}
	err = nsCtrlTpl.Execute(bufCtrl, data)
	if err != nil {
		println(err.Error())
		return
	}
	ioutil.WriteFile(nsDir+"/controllers/defaultController.go", bufCtrl.Bytes(), 0777)

	// write view file

	// write controller file
	nsViewTpl, _ := template.New("NsView").Parse(tplNsViewFile)
	bufView := &bytes.Buffer{}
	err = nsViewTpl.Execute(bufView, data)
	if err != nil {
		println(err.Error())
		return
	}
	ioutil.WriteFile(nsDir+"/views/default/index.html", bufView.Bytes(), 0777)
}

