package dbops

import (
	"fmt"
	"testing"
)

/*
init(dblogin, truncat tables) -> run tests -> clear data(truncate tables)
*/

var tempvid string

func clearTables()  {
	dbConn.Exec("trucate users")
	dbConn.Exec("trucate video_info")
	dbConn.Exec("trucate comments")
	dbConn.Exec("trucate sessions")
}

func TestMain(m *testing.M)  {
	clearTables()
	m.Run()
	clearTables()
}

func TestUserWorkFlow(t *testing.T)  {
	t.Run("Add", testAddUser)
	t.Run("Get", testGetUser)
	t.Run("Del", testDeleteUser)
	t.Run("Reget", testReGetUser)

}

func testAddUser(t *testing.T)  {
	err := AddUserCredential("ice", "123")
	if err != nil{
		t.Errorf("Error: %v", err)
	}

}

func testGetUser(t *testing.T)  {
	pwd, err := GetUserCredential("ice")
	if pwd !="123" || err != nil{
		t.Errorf("Error of get user")
	}
}

func testDeleteUser(t *testing.T) {
	err := DeleteUser("ice", "123")
	if err != nil{
		t.Errorf("Errot")
	}
}

func testReGetUser(t *testing.T)  {
	pwd, err := GetUserCredential("ice")
	if err != nil{
		t.Errorf("Error of RegetUser %v", err)
	}
	if pwd != ""{
		t.Errorf("Deleting user test failed")
	}
}

func TestVideoWorkFlow(t *testing.T) {
	clearTables()
	t.Run("PrepareUser", testAddUser)
	t.Run("AddVideo", testAddVideoInfo)
	t.Run("GetVideo", testGetVideoInfo)
	t.Run("DelVideo", testDeleteVideoInfo)
	t.Run("RegetVideo", testRegetVideoInfo)
}

func testAddVideoInfo(t *testing.T) {
	vi, err := AddNewVideo(1, "my-video")
	if err != nil {
		t.Errorf("Error of AddVideoInfo: %v", err)
	}
	tempvid = vi.Id
}

func testGetVideoInfo(t *testing.T) {
	_, err := GetVideoInfo(tempvid)
	if err != nil {
		t.Errorf("Error of GetVideoInfo: %v", err)
	}
}

func testDeleteVideoInfo(t *testing.T) {
	err := DeleteVideoInfo(tempvid)
	if err != nil {
		t.Errorf("Error of DeleteVideoInfo: %v", err)
	}
}

func testRegetVideoInfo(t *testing.T) {
	vi, err := GetVideoInfo(tempvid)
	if err != nil || vi != nil{
		t.Errorf("Error of RegetVideoInfo: %v", err)
	}
}

func TestComments(t *testing.T)  {
	clearTables()
	t.Run("AddUser", testAddUser)
	t.Run("AddComments", testAddComments)
	t.Run("ListComments", testListComments)
	clearTables()
}
//
//func getTestUserId() (int, error) {
//	stmtOut, _ := dbConn.Prepare("select id from users where login_name='ice'")
//	rows, _ := stmtOut.Query()
//	res := make([]int, 1)
//	for rows.Next() {
//		var id int
//		if err := rows.Scan(&id); err != nil{
//			return 0, err
//		}
//		res = append(res, id)
//	}
//	return res[1], nil
//}

func testAddComments(t *testing.T)  {
	vid := "123456"
	aid := 106
	content := "I like this video"
	err := AddNewComments(vid, aid, content)
	if err != nil{
		t.Errorf("Error of Addcomments %v", err)
	}
}

func testListComments(t *testing.T) {
	vid := "123456"
	limit := 10
	offset := 1
	res, err := ListComments(vid, limit, offset)
	if err != nil{
		t.Errorf("Error of ListComments: %v", err)
	}
	for i, ele := range res{
		fmt.Printf("comment: %d, %v, \n", i, ele)
	}
}