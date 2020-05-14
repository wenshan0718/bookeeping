package entity

//用户信息
type User struct {
	Id       int    `json:"id`       //用户id
	Name     string `json:"name"`    //用户名
	Password string `json:"password` //用户密码
}

func NewUser() *User {
	user := User{}
	return &user
}
