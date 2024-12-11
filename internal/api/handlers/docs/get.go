package docs

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

func (s *DocsHandler) Get(c *gin.Context) {
	token := c.Query("token")
	docIdStr := c.Param("id")
	docId, err := strconv.ParseInt(docIdStr, 10, 64)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": gin.H{
				"code": 400,
				"text": "Invalid document ID",
			},
		})
		return
	}
	ctx := c.Request.Context()
	data, doc, err := s.docServ.Get(ctx, token, docId)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error": gin.H{
				"code": 404,
				"text": err.Error(),
			},
		})
		return
	}
	if doc.IsFile {
		c.Header("Content-Type", doc.Mime)
		c.File(fmt.Sprintf("./uploads/%s", doc.Name))
		return
	}
	if doc.Mime == "json" {
		c.JSON(http.StatusOK, gin.H{
			"data": data,
		})
		return
	}
	// Для других mime типов возвращаем пустой объект в data
	c.JSON(http.StatusOK, gin.H{
		"data": gin.H{},
	})
	return
}
