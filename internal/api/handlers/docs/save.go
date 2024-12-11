package docs

import (
	"doc-server/internal/utils"
	"encoding/json"
	"errors"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
)

func (s *DocsHandler) Save(c *gin.Context) {
	if err := c.Request.ParseMultipartForm(10 << 20); err != nil { // Лимит на 10MB
		c.JSON(http.StatusBadRequest, gin.H{"error": "Failed to parse multipart form"})
		return
	}

	meta := c.Request.FormValue("meta")
	if meta == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Missing metadata"})
		return
	}

	model, err := parseMetadata(meta)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	file, handler, err := c.Request.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "File is required"})
		return
	}
	defer file.Close()

	model.Name = handler.Filename
	model.File = true

	ctx := c.Request.Context()
	claims, err := utils.ParseToken(model.Token, s.configApp.SecretKey())
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
		return
	}

	doc, err := converterMetaToDocumentServ(model, claims.Id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to convert metadata to document"})
		log.Printf("Error converting metadata: %v", err)
		return
	}

	data, err := s.docServ.Save(ctx, file, doc)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save document"})
		log.Printf("Error saving document: %v", err)
		return
	}

	response := save_response{
		Data: save_data{
			Json: data,
			File: doc.Name,
		},
	}
	c.JSON(http.StatusOK, response)
}

func parseMetadata(meta string) (metaData, error) {
	var model metaData
	if err := json.Unmarshal([]byte(meta), &model); err != nil {
		return metaData{}, errors.New("invalid metadata format")
	}
	if model.Token == "" {
		return metaData{}, errors.New("missing token in metadata")
	}
	return model, nil
}

type metaData struct {
	Name   string   `json:"name"`
	File   bool     `json:"file"`
	Public bool     `json:"public"`
	Token  string   `json:"token" validate:"required"`
	Mime   string   `json:"mime"`
	Grant  []string `json:"grant"`
}

type save_response struct {
	Data save_data `json:"data"`
}

type save_data struct {
	Json string `json:"json,omitempty"`
	File string `json:"file"`
}
