package main

import (
	"github.com/julienschmidt/httprouter"
	"net/http"
)

func RegisterHandlers() *httprouter.Router {
	router := httprouter.New()
	router.POST("/user", CreateUser)
	router.POST("/user/:user_name", Login)
	return router
}

func main()  {
	println("start video server...")
	http.ListenAndServe(":8001", RegisterHandlers())
}


// handler->validation{1.request, 2.user}->business logic->response
// 1. data model
// 2. error handling.