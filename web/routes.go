package web

import (
	"github.com/gorilla/mux"
	"net/http"
)

type Route struct {
	Name        string
	Method      string
	Pattern     string
	HandlerFunc http.HandlerFunc
}

type Routes []Route

var routes = Routes{
	Route{
		Name: "Login", Method: http.MethodGet, Pattern: "/login", HandlerFunc: Login,
	},
	Route{
		Name: "Register", Method: http.MethodGet, Pattern: "/register", HandlerFunc: Register,
	},
	Route{
		Name: "RegisterAction", Method: http.MethodPost, Pattern: "/register", HandlerFunc: RegisterAction,
	},
	Route{
		Name: "Index", Method: http.MethodGet, Pattern: "/", HandlerFunc: Index,
	},
	Route{
		Name: "Index", Method: http.MethodGet, Pattern: "/index", HandlerFunc: Index,
	},
	Route{
		Name: "Info", Method: http.MethodGet, Pattern: "/info", HandlerFunc: Info,
	},
}

func NewRouter() *mux.Router {
	r := mux.NewRouter().StrictSlash(true)
	for _, route := range routes {
		r.Methods(route.Method).Name(route.Name).Path(route.Pattern).HandlerFunc(route.HandlerFunc)
	}
	r.PathPrefix("/static").Handler(http.StripPrefix("/static/", http.FileServer(http.Dir("./web/static"))))
	return r
}
