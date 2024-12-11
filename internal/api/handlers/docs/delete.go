package docs

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

func (s *DocsHandler) Delete(c *gin.Context) {
	token := c.Query("token")
	if token == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "token is required"})
		return
	}
	docIdStr := c.Param("id")
	docId, err := strconv.ParseInt(docIdStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid document ID"})
		return
	}

	ctx := c.Request.Context()
	err = s.docServ.Delete(ctx, token, docId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to delete document"})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"response": gin.H{
			strconv.FormatInt(docId, 10): true,
		},
	})
}
