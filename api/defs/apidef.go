package defs

type UserCredential struct {
	Username string `json:"user_name"`
	Pwd string `json:"pwd"`
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

type Comment struct {
	Id string
	VideoId string
	Author string
	Content string
}