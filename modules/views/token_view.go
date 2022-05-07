package views

type BearerTokenRead struct {
	AccessToken string `json:"access_token"`
	TokenType   string `json:"token_type"`
}

func NewBearerTokenRead(token string) BearerTokenRead {
	return BearerTokenRead{
		AccessToken: token,
		TokenType:   "Bearer",
	}
}
