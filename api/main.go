package main

import (
	"github.com/julienschmidt/httprouter"
	"net/http"
	"video_server/api/session"
)

type middlewareHandler struct {
	r *httprouter.Router
}

func NewMiddleWareHanler(r *httprouter.Router) http.Handler {
	m := middlewareHandler{}
	m.r = r
	return m
}

func (m middlewareHandler) ServeHTTP(w http.ResponseWriter, r *http.Request)  {
	// check session
	validateUserSession(r)
	m.r.ServeHTTP(w, r)
}

func RegisterHandlers() *httprouter.Router {
	router := httprouter.New()
	router.POST("/user", CreateUser)
	router.POST("/user/:username", Login)
	router.GET("/user/:username", GetUserInfo)
	router.POST("/user/:username/videos", AddNewVideo)
	router.GET("/user/:username/videos", ListAllVideos)
	router.DELETE("/user/:username/videos/:vid-id", DeleteVideo)
	router.POST("/videos/:vid-id/comments", PostComment)
	router.GET("/videos/:vid-id/comments", ShowComments)
	router.GET("/videos/:vid-id/comments", ShowComments)
	return router
}

func Prepare() {
	session.LoadSessionsFromDB()
}

func main()  {
	println("start video server: http://127.0.0.1:8001...")
	Prepare()
	r := RegisterHandlers()
	mwh := NewMiddleWareHanler(r)
	http.ListenAndServe(":8001", mwh)
}


// handler->validation{1.request, 2.user}->business logic->response
// 1. data model
// 2. error handling.

/*
main->middleware->defs(message, err)->handlers->dbops->response
 */