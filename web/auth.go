package web

import (
	"fmt"
	"html/template"
	"net/http"
)

func Login(w http.ResponseWriter, r *http.Request) {
	t, err := template.ParseFiles("web/template/login.html", "web/template/common.html")
	if err != nil {
		fmt.Println(err)
		return
	}
	err = t.Execute(w, BaseTemplateArgs{Title: TITLE})
	if err != nil {
		fmt.Println(err)
	}
}
