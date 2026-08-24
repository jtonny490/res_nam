package handlers

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"res_nam/internal/services"
	"strconv"
)

type AuthorityHandler struct{ S services.AuthorityService }

func (h AuthorityHandler) Apply(c *gin.Context) {
	var x struct{ OrganizationName, Justification string }
	if c.ShouldBindJSON(&x) != nil {
		c.JSON(400, gin.H{"error": "invalid request"})
		return
	}
	r, e := h.S.Apply(uid(c), x.OrganizationName, x.Justification)
	if e != nil {
		errJSON(c, 400, e)
		return
	}
	c.JSON(http.StatusCreated, r)
}
func (h AuthorityHandler) Pending(c *gin.Context) {
	x, e := h.S.Pending()
	if e != nil {
		errJSON(c, 500, e)
		return
	}
	c.JSON(200, gin.H{"requests": x})
}
func (h AuthorityHandler) Review(c *gin.Context) {
	var x struct{ Status string }
	if c.ShouldBindJSON(&x) != nil {
		c.JSON(400, gin.H{"error": "invalid request"})
		return
	}
	id, _ := strconv.ParseUint(c.Param("id"), 10, 64)
	if e := h.S.Review(uint(id), uid(c), x.Status); e != nil {
		errJSON(c, 400, e)
		return
	}
	c.JSON(200, gin.H{"status": x.Status})
}
