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

func (a *App) queryUser(u User) (result bool, err error) {
	result = false
	if u.ID == 0 && u.Name == "" {
		return result, fmt.Errorf("query condition error")
	}

	db := a.getDB()
	var uid uint

	if u.ID == 0 {
		err = db.QueryRow(fmt.Sprintf("SELECT id FROM users WHERE name='%s'", u.Name)).Scan(&uid)

	} else {
		err = db.QueryRow(fmt.Sprintf("SELECT id FROM users WHERE id=%d", u.ID)).Scan(&uid)
	}

	switch {
	case err == sql.ErrNoRows:
		return result, nil
	case err != nil:
		return result, err
	default:
		result = true
	}

	return result, err
}

func (a *App) updateUser(u *User) error {
	var err error
	if u.ID == 0 && u.Name == "" {
		err = fmt.Errorf("update user condition error")
		return err
	}
	if u.ID != 0 && u.Name != "" {
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
	stmt, err := db.Prepare("UPDATE users SET password=? WHERE id=?")

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
	stmt, err := db.Prepare("UPDATE users SET password=? WHERE name=?")

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
