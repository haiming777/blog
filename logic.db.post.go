package main

func (a *App) createPost(p *Post) error {
	db := a.getDB()

	a.mutex.Lock()
	defer a.mutex.Unlock()

	r, err := db.Exec("INSERT INTO posts (code,summary,content,author,parent_category,sub_category,created_at,status) VALUES (?,?,?,?,?,?,?,?)", p.Code, p.Summary, p.Content, p.Author, p.ParentCategory, p.SubCategory, p.CreatedAt, p.Status)
	if err != nil {
		return err
	}

	pid, err := r.LastInsertId()
	if err != nil {
		return nil
	}

	p.ID = uint(pid)
	return nil

}
