package docs

import (
	"doc-server/internal/config"
	"doc-server/internal/services"
)

type DocsHandler struct {
	configApp config.ConfigApp
	authServ  services.AuthService
	docServ   services.DocumentService
}

func NewSongHandler(configApp config.ConfigApp,
	authServ services.AuthService,
	docServ services.DocumentService) *DocsHandler {
	return &DocsHandler{
		configApp: configApp,
		authServ:  authServ,
		docServ:   docServ,
	}
}
