package auth

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

func (s *AuthHandler) Register(c *gin.Context) {
	var model registerRequest
	if err := c.BindJSON(&model); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": gin.H{
				"code": http.StatusBadRequest,
				"text": "Invalid request body",
			},
		})
		return
	}

	if len(model.Login) < 8 {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": gin.H{
				"code": http.StatusBadRequest,
				"text": "Login must be at least 8 characters long",
			},
		})
		return
	}

	ctx := c.Request.Context()
	if err := s.authServ.Register(ctx, model.Token, model.Login, model.Pswd); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": gin.H{
				"code": http.StatusInternalServerError,
				"text": "Internal server error",
			},
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"response": gin.H{
			"login": model.Login,
		},
	})
}

type registerRequest struct {
	Token string `json:"token"`
	Login string `json:"login"`
	Pswd  string `json:"pswd"`
}
