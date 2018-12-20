package main

import (
	"net/http"
)

// createUserHandler 创建用户
func (a *App) createUserHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		outputJSON(w, APIStatus{
			ErrCode:    -1,
			ErrMessage: "Method not acceptable",
		})
		return
	}

	outputJSON(w, APIStatus{
		ErrCode: 0,
	})

}
