package handlers

import (
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"net/http"
	"res_nam/internal/models"
	"res_nam/internal/services"
	"strconv"
	"time"
)

type ReportHandler struct{ DB *gorm.DB }

func (h ReportHandler) List(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	if page < 1 {
		page = 1
	}
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "20"))
	if limit < 1 || limit > 100 {
		limit = 20
	}
	q := h.DB.Model(&models.Report{})
	if v := c.Query("category"); v != "" {
		q = q.Where("category = ?", v)
	}
	if v := c.Query("status"); v != "" {
		q = q.Where("status = ?", v)
	}
	var n int64
	q.Count(&n)
	var a []models.Report
	if e := q.Order("created_at desc").Offset((page - 1) * limit).Limit(limit).Find(&a).Error; e != nil {
		errJSON(c, 500, e)
		return
	}
	c.JSON(200, gin.H{"reports": a, "total": n, "page": page, "limit": limit})
}
func (h ReportHandler) Get(c *gin.Context) {
	var x models.Report
	if e := h.DB.Preload("User").Preload("Comments.User").Preload("Likes").First(&x, c.Param("id")).Error; e != nil {
		errJSON(c, 404, e)
		return
	}
	c.JSON(200, x)
}
func uid(c *gin.Context) uint { v, _ := c.Get("user_id"); return v.(uint) }
func (h ReportHandler) Create(c *gin.Context) {
	var x struct {
		Title, Description, PhotoURL, Category string
		Severity                               int
		Latitude, Longitude                    float64
	}
	if c.ShouldBindJSON(&x) != nil || x.Title == "" || x.Category == "" {
		c.JSON(400, gin.H{"error": "title and category are required"})
		return
	}
	if x.Severity < 1 || x.Severity > 5 {
		c.JSON(400, gin.H{"error": "severity must be between 1 and 5"})
		return
	}
	if !services.InLakeVictoriaCoverage(x.Latitude, x.Longitude) {
		c.JSON(400, gin.H{"error": "location is outside the Lake Victoria pilot coverage"})
		return
	}
	now := time.Now()
	r := models.Report{UserID: uid(c), Title: x.Title, Description: x.Description, PhotoURL: x.PhotoURL, Category: x.Category, Severity: x.Severity, Latitude: x.Latitude, Longitude: x.Longitude, Status: "open", LastActivityAt: now}
	if e := h.DB.Create(&r).Error; e != nil {
		errJSON(c, 500, e)
		return
	}
	c.JSON(http.StatusCreated, r)
}
func (h ReportHandler) Comment(c *gin.Context) {
	var x struct{ Body string }
	if c.ShouldBindJSON(&x) != nil || x.Body == "" {
		c.JSON(400, gin.H{"error": "body is required"})
		return
	}
	id, _ := strconv.Atoi(c.Param("id"))
	co := models.Comment{ReportID: uint(id), UserID: uid(c), Body: x.Body}
	if e := h.DB.Create(&co).Error; e != nil {
		errJSON(c, 500, e)
		return
	}
	c.JSON(201, co)
}
func (h ReportHandler) Like(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	l := models.Like{ReportID: uint(id), UserID: uid(c)}
	if e := h.DB.Where("report_id = ? AND user_id = ?", l.ReportID, l.UserID).First(&l).Error; e == nil {
		h.DB.Delete(&l)
		c.JSON(200, gin.H{"liked": false})
		return
	}
	h.DB.Create(&l)
	c.JSON(200, gin.H{"liked": true})
}
func (h ReportHandler) Status(c *gin.Context) {
	var x struct{ Status string }
	if c.ShouldBindJSON(&x) != nil || x.Status == "" {
		c.JSON(400, gin.H{"error": "status is required"})
		return
	}
	if e := h.DB.Model(&models.Report{}).Where("id = ?", c.Param("id")).Updates(map[string]interface{}{"status": x.Status, "last_activity_at": time.Now()}).Error; e != nil {
		errJSON(c, 500, e)
		return
	}
	c.JSON(200, gin.H{"status": x.Status})
}
