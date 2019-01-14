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
func (a *App) queryCategoriesFirstLevel() ([]Category, error) {
	db := a.getDB()

	a.mutex.RLock()
	defer a.mutex.RUnlock()

	rows, err := db.Query("SELECT * FROM categories WHERE parent_id=?", 0)
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

func (a *App) queryCategoryByID(c *Category) error {
	db := a.getDB()

	a.mutex.RLock()
	defer a.mutex.RUnlock()

	return db.QueryRow("SELECT * FROM categories WHERE id=?", c.ID).
		Scan(&c.ID, &c.Name, &c.ParentID)
}

func (a *App) queryCategoriesByParentID(pid uint) ([]Category, error) {
	db := a.getDB()

	a.mutex.RLock()
	defer a.mutex.RUnlock()

	var (
		id   uint
		name string
	)

	rows, err := db.Query("SELECT id,name FROM categories WHERE parent_id=?", pid)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	categories := make([]Category, 0)
	for rows.Next() {

		err = rows.Scan(&id, &name)
		if err != nil {
			return nil, err
		}

		category := Category{ID: id, Name: name, ParentID: pid}
		categories = append(categories, category)
	}

	return categories, nil
}
