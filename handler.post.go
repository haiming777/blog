package main

import (
	"database/sql"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"time"
)

func (a *App) createPostHandler(w http.ResponseWriter, r *http.Request) {
	resp := &APIStatus{ErrCode: 0}
	defer outputJSON(w, resp)
	if r.Method != "POST" {
		resp.ErrCode = -1
		resp.ErrMessage = "request method is error"
		return
	}

	if r.Body == nil {
		resp.ErrCode = -2
		resp.ErrMessage = "request body can not be nill"
		return
	}

	data, err := ioutil.ReadAll(ioutil.NopCloser(r.Body))
	if err != nil {
		resp.ErrCode = -3
		resp.ErrMessage = "read request body is error"
		return
	}
	defer r.Body.Close()

	// summary,content,author,parentCategory,subcategory,status

	req := struct {
		Summary    string `json:"summary"`
		Content    string `json:"content"`
		Author     string `json:"author"`
		CategoryID uint   `json:"category_id"`
		Status     string `json:"status"`
	}{}

	err = json.Unmarshal(data, &req)
	if err != nil {
		resp.ErrCode = -4
		resp.ErrMessage = "request param unmarshal error"
		return
	}

	switch {
	case req.Author == "":
		resp.ErrCode = -5
		resp.ErrMessage = "post author is empty"
		return
	case req.Content == "":
		resp.ErrCode = -6
		resp.ErrMessage = "post content is empty"
		return
	case req.CategoryID == 0:
		resp.ErrCode = -7
		resp.ErrMessage = "post category id is not exist"
		return
	}

	//查询分类
	category := &Category{ID: req.CategoryID}
	err = a.queryCategoryByID(category)
	if err != nil {
		resp.ErrCode = -8
		resp.ErrMessage = "post category query error"
		return
	}

	parentCategory := &Category{ID: category.ParentID}
	err = a.queryCategoryByID(parentCategory)
	if err != nil && err != sql.ErrNoRows {
		resp.ErrCode = -9
		resp.ErrMessage = "post category query error"
		return
	}

	code := getRandomString(10)
	post := &Post{
		Code:           code,
		Summary:        req.Summary,
		Content:        req.Content,
		Author:         req.Author,
		ParentCategory: int(parentCategory.ID),
		SubCategory:    int(category.ID),
		Status:         req.Status,
		CreatedAt:      time.Now(),
	}

	//save post

	err = a.createPost(post)
	if err != nil {
		resp.ErrCode = 10
		resp.ErrMessage = "create post error"
		return
	}

	return
}
