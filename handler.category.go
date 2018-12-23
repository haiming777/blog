package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

func (a *App) createCategoryHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		outputJSON(w, APIStatus{
			ErrCode:    -1,
			ErrMessage: "request method is error",
		})

		return
	}

	data, err := ioutil.ReadAll(ioutil.NopCloser(r.Body))

	if err != nil {
		outputJSON(w, APIStatus{
			ErrCode:    -1,
			ErrMessage: "read request body error",
		})
		return
	}
	defer r.Body.Close()

	req := struct {
		ParentID uint   `json:"parent_id"`
		Name     string `json:"name"`
	}{}

	err = json.Unmarshal(data, &req)
	if err != nil {
		outputJSON(w, APIStatus{
			ErrCode:    -1,
			ErrMessage: "request body unmashal to struct error",
		})
		return
	}

	if req.Name == "" {
		outputJSON(w, APIStatus{
			ErrCode:    -2,
			ErrMessage: "category name is empty",
		})
		return
	}

	//存储category

	cate := &Category{
		Name:     req.Name,
		ParentID: req.ParentID,
	}

	err = a.createCategory(cate)
	if err != nil {
		outputJSON(w, APIStatus{
			ErrCode:    -3,
			ErrMessage: fmt.Sprintf("db create categroy info error :%s", err.Error()),
		})
		return
	}

	outputJSON(w, APIStatus{ErrCode: 0, Data: cate})

}
