package app

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"log"
	"os"
	"res_nam/internal/config"
	"res_nam/internal/database"
	"res_nam/internal/handlers"
	"res_nam/internal/repositories"
	"res_nam/internal/router"
	"res_nam/internal/services"
	"time"
)

func New(c config.Config) (*gorm.DB, error) {
	dsn := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable", c.DBHost, c.DBPort, c.DBUser, c.DBPassword, c.DBName)
	var db *gorm.DB
	var err error
	for i := 0; i < 20; i++ {
		db, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
		if err == nil {
			break
		}
		time.Sleep(2 * time.Second)
	}
	if err != nil {
		return nil, err
	}
	if err = database.Run(db, c.MigrationsDir); err != nil {
		return nil, err
	}
	return db, nil
}
func BuildRouter(c config.Config, db *gorm.DB) *gin.Engine {
	_ = log.Print
	_ = os.MkdirAll("uploads", 0755)
	auth := handlers.AuthHandler{S: &services.AuthService{Users: repositories.UserRepository{DB: db}, Secret: c.JWTSecret}}
	return router.New(c, db, auth)
}
