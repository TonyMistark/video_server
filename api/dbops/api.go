package dbops

import (
	"database/sql"
	"log"
	_ "github.com/go-sql-driver/mysql"
	_ "database/sql"
	"time"
	"video_server/api/defs"
	"video_server/api/utils"
)


func AddUserCredential(loginName string, pwd string) error {
		stmtIns, err := dbConn.Prepare("insert into users (login_name, pwd) values (?, ?)")
		if err != nil{
			return err
		}
		stmtIns.Exec(loginName, pwd)
		stmtIns.Close()
		return nil
}

func GetUserCredential(loginName string) (string, error)  {
	stmtOut, err := dbConn.Prepare("select pwd from users where login_name = ?")
	if err != nil{
		log.Printf( "%s", err)
		return "", err
	}
	var pwd string
	stmtOut.QueryRow(loginName).Scan(&pwd)
	stmtOut.Close()
	return pwd, nil
}

func GetUser(username string) (*defs.User, error) {
	stmtOut, err := dbConn.Prepare("select id, pwd from users where login_name = ?")
	if err != nil{
		log.Printf("GetUser: %v", err)
		return nil, err
	}
	var id, pwd string
	err = stmtOut.QueryRow(username).Scan(&id, &pwd)
	if err != nil && err != sql.ErrNoRows{
		log.Printf("Query User error: %v", err)
		return nil, err
	}
	if err == sql.ErrNoRows{
		log.Printf("GetUser: not exists username: %s", username)
	}
	res := &defs.User{Id:id, LoginName:username, Pwd:pwd}
	defer stmtOut.Close()
	return res, nil


}

func DeleteUser(loginName string, pwd string) error {
	stmtDel, err := dbConn.Prepare("DELETE from users where login_name=? and pwd = ?")
	if err != nil{
		log.Printf("%s", err)
		return err
	}
	_, err = stmtDel.Exec(loginName, pwd)
	if err != nil{
		return err
	}
	defer stmtDel.Close()
	return nil
}

func AddNewVideo(aid int, name string) (*defs.VideoInfo, error) {
	// create uuid
	vid, err := utils.NewUUID()
	if err != nil{
		return nil, err
	}
	t := time.Now()
	ctime := t.Format("Jan 02 2006, 15:04:05") // M D y, HH:MM:SS
	stmtIns, err := dbConn.Prepare(`insert into video_info (id, author_id, name, display_ctime) values (?, ?, ?, ?)`)
	if err != nil{
		return nil, err
	}
	_, err = stmtIns.Exec(vid, aid, name, ctime)
	if err != nil{
		return nil, err
	}
	res := &defs.VideoInfo{Id: vid, AuthorId: aid, Name: name, DisplayCtime: ctime}
	defer stmtIns.Close()
	return res, nil
}


func GetVideoInfo(vid string) (*defs.VideoInfo, error) {
	stmtOut, err := dbConn.Prepare("SELECT author_id, name, display_ctime FROM video_info WHERE id=?")

	var aid int
	var dct string
	var name string

	err = stmtOut.QueryRow(vid).Scan(&aid, &name, &dct)
	if err != nil && err != sql.ErrNoRows{
		return nil, err
	}

	if err == sql.ErrNoRows {
		return nil, nil
	}

	defer stmtOut.Close()

	res := &defs.VideoInfo{Id: vid, AuthorId: aid, Name: name, DisplayCtime: dct}

	return res, nil
}


func ListVideoInfo(uname string, from, to int) ([]*defs.VideoInfo, error) {
	stmtOut, err := dbConn.Prepare(`SELECT video_info.id, video_info.author_id, video_info.name, video_info.display_ctime FROM video_info 
		INNER JOIN users ON video_info.author_id = users.id
		WHERE users.login_name = ? AND video_info.create_time > FROM_UNIXTIME(?) AND video_info.create_time <= FROM_UNIXTIME(?) 
		ORDER BY video_info.create_time DESC`)


	var res []*defs.VideoInfo

	if err != nil {
		return res, err
	}

	rows, err := stmtOut.Query(uname, from, to)
	if err != nil {
		log.Printf("%s", err)
		return res, err
	}

	for rows.Next() {
		var id, name, ctime string
		var aid int
		if err := rows.Scan(&id, &aid, &name, &ctime); err != nil {
			return res, err
		}

		vi := &defs.VideoInfo{Id: id, AuthorId: aid, Name: name, DisplayCtime: ctime}
		res = append(res, vi)
	}

	defer stmtOut.Close()

	return res, nil
}

func DeleteVideoInfo(vid string) error {
	stmtDel, err := dbConn.Prepare("DELETE FROM video_info WHERE id=?")
	if err != nil {
		return err
	}

	_, err = stmtDel.Exec(vid)
	if err != nil {
		return err
	}

	defer stmtDel.Close()
	return nil
}

func AddNewComments(vid string, aid int, content string) error {
	id, err := utils.NewUUID()
	if err != nil{
		return err
	}
	stmtIns, err := dbConn.Prepare("insert into comments (id, video_id, author_id, content) values (?, ?, ?, ?)")
	if err != nil{
		return err
	}
	_, err = stmtIns.Exec(id, vid, aid, content)
	if err != nil {
		return err
	}
	defer stmtIns.Close()
	return nil
}

func ListComments(vid string, limit, offset int) ([]*defs.Comment, error) {
	stmtOut, err := dbConn.Prepare(`select 
											 comments.id, users.login_name, comments.content 
										   from 
										  	  comments
										   inner join 
											  users 
										   on 
											  comments.author_id = users.id 
										   where 
										  	  video_id='123456' limit ? offset ?
										   order by 
											   comments.time 
											DESC`)
	var res []*defs.Comment
	rows, err := stmtOut.Query(limit, offset)
	if err != nil{
		return res, err
	}
	for rows.Next() {
		var id, name, content string
		if err := rows.Scan(&id, &name, &content); err != nil{
			return res, err
		}
		c := &defs.Comment{Id: id, VideoId: vid, Author: name, Content: content}
		res = append(res, c)
	}
	return res, nil
}