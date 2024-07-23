package server

import (
	"github.com/gorilla/mux"
	"net/http"
)

type Router interface {
	ServeHttp(w http.ResponseWriter, r *http.Request)
	Get(string, http.HandlerFunc)
	Post(string, http.HandlerFunc)
}

type router struct {
	mux *mux.Router
}

func (router *router) ServeHttp(w http.ResponseWriter, r *http.Request) {
	router.mux.ServeHTTP(w, r)
}

func (router *router) Get(path string, handlerFunc http.HandlerFunc) {
	router.mux.HandleFunc(path, handlerFunc).Methods(http.MethodGet)
}

func (router *router) Post(path string, handlerFunc http.HandlerFunc) {
	router.mux.HandleFunc(path, handlerFunc).Methods(http.MethodPost)
}
