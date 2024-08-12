package server

import (
	"net/http"

	"github.com/gorilla/mux"
)

type Router interface {
	ServeHttp(w http.ResponseWriter, r *http.Request)
	Get(string, http.HandlerFunc)
	Post(string, http.HandlerFunc)
	Delete(string, http.HandlerFunc)
	Update(string, http.HandlerFunc)
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

func (router *router) Delete(path string, handlerFunc http.HandlerFunc) {
	router.mux.HandleFunc(path, handlerFunc).Methods(http.MethodDelete)
}

func (router *router) Update(path string, handlerFunc http.HandlerFunc) {
	router.mux.HandleFunc(path, handlerFunc).Methods(http.MethodPut)
}
