package sisconf

type LoginData struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type LoginResponse struct {
	AuthenticationToken string
	AccessToken         string
}
