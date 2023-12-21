package domain

type User struct {
	Id       int64
	Email    string
	Password string
}

type UserProfile struct {
	Id       int64
	UserId   int64
	Nickname string
	Profile  string
	Birthday string
}
