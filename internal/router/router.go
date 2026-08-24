package router

import (
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"res_nam/internal/config"
	"res_nam/internal/handlers"
	"res_nam/internal/middleware"
	"res_nam/internal/repositories"
	"res_nam/internal/services"
)

func New(c config.Config, db *gorm.DB, auth handlers.AuthHandler) *gin.Engine {
	r := gin.Default()
	r.Use(middleware.CORS())
	r.Static("/uploads", "uploads")
	r.StaticFileFS("/", "frontend/index.html", gin.Dir("frontend", false))
	r.GET("/health", func(x *gin.Context) { x.JSON(200, gin.H{"status": "ok"}) })
	r.POST("/api/auth/register", auth.Register)
	r.POST("/api/auth/login", auth.Login)
	rr := handlers.ReportHandler{S: services.ReportService{Reports: repositories.ReportRepository{DB: db}}}
	ah := handlers.AuthorityHandler{S: services.AuthorityService{Requests: repositories.AuthorityRepository{DB: db}}}
	kp := handlers.KijaniHandler{DB: db, Provider: services.MockKijaniProvider{}}
	r.GET("/api/reports", rr.List)
	r.GET("/api/reports/:id", rr.Get)
	p := r.Group("/api", middleware.Auth(c.JWTSecret))
	p.POST("/reports", rr.Create)
	p.POST("/reports/:id/comments", rr.Comment)
	p.POST("/reports/:id/like", rr.Like)
	p.PATCH("/reports/:id/status", middleware.RequireRole("authority", "admin"), rr.Status)
	p.GET("/reports/:id/assessment", kp.Assess)
	p.POST("/authority-requests", ah.Apply)
	p.GET("/authority-requests", middleware.RequireRole("admin"), ah.Pending)
	p.PATCH("/authority-requests/:id", middleware.RequireRole("admin"), ah.Review)
	return r
}
