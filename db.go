package main

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/mattn/go-sqlite3"
)

func (a *App) initDB() {
	var sqlStmts = []string{
		` CREATE TABLE IF NOT EXISTS posts(
			id INTEGER primary key AUTOINCREMENT,
			code TEXT,
			title TEXT,
			summary TEXT,
			content TEXT,
			author TEXT,
			parent_category INTEGER,
			status TEXT,
			sub_category INTEGER,
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
			id INTEGER primary key AUTOINCREMENT,
			name TEXT,
			status INTEGER,
			encrypted_password TEXT
		)
		`,
		"CREATE INDEX IF NOT EXISTS idx_posts_code ON posts(code);",
	}

	var err error
	app.db, err = sql.Open("sqlite3", fmt.Sprintf("%s%s?cache=shared&mode=rwc", cfg.DataFolder, cfg.SqliteDB))
	if err != nil {
		log.Fatal(err)
	}
	for _, sqlStmt := range sqlStmts {
		_, err := app.db.Exec(sqlStmt)
		if err != nil {
			log.Printf("[ERR] - [InitDB] %q: %s\n", err, sqlStmt)
			return
		}
	}
}
