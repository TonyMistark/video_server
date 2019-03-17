package main

import (
	"net/http"

	"github.com/julienschmidt/httprouter"
)

func RegisterHandlers() *httprouter.Router {
	router := httprouter.New()

	router.GET("/videos/:vid-vid", streamHandler)
	router.POST("/upload/:vie-id", uploadHandler)
	return router
}

func main() {
	r := RegisterHandlers()
	http.ListenAndServe(":9000", r)
}
