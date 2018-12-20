package main

import (
	"database/sql"
)

// createUser 创建用户
func (a *App) createUser(u *User) (err error) {
	a.mutex.Lock()
	defer a.mutex.Unlock()

	db := a.getDB()

	var r sql.Result
	r, err = db.Exec("INSERT INTO users(name, encrypted_password, status) values(?,?,?)", u.Name, u.EncryptedPassword, u.Status)
	if err != nil {
		return
	}

	var id int64
	id, err = r.LastInsertId()
	if err != nil {
		return nil
	}
	u.ID = uint(id)

	return
}
