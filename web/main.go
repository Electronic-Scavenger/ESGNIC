package web

import (
	"fmt"
	"github.com/sirupsen/logrus"
	"html/template"
	"net/http"
)

const TITLE = "ESGNIC"

var gatewayServer http.Server

type BaseTemplateArgs struct {
	Title string
}

func Start() {
	gatewayServer.Addr = "0.0.0.0:80"
	http.Handle("/", NewRouter())
	go gatewayServer.ListenAndServe()
}

func getCommonTemplate(filename string) (*template.Template, error) {
	t, err := template.ParseFiles(filename, "web/template/nav.html", "web/template/common.html")
	if err != nil {
		logrus.Errorf("failed to parse template: %s", err.Error())
		return nil, err
	}
	return t, nil
}

func Index(w http.ResponseWriter, r *http.Request) {
	t, err := getCommonTemplate("web/template/index.html")
	if err != nil {
		fmt.Println(err)
		return
	}
	err = t.Execute(w, BaseTemplateArgs{
		Title: TITLE,
	})
	if err != nil {
		fmt.Println(err)
	}
}
