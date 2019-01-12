package main

import (
	"database/sql"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"time"
)

//CategoryInfo 返回的分类信息
type CategoryInfo struct {
	ID             uint          `json:"id"`
	Name           string        `json:"name"`
	ParentCategory *CategoryInfo `json:"parent_category"`
}

//PostListInfo 帖子列表
type PostListInfo struct {
	ID        uint         `json:"id"`
	Code      string       `json:"code"`
	Summary   string       `json:"summary"`
	Author    string       `json:"author"`
	CreatedAt time.Time    `json:"created_at"`
	Category  CategoryInfo `json:"category"`
}

type ResponsePostDetail struct {
	ID        uint         `json:"id"`
	Code      string       `json:"code"`
	Summary   string       `json:"summary"`
	Content   string       `json:"content"`
	Author    string       `json:"author"`
	CreatedAt time.Time    `json:"created_at"`
	UpdatedAt time.Time    `json:"updated_at"`
	Status    string       `json:"status"`
	Category  CategoryInfo `json:"category"`
}

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

//根据分类查询帖子
func (a *App) queryPostListHandler(w http.ResponseWriter, r *http.Request) {
	resp := &APIStatus{ErrCode: 0}
	defer outputJSON(w, resp)

	if r.Method == "GET" {
		resp.ErrCode = -1
		resp.ErrMessage = "request method is error"
		return
	}

	if r.Body == nil {
		resp.ErrCode = -2
		resp.ErrMessage = "request body is nil"
		return
	}

	data, err := ioutil.ReadAll(ioutil.NopCloser(r.Body))
	if err != nil {
		resp.ErrCode = -3
		resp.ErrMessage = "request body read error"
		return
	}
	defer r.Body.Close()

	req := struct {
		Author     string `json:"author"`
		CategoryID uint   `json:"category_id"`
	}{}

	err = json.Unmarshal(data, &req)
	if err != nil {
		resp.ErrCode = -4
		resp.ErrMessage = "request param unmarshal error"
		return
	}

	var posts []PostListInfo
	if req.Author == "" && req.CategoryID == 0 {
		//查询所有帖子
		posts, err = a.queryAllPostList()
		if err != nil {
			resp.ErrCode = -5
			resp.ErrMessage = "query all posts list error"
			return
		}

		resp.ErrCode = 0
		resp.Data = posts
		return
	}

	if req.Author != "" {
		posts, err = a.queryPostListByAuthor(req.Author)
		if err != nil {
			resp.ErrCode = -6
			resp.ErrMessage = "query posts list by author error"
			return
		}

		resp.ErrCode = 0
		resp.Data = posts
		return
	}

	if req.CategoryID != 0 {
		posts, err = a.queryPostListByCategoryID(req.CategoryID)
		if err != nil {
			resp.ErrCode = -7
			resp.ErrMessage = "query posts list by categoryID error"
			return
		}

		resp.ErrCode = 0
		resp.Data = posts
		return
	}

	return
}

//查询帖子详情
func (a *App) queryPostDetailHandler(w http.ResponseWriter, r *http.Request) {
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

	data, err := ioutil.ReadAll(ioutil.NopCloser(r.Body))
	if err != nil {
		resp.ErrCode = -3
		resp.ErrMessage = "read request param error"
		return
	}
	defer r.Body.Close()

	req := struct {
		ID uint `json:"id"`
	}{}

	err = json.Unmarshal(data, &req)
	if err != nil {
		resp.ErrCode = -4
		resp.ErrMessage = "request param unmarshal error"
		return
	}

	if req.ID == 0 {
		resp.ErrCode = -5
		resp.ErrMessage = "request param is error"
		return
	}

	postDetail := &ResponsePostDetail{ID: req.ID}
	err = a.queryPostDetail(postDetail)
	if err != nil {
		resp.ErrCode = -6
		resp.ErrMessage = "query post detail error"
		return
	}

	resp.ErrCode = 0
	resp.Data = postDetail
	return

}
