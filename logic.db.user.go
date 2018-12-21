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
	smt, err := db.Prepare("INSERT INTO users(name, encrypted_password, status) values(?,?,?)")
	if err != nil {
		return
	}

	r, err = smt.Exec(u.Name, u.EncryptedPassword, u.Status)
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

//queryUser 查询用户，返回结果：bool操作用户信息判断；err错误信息反馈；user返回对外用户信息
func (a *App) queryUser(u User) (result bool, user *User, err error) {
	result = false
	if u.ID == 0 && u.Name == "" {
		err = fmt.Errorf("query condition error")
		return result, nil, err
	}

	db := a.getDB()
	var name string
	var encryptedPassword string
	var status uint

	if u.ID == 0 {
		err = db.QueryRow(fmt.Sprintf("SELECT name,encrypted_password,status FROM users WHERE name='%s'", u.Name)).Scan(&name, &encryptedPassword, &status)

	} else {
		err = db.QueryRow(fmt.Sprintf("SELECT name,encrypted_password,status FROM users WHERE id=%d", u.ID)).Scan(&name, &encryptedPassword, &status)
	}

	switch {
	case err == sql.ErrNoRows:
		return result, nil, nil
	case err != nil:
		return result, nil, err
	default:
		result = true
		user = &User{
			Name:              name,
			EncryptedPassword: encryptedPassword,
			Status:            status,
		}
	}

	return result, user, err
}

func (a *App) updateUser(u *User) error {
	var err error
	if u.ID == 0 && u.Name == "" {
		err = fmt.Errorf("update user condition error")
		return err
	}
	if u.ID != 0 && u.Name != "" {
		qu := User{Name: u.Name}
		result, _, _ := a.queryUser(qu)
		if result {
			err = fmt.Errorf("Name already exist")
			return err
		}
		err = a.updateUserName(u)
	}

	if u.EncryptedPassword != "" {

		if u.ID != 0 {
			err = a.updateUserPasswordWithID(u)
		} else {
			if u.Name != "" {
				err = a.updatePasswordWithName(u)
			}
		}
	}

	if u.Status != 0 {
		if u.ID != 0 {
			err = a.updateStatusWithID(u)
		} else {
			if u.Name != "" {
				err = a.updateStatusWithName(u)
			}
		}
	}

	return nil
}

func (a *App) updateUserName(u *User) error {
	a.mutex.Lock()
	defer a.mutex.Unlock()

	db := a.getDB()
	stmt, err := db.Prepare("UPDATE users SET name=? WHERE id=?")

	if err != nil {
		return err
	}

	r, err := stmt.Exec(u.Name, u.ID)
	if err != nil {
		return err
	}
	_, err = r.RowsAffected()
	if err != nil {
		return err
	}

	return nil
}

func (a *App) updateUserPasswordWithID(u *User) error {
	a.mutex.Lock()
	defer a.mutex.Unlock()

	db := a.getDB()
	stmt, err := db.Prepare("UPDATE users SET encrypted_password=? WHERE id=?")

	if err != nil {
		return err
	}

	r, err := stmt.Exec(u.EncryptedPassword, u.ID)
	if err != nil {
		return err
	}
	_, err = r.RowsAffected()
	if err != nil {
		return err
	}

	return nil
}

func (a *App) updatePasswordWithName(u *User) error {
	a.mutex.Lock()
	defer a.mutex.Unlock()

	db := a.getDB()
	stmt, err := db.Prepare("UPDATE users SET encrypted_password=? WHERE name=?")

	if err != nil {
		return err
	}

	r, err := stmt.Exec(u.EncryptedPassword, u.Name)
	if err != nil {
		return err
	}
	uid, err := r.RowsAffected()
	if err != nil {
		return err
	}
	u.ID = uint(uid)

	return nil
}

func (a *App) updateStatusWithID(u *User) error {
	a.mutex.Lock()
	defer a.mutex.Unlock()

	db := a.getDB()
	stmt, err := db.Prepare("UPDATE users SET status=? WHERE id=?")

	if err != nil {
		return err
	}

	r, err := stmt.Exec(u.Status, u.ID)
	if err != nil {
		return err
	}
	_, err = r.RowsAffected()
	if err != nil {
		return err
	}

	return nil
}

func (a *App) updateStatusWithName(u *User) error {
	a.mutex.Lock()
	defer a.mutex.Unlock()

	db := a.getDB()
	stmt, err := db.Prepare("UPDATE users SET status=? WHERE name=?")

	if err != nil {
		return err
	}

	r, err := stmt.Exec(u.Status, u.Name)
	if err != nil {
		return err
	}
	uid, err := r.RowsAffected()
	if err != nil {
		return err
	}
	u.ID = uint(uid)
	return nil
}
