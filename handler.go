package main

import (
	"log"
	"net/http"
)

func (a *App) initHandler() {
	a.registerRouter("/user/create", a.createUserHandler)
	a.registerRouter("/user/update", a.updataUserHandler)
	a.registerRouter("/signin", a.signin)

}

func (a *App) registerRouter(pattern string, handler func(http.ResponseWriter, *http.Request)) {
	log.Printf("handler:[%s]\n", pattern)
	http.HandleFunc(pattern, handler)
}

func (a *App) listenAndServer() {
	if err := a.HTTPServer.ListenAndServe(); err != nil {
		log.Fatal("[ERR] - ListenAndServe:", err)
	}
}
