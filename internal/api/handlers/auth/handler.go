package auth

import "doc-server/internal/services"

type AuthHandler struct {
	authServ services.AuthService
}

func NewSongHandler(authServ services.AuthService) *AuthHandler {
	return &AuthHandler{authServ: authServ}
}
