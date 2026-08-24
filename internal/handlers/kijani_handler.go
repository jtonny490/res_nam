package handlers

import (
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"res_nam/internal/models"
	"res_nam/internal/services"
)

type KijaniHandler struct {
	DB       *gorm.DB
	Provider services.KijaniProvider
}

func (h KijaniHandler) Assess(c *gin.Context) {
	var r models.Report
	if e := h.DB.First(&r, c.Param("id")).Error; e != nil {
		errJSON(c, 404, e)
		return
	}
	a, e := h.Provider.Assess(c.Request.Context(), r)
	if e != nil {
		errJSON(c, 503, e)
		return
	}
	c.JSON(200, a)
}
