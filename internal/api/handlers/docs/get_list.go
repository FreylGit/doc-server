package docs

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
	"time"
)

func (s *DocsHandler) GetList(c *gin.Context) {
	token := c.Query("token")
	if token == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "token is required"})
		return
	}

	limitStr := c.Query("limit")
	limit, err := strconv.ParseInt(limitStr, 10, 64)
	if err != nil || limit <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid limit"})
		return
	}

	ctx := c.Request.Context()
	filter := initFilter(c)

	list, err := s.docServ.GetList(ctx, token, filter, limit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to retrieve document list"})
		return
	}

	var transformedList []documentResponse
	for _, doc := range list {
		transformedDoc := documentResponse{
			Id:        doc.Id,
			Name:      doc.Name,
			Mime:      doc.Mime,
			IsPublic:  doc.IsPublic,
			IsFile:    doc.IsFile,
			UserId:    doc.UserId,
			CreatedAt: doc.CreatedAt,
			Grant:     extractLogins(doc.Grant),
		}
		transformedList = append(transformedList, transformedDoc)
	}

	c.JSON(http.StatusOK, gin.H{"data": transformedList})
}

type documentResponse struct {
	Id        int64     `json:"id"`
	Name      string    `json:"name"`
	Mime      string    `json:"mime"`
	IsPublic  bool      `json:"is_public"`
	IsFile    bool      `json:"is_file"`
	UserId    int64     `json:"user_id"`
	CreatedAt time.Time `json:"created_at"`
	Grant     []string  `json:"grant"`
}

func initFilter(c *gin.Context) map[string]interface{} {
	filters := map[string]interface{}{}
	if login := c.Query("login"); login != "" {
		filters["login"] = login
	}
	if mime := c.Query("mime"); mime != "" {
		filters["mime"] = mime
	}
	if isPublic, err := strconv.ParseBool(c.Query("public")); err == nil {
		filters["is_public"] = isPublic
	}
	if isFile, err := strconv.ParseBool(c.Query("file")); err == nil {
		filters["is_file"] = isFile
	}

	return filters
}
