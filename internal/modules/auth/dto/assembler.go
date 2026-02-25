package dto

import "time"

func ToLoginResponse(accessToken string, accessExpiry time.Duration, refreshToken string, refreshExpiry time.Duration) *LoginResponse {
	return &LoginResponse{
		AccessToken:      accessToken,
		RefreshToken:     refreshToken,
		TokenType:        "Bearer",
		ExpiresIn:        int(accessExpiry.Seconds()),
		RefreshExpiresIn: int(refreshExpiry.Seconds()),
	}
}
