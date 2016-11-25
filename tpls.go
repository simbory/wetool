package main

var (
	tplMain = `package main

import (
	"github.com/Simbory/wemvc"
	_ "{{.pkgPath}}/controllers"
)

func main() {
	wemvc.StaticDir("/content/")
	wemvc.Run(8080)
}`

	tplNsInitFile = `package {{.pkgName}}

import (
	"github.com/Simbory/wemvc"
	{{.pkgName}}Ctrls "{{.nsCtrlPkg}}"
)

func init() {
	ns := wemvc.Namespace("{{.nsName}}")
	ns.Route("/default/{{.startTag}}action=index{{.endTag}}", {{.pkgName}}Ctrls.DefaultController{})
}`

	tplNsSettingFile = `<?xml version="1.0" encoding="utf-8"?>
<settings>
	<add key="Namespace.Name" value="{{.nsName}}" />
</settings>`

	tplNsCtrlFile = `package controllers

import (
	"github.com/Simbory/wemvc"
)

// DefaultController the default controller for 'admin' namespace
type DefaultController struct {
	wemvc.Controller
}

// GetIndex the index action for http GET method
func (def DefaultController) GetIndex() interface{} {
	return def.View()
}

// PostIndex the index action for http POST method
func (def DefaultController) PostIndex() interface{} {
	def.ViewData["PostMsg"] = def.Request().Form.Get("msg")
	return def.View()
}
`
	tplCtrlFile = `package controllers

import "github.com/Simbory/wemvc"

type {{.structName}} struct {
	wemvc.Controller
}

// GetIndex the index action for http GET method
func ({{.structParam}} {{.structName}}) GetIndex() interface{} {
	return {{.structParam}}.View()
}`

	tplNsViewFile = `<!DOCTYPE html>
<html>
<head>
    <meta http-equiv="Content-Type" content="text/html; charset=utf-8" />
    <meta charset="utf-8" />
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>Default page for namespace {{.nsName}} - My wemvc application</title>
</head>
<body>
    <div>
        <a href="/">Home</a> &gt; {{.nsName}} &gt; <a href="#">Default</a> &gt; Index
    </div>
    <h1>The Default page for namespace "{{.nsName}}"</h1>
    {{"{{"}}if .PostMsg{{"}}"}}<p style="color:red;">{{"{{"}}.PostMsg{{"}}"}}</p>{{"{{"}}end{{"}}"}}
    <form action="" method="POST">
    	<label for="msg">Your post message:</label>
    	<input type="text" name="msg" id="msg"/>
    	<button type="submit">Submit</button>
    </form>
</body>
</html>`

	staticFiles = map[string]string{
		"/config.xml": `<?xml version="1.0" encoding="utf-8"?>
<configuration>
    <defaultUrl>index.html;index.htm</defaultUrl>
    <!-- Connection string setting -->
    <connStrings>
        <add name="default" type="mysql" connString="connString1"/>
    </connStrings>
    <!-- Application setting -->
    <settings>
        <add key="DebugMode" value="true" />
    </settings>
    <session manager="memory" cookieName="Session_ID" enableSetCookie="true" sessionIDLength="32" />
</configuration>`,
		"/controllers/homeController.go": `package controllers

import "github.com/Simbory/wemvc"

// HomeController home controller
type HomeController struct {
	wemvc.Controller
}

// GetIndex use http GET method to visit index action
func (home HomeController) GetIndex() interface{} {
	home.ViewData["Message"] = "Welcome to WEMVC 1.0"
	return home.View()
}

// GetAbout use http GET method to visit about action
func (home HomeController) GetAbout() interface{} {
	home.ViewData["Message"] = "About wemvc"
	return home.View()
}

// GetAbout use http GET method to visit about action
func (home HomeController) GetContact() interface{} {
	home.ViewData["Message"] = "Contact us"
	return home.View()
}`,
		"/models/models.go": `package models

type LoginModel struct {
	Email 		string "field:email"
	Password 	string "field:pwd"
	RememberMe 	string "field:remember"
}

`,
		"/controllers/init.go": `package controllers

import "github.com/Simbory/wemvc"

func init() {
	wemvc.Route("/<action=index>", HomeController{})
}`,
		"/views/shared/layout.html": `<!DOCTYPE html>
<html>
<head>
    <meta http-equiv="Content-Type" content="text/html; charset=utf-8" />
    <meta charset="utf-8" />
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>{{template "PageTitle" .}} - My wemvc application</title>
    <link href="//cdn.bootcss.com/bootstrap/3.3.7/css/bootstrap.min.css" rel="stylesheet" />
    <link href="/content/css/site.css" rel="stylesheet" />
    <script src="//cdn.bootcss.com/modernizr/2.8.3/modernizr.min.js"></script>
    {{template "Head" .}}
</head>
<body>
    <div class="navbar navbar-inverse navbar-fixed-top">
        <div class="container">
            <div class="navbar-header">
                <button type="button" class="navbar-toggle" data-toggle="collapse" data-target=".navbar-collapse">
                    <span class="icon-bar"></span>
                    <span class="icon-bar"></span>
                    <span class="icon-bar"></span>
                </button>
                <a class="navbar-brand" href="/">Application Name</a>
            </div>
            <div class="navbar-collapse collapse">
                <ul class="nav navbar-nav">
                    <li><a href="/">Home</a></li>
                    <li><a href="/about/">About</a></li>
                    <li><a href="/contact/">Contact</a></li>
                </ul>
            </div>
        </div>
    </div>
    <div class="container body-content">
        {{template "Body".}}
        <hr/>
        <footer>
            <p>&copy; 2016 - My wemvc application</p>
        </footer>
    </div>
    <script src="//cdn.bootcss.com/jquery/2.2.4/jquery.min.js"></script>
    <script src="//cdn.bootcss.com/bootstrap/3.3.7/js/bootstrap.min.js"></script>
    <script src="//cdn.bootcss.com/respond.js/1.4.2/respond.min.js"></script>
    {{template "Scripts" .}}
</body>
</html>`,
		"/views/home/index.html": `{{define "PageTitle"}}Homepage{{end}}

{{define "Head"}}{{end}}

{{define "Body"}}
    <div class="jumbotron" style="min-height: 300px;"></div>
    <h1>{{.Message}}</h1>
{{end}}

{{define "Scripts"}}
    <script src="/content/js/home.js"></script>
{{end}}

{{template "shared/layout.html" .}}`,
		"/views/home/contact.html": `{{define "PageTitle"}}Homepage{{end}}

{{define "Head"}}{{end}}

{{define "Body"}}
    <h1>{{.Message}}</h1>
    <ol>
        <li>Email: xxx@xxx.com</li>
        <li>Phone number: 4000 0000</li>
        <li>Address: No. xxx, XXX Road</li>
    </ol>
{{end}}

{{define "Scripts"}}{{end}}

{{template "shared/layout.html" .}}`,
		"/views/home/about.html": `{{define "PageTitle"}}Homepage{{end}}

{{define "Head"}}{{end}}

{{define "Body"}}
    <h1>{{.Message}}</h1>
    <article>
        Input your message here...
    </article>
{{end}}

{{define "Scripts"}}{{end}}

{{template "shared/layout.html" .}}`,
		"/content/js/home.js":   "",
		"/content/css/site.css": `body {padding-top: 50px;padding-bottom: 20px;}.body-content {padding-left: 15px;padding-right: 15px;}`,
	}
)
