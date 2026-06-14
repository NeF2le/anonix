package schemas

type RegisterSchema struct {
	Login    string `json:"login"`
	Password string `json:"password"`
	RoleId   int    `json:"role_id"`
}

type LoginSchema struct {
	Login    string `json:"login"`
	Password string `json:"password"`
}

type RefreshSchema struct {
	RefreshToken string `json:"refresh_token"`
	AccessToken  string `json:"access_token"`
}

type IsAdminSchema struct {
	UserId string `json:"user_id"`
}

type RegisterRespSchema struct {
	UserId string `json:"user_id"`
}

type LoginRespSchema struct {
	UserId       string        `json:"user_id"`
	AccessToken  string        `json:"access_token"`
	RefreshToken string        `json:"refresh_token"`
	Roles        []*RoleSchema `json:"roles"`
}

type RefreshRespSchema struct {
	RefreshToken string `json:"refresh_token"`
	AccessToken  string `json:"access_token"`
}

type IsAdminRespSchema struct {
	Result bool `json:"result"`
}

type DeleteUserSchema struct {
	UserId string `json:"user_id"`
}

type DeleteUserRespSchema struct{}

type AssignRoleSchema struct {
	UserId string `json:"user_id"`
	RoleId int32  `json:"role_id"`
}

type AssignRoleRespSchema struct{}

type RemoveRoleSchema struct {
	UserId string `json:"user_id"`
	RoleId int32  `json:"role_id"`
}

type RemoveRoleRespSchema struct{}

type GetRolesListRespSchema struct {
	Roles []*RoleSchema `json:"roles"`
}

type RoleSchema struct {
	Id   int32  `json:"id"`
	Name string `json:"name"`
}

type UserSchema struct {
	Id             string        `json:"id"`
	Login          string        `json:"login"`
	Roles          []*RoleSchema `json:"roles"`
	ClearanceLevel int32         `json:"clearance_level"`
}

type UpdateClearanceSchema struct {
	UserId         string `json:"user_id"`
	ClearanceLevel int32  `json:"clearance_level"`
}

type UpdateClearanceRespSchema struct{}

type GetUsersRespSchema struct {
	Users []*UserSchema `json:"users"`
}

type GetUserRolesSchema struct {
	UserId string `json:"user_id"`
}

type GetUserRolesRespSchema struct {
	Roles []*RoleSchema `json:"roles"`
}
