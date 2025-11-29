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
	UserId       string `json:"user_id"`
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}

type RefreshRespSchema struct {
	RefreshToken string `json:"refresh_token"`
	AccessToken  string `json:"access_token"`
}

type IsAdminRespSchema struct {
	Result bool `json:"result"`
}
