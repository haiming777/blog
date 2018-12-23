package main

func (a *App) createCategory(c *Category) error {
	a.mutex.Lock()
	defer a.mutex.Unlock()
	db := a.getDB()

	stmt, err := db.Prepare("INSERT INTO categories (name,parent_id) VALUES (?,?)")
	if err != nil {
		return err
	}

	r, err := stmt.Exec(c.Name, c.ParentID)
	if err != nil {
		return err
	}

	cid, err := r.LastInsertId()
	if err != nil {
		return nil
	}

	c.ID = uint(cid)
	return nil
}
