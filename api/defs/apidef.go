package defs

type UserCredential struct {
	Username string `json:"user_name"`
	Pwd string `json:"pwd"`
}

type SignedUp struct {
	Success bool `json:"success"`
	SessionId string `json:"session_id"`
}

type SignedIn struct {
	Success bool `json:"success"`
	SessionId string `json:"session_id"`
}


type User struct {
	Id string
	LoginName string
	Pwd string
}

// {
//   user_name: xxx,
//   pwd: xxx,
// }

type VideoInfo struct {
	Id string
	AuthorId int
	Name string
	DisplayCtime string
}

type VideosInfo struct {
	Videos []*VideoInfo `json:"videos"`
}

type NewVideo struct {
	AuthorId int `json:"author_id"`
	Name string `json:"name"`
}


type Comment struct {
	Id string
	VideoId string
	Author string
	Content string
}

type Comments struct {
	Comments []*Comment `json:"comments"`
}


type NewComment struct {
	AuthorId int `json:"author_id"`
	Content string `json:"content"`
}

type SimpleSession struct {
	Username string	//login name
	TTL int64
}