package db

import (
	modelsServ "doc-server/internal/models"
	"doc-server/internal/storage/db/pg/models"
)

func ConvertUserRepoToUserServ(user models.User) modelsServ.User {
	return modelsServ.User{
		Id:        user.Id,
		Login:     user.Login,
		Password:  user.Password,
		CreatedAt: user.CreatedAt,
	}
}

func ConverterDocumentRepoToDocumentServ(document models.Document) modelsServ.Document {
	return modelsServ.Document{
		Id:        document.Id,
		Name:      document.Name,
		Mime:      document.Mime,
		IsPublic:  document.IsPublic,
		IsFile:    document.IsFile,
		UserId:    document.UserId,
		CreatedAt: document.CreatedAt,
		Grant:     ConvertGrantsRepoToGrantsServ(document.Grant),
	}
}

func ConvertGrantsRepoToGrantsServ(grants []models.Grant) []modelsServ.Grant {
	result := make([]modelsServ.Grant, 0, len(grants))
	for i := 0; i < len(grants); i++ {
		result = append(result, ConvertGrantRepoToGrantServ(grants[i]))
	}
	return result
}

func ConvertGrantRepoToGrantServ(grant models.Grant) modelsServ.Grant {
	return modelsServ.Grant{
		Id:         grant.Id,
		Login:      grant.Login,
		DocumentId: grant.DocumentId,
		Permission: grant.Permission,
		CreatedAt:  grant.CreatedAt,
	}
}
