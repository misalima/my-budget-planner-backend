package dto

// RefreshTokenDTO representa o body necessário para o refresh token.
type RefreshTokenDTO struct {
	Token string `json:"token"`
}
