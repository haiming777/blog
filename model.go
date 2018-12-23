package main

import (
	"time"
)

//APIStatus - general api result
type APIStatus struct {
	ErrCode    int         `json:"code"`
	ErrMessage string      `json:"msg,omitempty"`
	Data       interface{} `json:"response,omitempty"`
}

// Post 文章
type Post struct {
	ID             uint      `json:"id,omitempty"`
	Code           string    `json:"code,omitempty"`
	Summary        string    `json:"summary,omitempty"`
	Content        string    `json:"content,omitempty"`
	Author         string    `json:"author,omitempty"`
	ParentCategory int       `json:"parent_category,omitempty"`
	SubCategory    int       `json:"sub_category,omitempty"`
	UpdatedAt      time.Time `json:"updated_at,omitempty"`
	CreatedAt      time.Time `json:"created_at,omitempty"`
	Status         string    `json:"status,omitempty"`
}

// Category 文章分类
type Category struct {
	ID       uint   `json:"id,omitempty"`
	Name     string `json:"name,omitempty"`
	ParentID uint   `json:"parent_id,omitempty"`
}

// User 用户
type User struct {
	ID                uint   `json:"id,omitempty"`
	Name              string `json:"name,omitempty"`
	EncryptedPassword string `json:"encrypted_password,omitempty"`
	Status            uint   `json:"status,omitempty"`
}
