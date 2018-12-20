package main

import (
	"database/sql"
	"net/http"
	"sync"
	"time"
)

// App blog app
type App struct {
	db         *sql.DB
	mutex      *sync.RWMutex
	HTTPServer *http.Server
	config     *Config
}

// NewApp 生成app
func NewApp(config *Config) *App {
	return &App{
		config: config,
		mutex:  &sync.RWMutex{},
		HTTPServer: &http.Server{
			Addr:         config.Port,
			ReadTimeout:  5 * time.Second,
			WriteTimeout: 10 * time.Second,
		},
	}
}
