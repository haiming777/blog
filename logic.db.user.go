package main

import (
	"database/sql"
	"fmt"
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

func (a *App) queryUser(u User) (result bool, err error) {
	result = true

	a.mutex.Lock()
	defer a.mutex.Unlock()
	db := a.getDB()

	var rows *sql.Rows
	if u.ID != 0 {
		rows, err = db.Query(fmt.Sprintf("SELECT id users WHERE id=%d", u.ID))
	} else if u.Name != "" {
		rows, err = db.Query(fmt.Sprintf("SELECT id users WHERE name=%s", u.Name))
	}

	if err != nil {
		return result, err
	}

	for rows.Next() {
		var id int
		err = rows.Scan(&id)
		if err != nil {
			result = false
			return result, err
		}

		if id == 0 {
			result = false
			return result, err
		}
	}

	return result, err
}

func (a *App) updateUser(u User) error {
	a.mutex.Lock()
	defer a.mutex.Unlock()
	db := a.getDB()

	sqlStr := "UPDATE users"
	updateStr := " SET "
	conditionStr := " WHERE "

	if u.ID != 0 {
		conditionStr = fmt.Sprintf("%sid=%d", conditionStr, u.ID)
		if u.Name != "" {
			updateStr = fmt.Sprintf("%sname=%s", updateStr, u.Name)
		}

		if u.EncryptedPassword != "" {
			updateStr = fmt.Sprintf("%s,encrypted_password=%s", updateStr, u.EncryptedPassword)
		}

		// if u.Status != 0 {
		// 	updateStr = fmt.Sprintf("%s,status=%d", updateStr, u.Status)
		// }
	} else if u.Name != "" {
		conditionStr = fmt.Sprintf("%sname=%s", conditionStr, u.Name)
		if u.EncryptedPassword != "" {
			updateStr = fmt.Sprintf("%sname=%s", updateStr, u.Name)

		}

		// if u.Status != 0 {
		// 	updateStr = fmt.Sprintf("%s,status=%d", updateStr, u.Status)
		// }
	}

	if updateStr == " SET " || conditionStr == " WHERE " {
		return fmt.Errorf("update info error or update condition info error")
	}

	sqlStr = fmt.Sprintf("%s%s%s", sqlStr, updateStr, conditionStr)

	_, err := db.Exec(sqlStr)
	if err != nil {
		return err
	}

	return nil
}
