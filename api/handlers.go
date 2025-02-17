package main

import (
	"encoding/json"
	"github.com/julienschmidt/httprouter"
	"io/ioutil"
	"log"
	"net/http"
	"video_server/api/dbops"
	"video_server/api/defs"
	"video_server/api/session"
	"video_server/api/utils"
)

func CreateUser(w http.ResponseWriter, r *http.Request, p httprouter.Params)  {
	res, _ := ioutil.ReadAll(r.Body)
	ubody := &defs.UserCredential{}
	if err := json.Unmarshal(res, ubody); err != nil{
		sendErrorResponse(w, defs.ErrorRequestBodyParseFailed)
	}
	if err := dbops.AddUserCredential(ubody.Username, ubody.Pwd); err != nil{
		sendErrorResponse(w, defs.ErrorDBError)
	}
	id := session.GenerateNewSessionId(ubody.Username)
	su := &defs.SignedUp{Success:true, SessionId:id}
	if resp, err := json.Marshal(su); err != nil{
		sendErrorResponse(w, defs.ErrorInternalFaults)
		return
	} else {
		sendNormalResponse(w, string(resp), 201)
	}
}


func GetUserInfo(w http.ResponseWriter, r *http.Request, p httprouter.Params)  {
	if !ValidateUser(w, r){
		log.Printf("Unathorized user \n")
		return
	}
	username := p.ByName("username")
	user, err := dbops.GetUser(username)
	if err != nil{
		log.Printf("Error in GetUserInfo: %v", err)
		sendErrorResponse(w, defs.ErrorDBError)
		return
	}
	ui := &defs.User{Id:user.Id}
	if resp, err := json.Marshal(ui); err != nil{
		sendErrorResponse(w, defs.ErrorInternalFaults)
	} else {
		sendNormalResponse(w, string(resp), http.StatusOK)
	}

}


func Login(w http.ResponseWriter, r *http.Request, p httprouter.Params)  {
	res, _ := ioutil.ReadAll(r.Body)
	log.Printf("Login.request body: %s", res)
	ubody := &defs.UserCredential{}
	if err := json.Unmarshal(res, ubody); err != nil{
		log.Printf("%s", err)
		sendErrorResponse(w, defs.ErrorRequestBodyParseFailed)
	}
	username := p.ByName("username")
	log.Printf("Login url name: %s", username)
	log.Printf("Login body name: %s", ubody.Username)
	if username != ubody.Username{
		sendErrorResponse(w, defs.ErrorNotAuthUser)
		return
	}
	log.Printf("%s", ubody.Username)
	pwd, err := dbops.GetUserCredential(ubody.Username)
	log.Printf("Login pwd: %s", pwd)
	log.Printf("Login body pwd: %s", ubody.Pwd)
	if err != nil || len(pwd) == 0 || pwd != ubody.Pwd{
		sendErrorResponse(w, defs.ErrorNotAuthUser)
		return
	}
	id := session.GenerateNewSessionId(username)
	si := &defs.SignedIn{Success:true, SessionId:id}
	if resp, err := json.Marshal(si); err != nil{
		sendErrorResponse(w, defs.ErrorInternalFaults)
	} else {
		sendNormalResponse(w, string(resp), http.StatusOK)
	}

}

func AddNewVideo(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	if !ValidateUser(w, r){
		log.Printf("Unathorized user\n")
		return
	}
	res, _ := ioutil.ReadAll(r.Body)
	nvbody := &defs.NewVideo{}
	if err := json.Unmarshal(res, nvbody); err != nil{
		log.Printf("Error in AddNewVideo: %v", err)
		sendErrorResponse(w, defs.ErrorRequestBodyParseFailed)
		return
	}
	vi, err := dbops.AddNewVideo(nvbody.AuthorId, nvbody.Name)
	log.Printf("Author id: %d, name: %s \n", nvbody.AuthorId, nvbody.Name)
	if err != nil{
		log.Printf("Error in AddNewVideo: %v", err)
		sendErrorResponse(w, defs.ErrorDBError)
		return
	}
	if resp, err := json.Marshal(vi); err != nil{
		sendErrorResponse(w, defs.ErrorInternalFaults)
	} else {
		sendNormalResponse(w, string(resp), http.StatusCreated)
	}
}

func ListAllVideos(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	if !ValidateUser(w, r){
		log.Printf("Unathorized user\n")
		return
	}
	username := p.ByName("username")
	vs, err := dbops.ListVideoInfo(username, 0, utils.GetCurrentTimestampSec())
	if err != nil{
		log.Printf("Error in ListAllVideos: %v", err)
		sendErrorResponse(w, defs.ErrorDBError)
		return
	}
	vsi := &defs.VideosInfo{Videos:vs}
	if resp, err := json.Marshal(vsi); err != nil {
		log.Printf("Error in ListAllVideos: %v", err)
		sendErrorResponse(w, defs.ErrorInternalFaults)
		return
	} else {
		sendNormalResponse(w, string(resp), http.StatusOK)
	}
}

func DeleteVideo(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	if !ValidateUser(w, r){
		log.Printf("Unathorized user\n")
		return
	}
	vid := p.ByName("vid-id")
	err := dbops.DeleteVideoInfo(vid)
	if err != nil{
		log.Printf("Error in DeleteVideo: %v", err)
		sendErrorResponse(w, defs.ErrorDBError)
		return
	}
	go utils.SendDeleteVideoRequest(vid)
	sendNormalResponse(w, "success", http.StatusNoContent)
}

func PostComment(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	if !ValidateUser(w, r){
		log.Printf("Unathorized user\n")
		return
	}
	reqBody, _ := ioutil.ReadAll(r.Body)
	cbody := &defs.NewComment{}
	if err := json.Unmarshal(reqBody, cbody); err != nil{
		log.Printf("Error in PostComment: %v", err)
		sendErrorResponse(w, defs.ErrorRequestBodyParseFailed)
		return
	}
	vid := p.ByName("vid-id")
	if err := dbops.AddNewComments(vid, cbody.AuthorId, cbody.Content); err != nil{
		log.Printf("Error in PostComment db error: %v", err)
		sendErrorResponse(w, defs.ErrorDBError)
		return
	}
	sendNormalResponse(w, "success", http.StatusCreated)
}

func ShowComments(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	if !ValidateUser(w, r){
		log.Printf("Unathorized user\n")
		return
	}
	vid := p.ByName("vid-id")
	cm, err := dbops.ListComments(vid, 0, utils.GetCurrentTimestampSec())
	if err != nil{
		log.Printf("Error in ShowComments: %s", err)
		sendErrorResponse(w, defs.ErrorDBError)
		return
	}
	cms := &defs.Comments{Comments:cm}
	if resp, err := json.Marshal(cms); err != nil{
		log.Printf("Error in ShowComments: %v", err)
		sendErrorResponse(w, defs.ErrorInternalFaults)
	} else {
		sendNormalResponse(w, string(resp), http.StatusOK)
	}
}