package main

import "log"

func initDB() {
	var sql = []string{
		` CREATE TABLE IF NOT EXISTS posts(
			code TEXT primary key,
			title TEXT,
			summary TEXT,
			content TEXT,
			author TEXT,
			parent_category TEXT,
			status TEXT,
			sub_category TEXT,
			created_at DATETIME,
			updated_at DATETIME
		)
		`,
		` CREATE TABLE IF NOT EXISTS categories(
			id INTEGER primary key AUTOINCREMENT,
			name TEXT ,
			parent_id INTEGER
		)
		`,
		`CREATE TABLE IF NOT EXISTS users(
			user TEXT primary key,
			password TEXT
		)
		`,
	}

	for _, v := range sql {
		log.Println("[initDB]:", v)
	}
}
