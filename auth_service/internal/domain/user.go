package domain

type User struct {
	ID             string
	Login          string
	PassHash       []byte
	Roles          []*Role
	ClearanceLevel int
}
