package handlers

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"res_nam/internal/services"
	"strconv"
)

type ReportHandler struct{ S services.ReportService }

func uid(c *gin.Context) uint    { v, _ := c.Get("user_id"); id, _ := v.(uint); return id }
func role(c *gin.Context) string { v, _ := c.Get("role"); role, _ := v.(string); return role }
func (h ReportHandler) List(c *gin.Context) {
	p, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	if p < 1 {
		p = 1
	}
	l, _ := strconv.Atoi(c.DefaultQuery("limit", "20"))
	if l < 1 || l > 100 {
		l = 20
	}
	a, n, e := h.S.List(c.Query("category"), c.Query("status"), p, l)
	if e != nil {
		errJSON(c, 500, e)
		return
	}
	c.JSON(200, gin.H{"reports": a, "total": n, "page": p, "limit": l})
}
func (h ReportHandler) Get(c *gin.Context) {
	id, _ := strconv.ParseUint(c.Param("id"), 10, 64)
	x, e := h.S.Get(uint(id))
	if e != nil {
		errJSON(c, 404, e)
		return
	}
	c.JSON(200, x)
}
func (h ReportHandler) Create(c *gin.Context) {
	var x services.CreateReportInput
	if c.ShouldBindJSON(&x) != nil {
		c.JSON(400, gin.H{"error": "invalid request"})
		return
	}
	r, e := h.S.Create(uid(c), x)
	if e != nil {
		errJSON(c, 400, e)
		return
	}
	c.JSON(http.StatusCreated, r)
}
func (h ReportHandler) Comment(c *gin.Context) {
	var x struct{ Body string }
	if c.ShouldBindJSON(&x) != nil {
		c.JSON(400, gin.H{"error": "invalid request"})
		return
	}
	id, _ := strconv.ParseUint(c.Param("id"), 10, 64)
	co, e := h.S.Comment(uint(id), uid(c), role(c), x.Body)
	if e != nil {
		errJSON(c, 400, e)
		return
	}
	c.JSON(http.StatusCreated, co)
}
func (h ReportHandler) Like(c *gin.Context) {
	id, _ := strconv.ParseUint(c.Param("id"), 10, 64)
	liked, e := h.S.ToggleLike(uint(id), uid(c))
	if e != nil {
		errJSON(c, 500, e)
		return
	}
	c.JSON(200, gin.H{"liked": liked})
}
func (h ReportHandler) Status(c *gin.Context) {
	var x struct{ Status string }
	if c.ShouldBindJSON(&x) != nil {
		c.JSON(400, gin.H{"error": "invalid request"})
		return
	}
	id, _ := strconv.ParseUint(c.Param("id"), 10, 64)
	if e := h.S.UpdateStatus(uint(id), x.Status); e != nil {
		errJSON(c, 400, e)
		return
	}
	c.JSON(200, gin.H{"status": x.Status})
}

func (h ReportHandler) Resolve(c *gin.Context) {
	id, _ := strconv.ParseUint(c.Param("id"), 10, 64)
	if e := h.S.Resolve(uint(id)); e != nil {
		errJSON(c, 400, e)
		return
	}
	c.JSON(200, gin.H{"status": "resolved"})
}
