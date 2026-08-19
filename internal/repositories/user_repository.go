package repositories

import (
	"gorm.io/gorm"
	"res_nam/internal/models"
)

type UserRepository struct{ DB *gorm.DB }

func (r UserRepository) Create(u *models.User) error { return r.DB.Create(u).Error }
func (r UserRepository) FindByEmail(e string) (*models.User, error) {
	var u models.User
	err := r.DB.Where("email = ?", e).First(&u).Error
	return &u, err
}
