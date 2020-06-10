package web

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"github.com/Electronic-Scavenger/ESGNIC/orm"
	"github.com/sirupsen/logrus"
	"hash"
	"html/template"
	"net/http"
	"time"
)

const cookieMaxAge = 3600
const cookieKey = "fds46489d6"

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

func LoginAction(w http.ResponseWriter, r *http.Request) {
	req := UserForm{}
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		logrus.Error(err)
		WriteResponseCode(w, 500, "格式错误")
		return
	}
	req.Password = passwordHash(req.Password)
	userinfo := orm.User{}
	if err = orm.DB.Where(&orm.User{QQ: req.QQ, Password: req.Password}).Take(&userinfo).Error; err != nil {
		WriteResponseCode(w, 500, "用户名或密码错误")
		return
	}
	login(w, userinfo.QQ)
	WriteResponse(w, nil)
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

	userinfo.Password = passwordHash(req.Password)
	orm.DB.Model(&orm.User{}).Update(&userinfo)

	WriteResponseCode(w, 204, nil)
	return
}

func LogoutAction(w http.ResponseWriter, r *http.Request) {
	logout(w)
}

func passwordHash(password string) string {
	return fmt.Sprintf("%X", sha256.Sum256([]byte(password)))
}

func writeCookie(w http.ResponseWriter, name, value string, maxAge int) {
	cookie := http.Cookie{
		Name:   name,
		Value:  value,
		MaxAge: maxAge,
		SameSite: http.SameSiteStrictMode,
	}
	http.SetCookie(w, &cookie)
}

func login(w http.ResponseWriter, qq string) {
	expire := time.Now().Add(time.Second * cookieMaxAge).Format(time.RFC3339)
	hash := hmac.New(func() hash.Hash {
		return sha256.New()
	}, []byte(cookieKey))
	hash.Write([]byte(qq))
	token := fmt.Sprintf("%x", hash.Sum([]byte(expire)))
	writeCookie(w, "t", expire, cookieMaxAge)
	writeCookie(w, "qq", qq, cookieMaxAge)
	writeCookie(w, "token", token, cookieMaxAge)
}

func logout(w http.ResponseWriter) {
	writeCookie(w, "t", "", -1)
	writeCookie(w, "qq", "", -1)
	writeCookie(w, "token", "", -1)
}

func auth(w http.ResponseWriter, r *http.Request) (userinfo *orm.User) {
	defer func() {
		if userinfo == nil {
			logout(w)
			http.Redirect(w, r, "/login", http.StatusFound)
		}
	}()
	qq, _ := r.Cookie("qq")
	token, _ := r.Cookie("token")
	expire, _ := r.Cookie("t")
	if qq != nil && token != nil && expire != nil {
		// validate
		hash := hmac.New(func() hash.Hash {
			return sha256.New()
		}, []byte(cookieKey))
		hash.Write([]byte(qq.Value))
		targetToken := fmt.Sprintf("%x", hash.Sum([]byte(expire.Value)))
		if token.Value != targetToken {
			logrus.Error("invalid token")
			return nil
		}

		t, _ := time.Parse(time.RFC3339, expire.Value)
		if time.Now().After(t) {
			logrus.Error("cookie timeout")
			return nil
		}
		userinfo = new(orm.User)
		if err := orm.DB.Where(&orm.User{QQ: qq.Value}).Take(userinfo).Error; err != nil {
			logrus.WithField("qq", qq.Value).Error("user not exist")
			return nil
		}
	}
	return
}
