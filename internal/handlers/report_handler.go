package handlers

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"res_nam/internal/services"
	"strconv"
	"strings"
)

type ReportHandler struct{ S services.ReportService }

func uid(c *gin.Context) uint    { v, _ := c.Get("user_id"); id, _ := v.(uint); return id }
func role(c *gin.Context) string { v, _ := c.Get("role"); role, _ := v.(string); return role }

func atoi(s string) int     { n, _ := strconv.Atoi(s); return n }
func atof(s string) float64 { n, _ := strconv.ParseFloat(s, 64); return n }

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
	if strings.HasPrefix(c.ContentType(), "multipart/form-data") {
		if err := c.Request.ParseMultipartForm(10 << 20); err != nil {
			c.JSON(400, gin.H{"error": "invalid request"})
			return
		}
		x.Title = c.PostForm("title")
		x.Description = c.PostForm("description")
		x.Category = c.PostForm("category")
		x.Severity = atoi(c.PostForm("severity"))
		x.Latitude = atof(c.PostForm("latitude"))
		x.Longitude = atof(c.PostForm("longitude"))
		if f, err := c.FormFile("photo"); err == nil && f != nil {
			url, err := saveUpload(c, f)
			if err != nil {
				c.JSON(400, gin.H{"error": err.Error()})
				return
			}
			x.PhotoURL = url
		}
	} else if c.ShouldBindJSON(&x) != nil {
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
func (h ReportHandler) Update(c *gin.Context) {
	var in services.UpdateReportInput
	if c.ShouldBindJSON(&in) != nil {
		c.JSON(400, gin.H{"error": "invalid request"})
		return
	}
	id, _ := strconv.ParseUint(c.Param("id"), 10, 64)
	r, e := h.S.Update(uint(id), uid(c), role(c), in)
	if e != nil {
		domainErr(c, e)
		return
	}
	c.JSON(200, r)
}
func (h ReportHandler) Delete(c *gin.Context) {
	id, _ := strconv.ParseUint(c.Param("id"), 10, 64)
	if e := h.S.Delete(uint(id), uid(c), role(c)); e != nil {
		domainErr(c, e)
		return
	}
	c.Status(http.StatusNoContent)
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
func (h ReportHandler) Comments(c *gin.Context) {
	id, _ := strconv.ParseUint(c.Param("id"), 10, 64)
	x, e := h.S.Comments(uint(id))
	if e != nil {
		errJSON(c, 400, e)
		return
	}
	c.JSON(200, gin.H{"comments": x})
}
func (h ReportHandler) Like(c *gin.Context) {
	id, _ := strconv.ParseUint(c.Param("id"), 10, 64)
	if e := h.S.Like(uint(id), uid(c)); e != nil {
		errJSON(c, 400, e)
		return
	}
	c.JSON(200, gin.H{"liked": true})
}
func (h ReportHandler) Unlike(c *gin.Context) {
	id, _ := strconv.ParseUint(c.Param("id"), 10, 64)
	if e := h.S.Unlike(uint(id), uid(c)); e != nil {
		errJSON(c, 400, e)
		return
	}
	c.JSON(200, gin.H{"liked": false})
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

// Pins serves the lightweight map overlay data (authority/admin only).
func (h ReportHandler) Pins(c *gin.Context) {
	pins, e := h.S.Pins()
	if e != nil {
		errJSON(c, 500, e)
		return
	}
	c.JSON(200, gin.H{"pins": pins})
}
