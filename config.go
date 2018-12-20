package main

//Config - env configuration
type Config struct {
	Name       string `json:"name,omitempty"`
	Port       string `json:"port,omitempty"`
	SqliteDB   string `json:"sqlite-db,omitempty"`
	DataFolder string `json:"data,omitempty"`
}

var cfg = Config{
	Name:       "BLOG",
	Port:       ":8080",
	SqliteDB:   "blog.db",
	DataFolder: "./data/",
}
