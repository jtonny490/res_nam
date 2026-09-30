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
	r.Use(middleware.CORS(c.CORSOrigin))
	r.Static("/uploads", "uploads")
	r.Static("/static", "static")
	r.NoRoute(func(ctx *gin.Context) {
		ctx.File("static/index.html")
	})
	r.GET("/health", func(x *gin.Context) { x.JSON(200, gin.H{"status": "ok"}) })
	r.POST("/api/auth/register", auth.Register)
	r.POST("/api/auth/login", auth.Login)

	reportRepo := repositories.ReportRepository{DB: db}
	rr := handlers.ReportHandler{S: services.ReportService{Reports: reportRepo}}
	ah := handlers.AuthorityHandler{S: services.AuthorityService{Requests: repositories.AuthorityRepository{DB: db}}}
	kp := handlers.KijaniHandler{DB: db, Provider: services.MockKijaniProvider{}}
	adm := handlers.AdminHandler{S: services.AdminService{
		UserStore:   repositories.UserRepository{DB: db},
		ReportStore: reportRepo,
	}}

	r.GET("/api/reports", rr.List)
	r.GET("/api/reports/:id", rr.Get)
	r.GET("/api/reports/:id/comments", rr.Comments)

	p := r.Group("/api", middleware.Auth(c.JWTSecret))
	p.POST("/reports", rr.Create)
	p.PATCH("/reports/:id", rr.Update)
	p.DELETE("/reports/:id", rr.Delete)
	p.POST("/reports/:id/comments", rr.Comment)
	p.POST("/reports/:id/like", rr.Like)
	p.DELETE("/reports/:id/like", rr.Unlike)
	p.POST("/reports/:id/likes", rr.Like)
	p.DELETE("/reports/:id/likes", rr.Unlike)
	p.PATCH("/reports/:id/status", middleware.RequireRole("authority", "admin"), rr.Status)
	p.PATCH("/reports/:id/resolve", middleware.RequireRole("admin"), rr.Resolve)
	p.GET("/reports/:id/assessment", kp.Assess)
	p.POST("/authority-requests", ah.Apply)
	p.GET("/authority-requests", middleware.RequireRole("admin"), ah.Pending)
	p.PATCH("/authority-requests/:id", middleware.RequireRole("admin"), ah.Review)
	p.GET("/map/reports", middleware.RequireRole("authority", "admin"), rr.Pins)

	admin := p.Group("/admin", middleware.RequireRole("admin"))
	admin.GET("/users", adm.Users)
	admin.PATCH("/users/:id", adm.UpdateUser)
	admin.GET("/analytics", adm.Analytics)
	return r
}
