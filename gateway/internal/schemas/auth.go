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
