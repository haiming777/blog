package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
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

func (a *App) categoryListHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		outputJSON(w, APIStatus{
			ErrCode:    -1,
			ErrMessage: "request method is error",
		})
		return
	}

	//查询所有category数据并返回
	categories, err := a.queryCategoriesFirstLevel()
	if err != nil {
		outputJSON(w, APIStatus{
			ErrCode:    -1,
			ErrMessage: fmt.Sprintf("query category error:%s", err.Error()),
		})

		return
	}

	outputJSON(w, APIStatus{
		ErrCode: 0,
		Data:    categories,
	})
}

func (a *App) queryCategoriesByParentIDHandler(w http.ResponseWriter, r *http.Request) {
	resp := &APIStatus{ErrCode: 0}
	defer outputJSON(w, resp)
	if r.Method != "GET" {
		resp.ErrCode = -1
		resp.ErrMessage = "request method is error"
		return
	}

	if r.Body == nil {
		resp.ErrCode = -2
		resp.ErrMessage = "request body can not be nil"
		return
	}

	parentID, err := strconv.ParseUint(r.URL.Query().Get("id"), 10, 64)
	if err != nil {
		resp.ErrCode = -3
		resp.ErrMessage = "request param strconv error"
		return
	}

	if parentID == 0 {
		resp.ErrCode = -3
		resp.ErrMessage = "request param error"
		return
	}

	categories, err := a.queryCategoriesByParentID(uint(parentID))
	if err != nil {
		resp.ErrCode = -4
		resp.ErrMessage = "query category by parent id error"
		return
	}
	resp.Data = categories
	return
}
