package web

import (
	"fmt"
	"net/http"
)

func Info(w http.ResponseWriter, r *http.Request) {
	t, err := getCommonTemplate("web/template/info.html")
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