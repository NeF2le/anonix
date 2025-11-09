package domain

type User struct {
	ID       string
	Login    string
	PassHash []byte
	RoleId   int
}
