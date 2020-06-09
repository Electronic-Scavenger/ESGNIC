package web

import (
	"encoding/json"
	"github.com/sirupsen/logrus"
	"net/http"
)

func WriteResponse(w http.ResponseWriter, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	b, err := json.Marshal(data)
	if err != nil {
		logrus.Warnf("error when marshal response message, error:%v\n", err)
	}
	w.Write(b)
}

func WriteResponseCode(w http.ResponseWriter, code int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	b, err := json.Marshal(data)
	if err != nil {
		logrus.Warnf("error when marshal response message, error:%v\n", err)
	}
	w.WriteHeader(code)
	w.Write(b)
}
