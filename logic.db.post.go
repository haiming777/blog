package main

import (
	"time"
)

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

func (a *App) queryAllPostList() ([]PostListInfo, error) {
	db := a.getDB()

	a.mutex.RLock()
	defer a.mutex.RUnlock()

	rows, err := db.Query("SELECT posts.id,posts.code,posts.summary,posts.author,posts.created_at,categories.id,categories.name FROM posts INNER JOIN categories WHERE posts.sub_category=categories.id")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var (
		pid          uint
		code         string
		summary      string
		author       string
		categoryID   uint
		categoryName string
		createdAt    time.Time
	)

	postList := make([]PostListInfo, 0)
	for rows.Next() {
		err = rows.Scan(&pid, &code, &summary, &author, &createdAt, &categoryID, &categoryName)
		if err != nil {
			return nil, err
		}

		post := PostListInfo{
			ID:        pid,
			Code:      code,
			Summary:   summary,
			Author:    author,
			CreatedAt: createdAt,
			Category:  CategoryInfo{ID: categoryID, Name: categoryName},
		}
		postList = append(postList, post)
	}

	return postList, nil

}

func (a *App) queryPostListByAuthor(author string) ([]PostListInfo, error) {
	db := a.getDB()

	a.mutex.RLock()
	defer a.mutex.RUnlock()

	rows, err := db.Query("SELECT posts.id,posts.code,posts.summary,posts.created_at,categories.id,categories.name FROM posts INNER JOIN categories WHERE posts.sub_category=categories.id AND posts.author=?", author)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var (
		pid          uint
		code         string
		summary      string
		categoryID   uint
		categoryName string
		createdAt    time.Time
	)

	postList := make([]PostListInfo, 0)
	for rows.Next() {
		err = rows.Scan(&pid, &code, &summary, &createdAt, &categoryID, &categoryName)
		if err != nil {
			return nil, err
		}

		post := PostListInfo{
			ID:        pid,
			Code:      code,
			Summary:   summary,
			Author:    author,
			CreatedAt: createdAt,
			Category:  CategoryInfo{ID: categoryID, Name: categoryName},
		}
		postList = append(postList, post)
	}

	return postList, nil

}

func (a *App) queryPostListByCategoryID(cid uint) ([]PostListInfo, error) {
	db := a.getDB()

	a.mutex.RLock()
	defer a.mutex.RUnlock()

	rows, err := db.Query("SELECT posts.id,posts.code,posts.summary,posts.author,posts.created_at,categories.id,categories.name FROM posts INNER JOIN categories WHERE posts.sub_category=categories.id AND posts.sub_category=?", cid)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var (
		pid          uint
		code         string
		summary      string
		author       string
		categoryID   uint
		categoryName string
		createdAt    time.Time
	)

	postList := make([]PostListInfo, 0)
	for rows.Next() {
		err = rows.Scan(&pid, &code, &summary, &createdAt, &categoryID, &categoryName)
		if err != nil {
			return nil, err
		}

		post := PostListInfo{
			ID:        pid,
			Code:      code,
			Summary:   summary,
			Author:    author,
			CreatedAt: createdAt,
			Category:  CategoryInfo{ID: categoryID, Name: categoryName},
		}
		postList = append(postList, post)
	}

	return postList, nil

}
