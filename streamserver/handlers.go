package main

import (
	"github.com/gpmgo/gopm/modules/log"
	"github.com/julienschmidt/httprouter"
	"html/template"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"time"
)

func streamHandler(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	vid := p.ByName("vid-id")
	vl := VIDEO_DIR + vid
	println("vl: ", vl)
	video, err := os.Open(vl)
	if err != nil{
		log.Print(0, "Error when try to open file: %v", err)
		sendErrorResponse(w, http.StatusInternalServerError, "Internal Server Error")
		return
	}
	w.Header().Set("Content-Type", "video/mp4")
	http.ServeContent(w, r, "", time.Now(), video)
	defer video.Close()


	targetUrl := "http:/x/b/c" + p.ByName("vid-id")
	http.Redirect(w, r, targetUrl, 301)
}

func uploadHandler(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	r.Body = http.MaxBytesReader(w, r.Body, MAX_UPLOAD_SIZE)
	if err := r.ParseMultipartForm(MAX_UPLOAD_SIZE); err != nil{
		sendErrorResponse(w, http.StatusBadRequest, "File is too Big")
		return
	}
	file, _, err := r.FormFile("file")	//<form name="file"
	if err != nil{
		sendErrorResponse(w, http.StatusInternalServerError, "Internal server error")
		return
	}
	data, err := ioutil.ReadAll(file)
	if err != nil{
		log.Print(0, "read file error: %v", err)
		sendErrorResponse(w, http.StatusInternalServerError, "Internal server error")
	}

	fn := p.ByName("vid-id")
	err = ioutil.WriteFile(VIDEO_DIR + fn, data, 0666)
	if err != nil{
		log.Print(0, "Write file error: %v", err)
		sendErrorResponse(w, http.StatusInternalServerError, "Internal server error")
		return
	}
	ossfn := "video/" + fn
	path := "./videos/" + fn
	bn := "ice-video"
	ret := UploadToOss(ossfn, path, bn)
	if !ret {
		sendErrorResponse(w, http.StatusInternalServerError, "Internal Service Error")
		return
	}
	os.Remove(path)

	w.WriteHeader(http.StatusCreated)
	io.WriteString(w, "Uploaded successfully")

}

func testPageHandler(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	t, _ := template.ParseFiles("./videos/upload.html")
	t.Execute(w, nil)
}