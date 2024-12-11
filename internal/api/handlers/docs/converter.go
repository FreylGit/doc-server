package docs

import (
	"doc-server/internal/models"
	"strconv"
)

func converterMetaToDocumentServ(meta metaData, userId string) (models.Document, error) {
	id, err := strconv.ParseInt(userId, 10, 64)
	return models.Document{
		Name:     meta.Name,
		Mime:     meta.Mime,
		IsPublic: meta.Public,
		IsFile:   meta.File,
		UserId:   id,
		Grant:    converterGrantsReqToGrantsServ(meta.Grant),
	}, err
}

func converterGrantsReqToGrantsServ(grants []string) []models.Grant {
	model := make([]models.Grant, len(grants))
	for i := 0; i < len(grants); i++ {
		model[i] = models.Grant{
			Login: grants[i],
		}
	}
	return model
}

func extractLogins(grants []models.Grant) []string {
	var logins []string
	for _, grant := range grants {
		logins = append(logins, grant.Login)
	}
	return logins
}
