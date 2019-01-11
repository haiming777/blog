package main

func (a *App) createCategory(c *Category) error {
	a.mutex.Lock()
	defer a.mutex.Unlock()
	db := a.getDB()

	r, err := db.Exec("INSERT INTO categories (name,parent_id) VALUES (?,?)", c.Name, c.ParentID)
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

//queryCategories 查询所有分类
func (a *App) queryCategories() ([]Category, error) {
	db := a.getDB()

	a.mutex.RLock()
	defer a.mutex.RUnlock()

	rows, err := db.Query("SELECT * FROM categories")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var (
		cid      uint
		name     string
		parentID uint
	)

	categories := make([]Category, 0)
	for rows.Next() {
		err = rows.Scan(&cid, &name, &parentID)
		if err != nil {
			return nil, err
		}
		category := Category{
			ID:       cid,
			Name:     name,
			ParentID: parentID,
		}
		categories = append(categories, category)
	}

	return categories, nil
}
