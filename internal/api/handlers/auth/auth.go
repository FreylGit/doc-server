package auth

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

type authRequest struct {
	Login string `json:"login"`
	Pswd  string `json:"pswd"`
}

type authResponse struct {
	Token string `json:"token"`
}

func (s *AuthHandler) Auth(c *gin.Context) {
	var model authRequest
	if err := c.BindJSON(&model); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": gin.H{
				"code": http.StatusBadRequest,
				"text": "Invalid request body",
			},
		})
		return
	}

	ctx := c.Request.Context()
	token, err := s.authServ.Auth(ctx, model.Login, model.Pswd)
	if err != nil {
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
			"token": token,
		},
	})
}
