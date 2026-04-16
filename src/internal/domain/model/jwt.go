package model

type JwtRequest struct {
	Login    string
	Password string
}

type JwtResponse struct {
	AccessToken  string
	RefreshToken string
}

type RefreshJwtRequest struct {
	RefreshToken string
}
