package main

import (
	"crypto/md5"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
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
		Name     string `json:"name"`
		Password string `json:"password"`
	}{}

	err = json.Unmarshal(data, &req)
	if err != nil {

		outputJSON(w, APIStatus{
			ErrCode:    -1,
			ErrMessage: "json format struct error",
		})
		return
	}

	if req.Name == "" {
		outputJSON(w, APIStatus{
			ErrCode:    -2,
			ErrMessage: "account is empty",
		})
		return
	}

	if req.Password == "" {
		outputJSON(w, APIStatus{
			ErrCode:    -3,
			ErrMessage: "password is empty",
		})
		return
	}

	//查询user name 是否已存在

	exist, err := a.queryUser(User{Name: req.Name})
	if err != nil {
		outputJSON(w, APIStatus{
			ErrCode:    -4,
			ErrMessage: fmt.Sprintf("db query error:%s", err.Error()),
		})
		return
	}

	if exist {
		outputJSON(w, APIStatus{
			ErrCode:    -4,
			ErrMessage: "name is already signup",
		})
		return
	}

	//创建user
	encryptedPassword := fmt.Sprintf("%x", md5.Sum([]byte(req.Password)))

	user := &User{
		Name:              req.Name,
		EncryptedPassword: encryptedPassword,
		Status:            1,
	}

	err = a.createUser(user)
	if err != nil {
		outputJSON(w, APIStatus{
			ErrCode:    -5,
			ErrMessage: fmt.Sprintf("create user error:%s", err.Error()),
		})
		return
	}

	outputJSON(w, APIStatus{
		ErrCode: 0,
	})

}

func (a *App) updataUserHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != "PUT" {
		outputJSON(w, APIStatus{
			ErrCode:    -1,
			ErrMessage: "request method error",
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
		ID       uint   `json:"id"`
		Name     string `json:"name"`
		Password string `json:"password"`
		State    uint   `json:"state"`
	}{}

	err = json.Unmarshal(data, &req)
	if err != nil {

		outputJSON(w, APIStatus{
			ErrCode:    -1,
			ErrMessage: "json format struct error",
		})
		return
	}

	if req.ID == 0 && req.Name == "" {
		outputJSON(w, APIStatus{
			ErrCode:    -1,
			ErrMessage: "param is error",
		})

		return
	}

	exist, err := a.queryUser(User{ID: req.ID, Name: req.Name})
	if err != nil {
		outputJSON(w, APIStatus{
			ErrCode:    -2,
			ErrMessage: fmt.Sprintf("db query error:%s", err.Error()),
		})
		return
	}

	if !exist {
		outputJSON(w, APIStatus{
			ErrCode:    -2,
			ErrMessage: "user name is not exist",
		})
	}

	//修改密码没有有效字符，则不做修改
	targetStr := strings.Replace(req.Password, " ", "", -1)
	encryptedPassword := targetStr
	if targetStr != "" {
		encryptedPassword = fmt.Sprintf("%x", md5.Sum([]byte(req.Password)))
	}

	user := &User{
		ID:                req.ID,
		Name:              req.Name,
		EncryptedPassword: encryptedPassword,
		Status:            req.State,
	}

	err = a.updateUser(user)
	if err != nil {
		outputJSON(w, APIStatus{
			ErrCode:    -3,
			ErrMessage: fmt.Sprintf("update user info error:%s", err.Error()),
		})
		return
	}

	outputJSON(w, APIStatus{
		ErrCode: 0,
	})
}
