package main

import (
	"httprouter"
	"net/http"
	"video_server/scheduler/taskrunner"
)

func RegisterHandlers() *httprouter.Router {
	router := httprouter.New()
	router.GET("/video-delete-record/:vid-id", vidoDelRecHandler)
	return router
}


func main() {
	println("start video server: http://127.0.0.1:9002...")
	go taskrunner.Start()
	r := RegisterHandlers()
	http.ListenAndServe(":9002", r)

}

