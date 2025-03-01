package sisconf

type LoginData struct {
	Email    string
	Password string
}

type LoginResponse struct {
	AuthenticationToken string
	AccessToken         string
}
