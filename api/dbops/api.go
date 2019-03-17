package dbops

import (
	"database/sql"
	_ "database/sql"
	"log"
	"time"

	"video_server/api/defs"

	"video_server/api/utils"

	_ "github.com/go-sql-driver/mysql"
)

func AddUserCredential(loginNmae string, pwd string) error {
	stmtIns, err := dbConn.Prepare("INSERT INTO user (login_name, pwd) values (?, ?)")
	if err != nil {
		return err
	}
	_, err = stmtIns.Exec(loginNmae, pwd)
	if err != nil {
		return err
	}
	defer stmtIns.Close()
	return nil
}

func GetUserCredential(loginName string) (string, error) {
	stmtOut, err := dbConn.Prepare("SELECT pwd from user WHERE login_name = ?")
	if err != nil {
		log.Printf("%s", err)
		return "", err
	}
	var pwd string
	err = stmtOut.QueryRow(loginName).Scan(&pwd)
	if err != nil && err != sql.ErrNoRows {
		return "", err
	}
	defer stmtOut.Close()
	return pwd, nil
}

func DeleteUser(loginName string, pwd string) error {
	stmtDel, err := dbConn.Prepare("DELETE FROM user WHERE login_name = ? AND pwd = ?")
	if err != nil {
		log.Printf("DeleteUser error: %s", err)
		return err
	}
	_, err = stmtDel.Exec(loginName, pwd)
	if err != nil {
		return err
	}
	defer stmtDel.Close()
	return nil
}

func AddVideo(aid int, name string) (*defs.VideoInfo, error) {
	// create uuid
	vid, err := utils.NewUUID()
	if err != nil {
		return nil, err
	}
	// createtime -> db ->
	t := time.Now()
	ctime := t.Format("Jan 02 2006, 15:04:05") // M D y, HH:MM:SS
	stmtIns, err := dbConn.Prepare(`INSERT INTO video_info (id, author_id, name, display_ctime VALUES (?, ?, ?, ?)) `)
	if err != nil {
		return nil, err
	}
	_, err = stmtIns.Exec(vid, aid, name, ctime)
	if err != nil {
		return nil, err
	}
	res := &defs.VideoInfo{Id: vid, AuthorId: aid, Name: name, DisplayCtime: ctime}
	defer stmtIns.Close()
	return res, nil
}

func GetVideo(vid int) (*defs.VideoInfo, error) {
	stmtOut, err := dbConn.Prepare("SELECT id, author_id, name, display_ctime FROM video_info WHERE id=?")
	if err != nil {
		return nil, err
	}
	var id, display_ctime, name string
	var author_id int
	err = stmtOut.QueryRow(vid).Scan(id, author_id, name, display_ctime)
	if err != nil {
		return nil, err
	}
	res := &defs.VideoInfo{Id: id, AuthorId: author_id, Name: name, DisplayCtime: display_ctime}
	defer stmtOut.Close()
	return res, err
}

func AddComments(vid string, aid int, content string) error {
	id, err := utils.NewUUID()
	if err != nil {
		return err
	}
	stmtIns, err := dbConn.Prepare("INSERT INT comments (id, video_id, author_id, content) VALUES (?, ?, ?, ?)")
	if err != nil {
		return err
	}
	_, err = stmtIns.Exec(id, vid, aid, content)
	if err != nil {
		return err
	}
	defer stmtIns.Close()
	return nil
}

func ListComments(vid string, from, to int) ([]*defs.Comment, error) {
	stmtOut, err := dbConn.Prepare(
		`SELECT comments.id, user.login_name, comments.content FROM comments
				INNER JOIN user on comments.author_id = user.id WHERE comments.video_id = ?
				AND comments.time > FROM_UNIXTIME(?) AND comments.time <= FROM_UNIXTIME(?);`)
	if err != nil {
		return nil, err
	}
	var res []*defs.Comment
	rows, err := stmtOut.Query(vid, from, to)
	if err != nil {
		return nil, err
	}
	for rows.Next() {
		var id, name, content string
		if err := rows.Scan(&id, &name, &content); err != nil {
			return nil, err
		}
		c := &defs.Comment{Id: id, VideoId: vid, Author: name, Content: content}
		res = append(res, c)
	}
	defer stmtOut.Close()
	return res, nil

}
