package main

import (
	"log"
	"net/http"
)

func (a *App) initHandler() {
	//user
	a.registerRouter("/user/create", a.createUserHandler)
	a.registerRouter("/user/update", a.updataUserHandler)
	a.registerRouter("/signin", a.signin)

	//category
	//创建分类
	a.registerRouter("/category/create", a.createCategoryHandler)
	//获取一级分类
	a.registerRouter("/categories", a.categoryListHandler)
	//获取子分类
	a.registerRouter("/sub-categories", a.queryCategoriesByParentIDHandler)

	//post
	//创建帖子
	a.registerRouter("/posts/create", a.createPostHandler)
	//查询帖子列表
	a.registerRouter("/posts", a.queryPostListHandler)
	//查询帖子详情
	a.registerRouter("/posts-detail", a.queryPostDetailHandler)

}

func (a *App) registerRouter(pattern string, handler func(http.ResponseWriter, *http.Request)) {
	log.Printf("handler:[%s]\n", pattern)
	http.HandleFunc(pattern, handler)
}

func (a *App) listenAndServer() {
	if err := a.HTTPServer.ListenAndServe(); err != nil {
		log.Fatal("[ERR] - ListenAndServe:", err)
	}
}
