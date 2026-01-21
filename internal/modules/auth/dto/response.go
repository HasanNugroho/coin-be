package dto

type TokenPair struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}

type LoginResponse struct {
	User      interface{} `json:"user"`
	TokenPair TokenPair   `json:"token_pair"`
}

type RegisterResponse struct {
	User interface{} `json:"user"`
}

type RefreshTokenResponse struct {
	AccessToken string `json:"access_token"`
}
