package main

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
)

// DEBUG debug
var DEBUG bool

func (a *App) getDB() *sql.DB {
	return a.db
}

//outputJSON - output json for http response
func outputJSON(w http.ResponseWriter, data interface{}) {
	w.Header().Set("Content-Type", "application/json;charset=utf-8")
	enc := json.NewEncoder(w)
	if DEBUG {
		enc.SetIndent("", " ")
	}
	if err := enc.Encode(data); err != nil {
		log.Println("[ERR] - JSON encode error:", err)
		http.Error(w, "500", http.StatusInternalServerError)
		return
	}
}
