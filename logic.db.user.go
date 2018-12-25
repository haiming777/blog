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

//queryUser 查询用户，返回结果：user返回对外用户信息;err错误信息反馈；
func (a *App) queryUser(u User) (user *User, err error) {
	if u.ID == 0 && u.Name == "" {
		err = fmt.Errorf("query condition error")
		return
	}

	if u.ID != 0 {
		user, err = a.queryUserWithID(u)
	} else {
		user, err = a.queryUserWithName(u)
	}

	if err != nil && err == sql.ErrNoRows {
		err = nil
		return
	}

	return
}

func (a *App) queryUserWithID(u User) (user *User, err error) {
	if u.ID == 0 && u.Name == "" {
		err = fmt.Errorf("query condition error")
		return
	}

	user = &User{}
	db := a.getDB()

	var uid, status uint
	var name string

	err = db.QueryRow("SELECT id,name,status FROM users WHERE id=? ", u.ID).Scan(&uid, &name, &status)
	if err == nil {
		user.ID = uid
		user.Name = name
		user.Status = status
	}

	return
}

func (a *App) queryUserWithName(u User) (user *User, err error) {
	if u.ID == 0 && u.Name == "" {
		err = fmt.Errorf("query condition error")
		return
	}

	user = &User{}
	db := a.getDB()

	var uid, status uint
	var name string
	err = db.QueryRow("SELECT id,name,status FROM users WHERE name=? ", u.Name).Scan(&uid, &name, &status)

	if err == nil {
		user.ID = uid
		user.Name = name
		user.Status = status
	}

	return
}

func (a *App) updateUser(u *User) error {
	var err error
	if u.ID == 0 && u.Name == "" {
		err = fmt.Errorf("update user condition error")
		return err
	}
	if u.ID != 0 && u.Name != "" {
		qu := User{Name: u.Name}
		user, _ := a.queryUser(qu)
		if user.ID != 0 {
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

	_, err := db.Exec("UPDATE users SET name=? WHERE id=?", u.Name, u.ID)
	if err != nil {
		return err
	}

	return nil
}

func (a *App) updateUserPasswordWithID(u *User) error {
	a.mutex.Lock()
	defer a.mutex.Unlock()

	db := a.getDB()

	_, err := db.Exec("UPDATE users SET encrypted_password=? WHERE id=?", u.EncryptedPassword, u.ID)
	if err != nil {
		return err
	}

	return nil
}

func (a *App) updatePasswordWithName(u *User) error {
	a.mutex.Lock()
	defer a.mutex.Unlock()

	db := a.getDB()

	r, err := db.Exec("UPDATE users SET encrypted_password=? WHERE name=?", u.EncryptedPassword, u.Name)
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

	_, err := db.Exec("UPDATE users SET status=? WHERE id=?", u.Status, u.ID)
	if err != nil {
		return err
	}

	return nil
}

func (a *App) updateStatusWithName(u *User) error {
	a.mutex.Lock()
	defer a.mutex.Unlock()

	db := a.getDB()

	r, err := db.Exec("UPDATE users SET status=? WHERE name=?", u.Status, u.Name)
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
