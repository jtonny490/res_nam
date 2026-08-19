package handlers

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"res_nam/internal/services"
)

type AuthHandler struct{ S *services.AuthService }

func (h AuthHandler) Register(c *gin.Context) {
	var x struct{ Name, Email, Password string }
	if c.BindJSON(&x) != nil {
		errJSON(c, 400, gin.Error{Err: gin.ErrorTypeBind})
		return
	}
	u, t, e := h.S.Register(x.Name, x.Email, x.Password)
	if e != nil {
		errJSON(c, 400, e)
		return
	}
	c.JSON(http.StatusCreated, gin.H{"user": gin.H{"id": u.ID, "name": u.Name, "email": u.Email, "role": u.Role}, "token": t})
}
func (h AuthHandler) Login(c *gin.Context) {
	var x struct{ Email, Password string }
	if c.BindJSON(&x) != nil {
		c.JSON(400, gin.H{"error": "invalid request"})
		return
	}
	u, t, e := h.S.Login(x.Email, x.Password)
	if e != nil {
		errJSON(c, 401, e)
		return
	}
	c.JSON(200, gin.H{"user": gin.H{"id": u.ID, "name": u.Name, "email": u.Email, "role": u.Role}, "token": t})
}
