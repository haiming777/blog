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

// queryUserByName 根据用户名查询用户信息
func (a *App) queryUserByName(name string, user *User) error {
	db := a.getDB()

	a.mutex.RLock()
	defer a.mutex.RUnlock()

	return db.QueryRow("SELECT id,name,status FROM users WHERE name=? ", name).
		Scan(&user.ID, &user.Name, &user.Status)
}

func (a *App) updateUserPassword(u *User) error {
	a.mutex.Lock()
	defer a.mutex.Unlock()

	db := a.getDB()

	_, err := db.Exec("UPDATE users SET encrypted_password=? WHERE name=?", u.Name, u.EncryptedPassword)
	if err != nil {
		return err
	}

	return nil
}

// func (a *App) updateUserPasswordWithID(u *User) error {
// 	a.mutex.Lock()
// 	defer a.mutex.Unlock()

// 	db := a.getDB()

// 	_, err := db.Exec("UPDATE users SET encrypted_password=? WHERE id=?", u.EncryptedPassword, u.ID)
// 	if err != nil {
// 		return err
// 	}

// 	return nil
// }

// func (a *App) updatePasswordWithName(u *User) error {
// 	a.mutex.Lock()
// 	defer a.mutex.Unlock()

// 	db := a.getDB()

// 	r, err := db.Exec("UPDATE users SET encrypted_password=? WHERE name=?", u.EncryptedPassword, u.Name)
// 	if err != nil {
// 		return err
// 	}

// 	uid, err := r.RowsAffected()
// 	if err != nil {
// 		return err
// 	}
// 	u.ID = uint(uid)

// 	return nil
// }

// func (a *App) updateStatusWithID(u *User) error {
// 	a.mutex.Lock()
// 	defer a.mutex.Unlock()

// 	db := a.getDB()

// 	_, err := db.Exec("UPDATE users SET status=? WHERE id=?", u.Status, u.ID)
// 	if err != nil {
// 		return err
// 	}

// 	return nil
// }

// func (a *App) updateStatusWithName(u *User) error {
// 	a.mutex.Lock()
// 	defer a.mutex.Unlock()

// 	db := a.getDB()

// 	r, err := db.Exec("UPDATE users SET status=? WHERE name=?", u.Status, u.Name)
// 	if err != nil {
// 		return err
// 	}

// 	uid, err := r.RowsAffected()
// 	if err != nil {
// 		return err
// 	}
// 	u.ID = uint(uid)
// 	return nil
// }
