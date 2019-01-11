package main

var app *App

func init() {
	app = NewApp(&cfg)
	app.initDB()
	app.initHandler()
}

func main() {
	// http.Handle(pattern, handler)
	// http.ServeFile(w, r, name)
	// http.ListenAndServe(addr, handler)

	app.listenAndServer()
}
