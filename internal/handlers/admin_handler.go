package handlers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

	"res_nam/internal/services"
)

type AdminHandler struct{ S services.AdminService }

func (h AdminHandler) Users(c *gin.Context) {
	users, e := h.S.ListUsers()
	if e != nil {
		errJSON(c, 500, e)
		return
	}
	c.JSON(200, gin.H{"users": users})
}

func (h AdminHandler) UpdateUser(c *gin.Context) {
	var x struct{ Role, Status string }
	if c.ShouldBindJSON(&x) != nil {
		c.JSON(400, gin.H{"error": "invalid request"})
		return
	}
	id, _ := strconv.ParseUint(c.Param("id"), 10, 64)
	if e := h.S.UpdateUser(uint(id), x.Role, x.Status); e != nil {
		errJSON(c, 400, e)
		return
	}
	c.JSON(http.StatusOK, gin.H{"id": id, "role": x.Role, "status": x.Status})
}

func (h AdminHandler) Analytics(c *gin.Context) {
	a, e := h.S.Analytics()
	if e != nil {
		errJSON(c, 500, e)
		return
	}
	c.JSON(200, a)
}
