package web

import (
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"github.com/Electronic-Scavenger/ESGNIC/orm"
	"github.com/sirupsen/logrus"
	"html/template"
	"net/http"
)

type UserForm struct {
	QQ       string `json:"qq"`
	Password string `json:"password"`
}

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

func Register(w http.ResponseWriter, r *http.Request) {
	t, err := template.ParseFiles("web/template/register.html", "web/template/common.html")
	if err != nil {
		fmt.Println(err)
		return
	}
	err = t.Execute(w, BaseTemplateArgs{Title: TITLE})
	if err != nil {
		fmt.Println(err)
	}
}

func RegisterAction(w http.ResponseWriter, r *http.Request) {
	req := UserForm{}
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		logrus.Error(err)
		WriteResponseCode(w, 500, "格式错误")
		return
	}
	userinfo := orm.User{}

	if err = orm.DB.Where(&orm.User{QQ: req.QQ}).Take(&userinfo).Error; err != nil {
		WriteResponseCode(w, 500, "无效的QQ，请联系管理员")
		return
	} else if userinfo.Password != "" {
		WriteResponseCode(w, 500, "该QQ已注册，请直接登录")
		return
	}
	
	userinfo.Password = fmt.Sprintf("%X", sha256.Sum256([]byte(req.Password)))
	orm.DB.Model(&orm.User{}).Update(&userinfo)

	WriteResponseCode(w, 204, nil)
	return
}
