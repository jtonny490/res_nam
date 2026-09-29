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
func (r UserRepository) List() ([]models.User, error) {
	var users []models.User
	err := r.DB.Order("created_at").Find(&users).Error
	return users, err
}
func (r UserRepository) UpdateRoleStatus(id uint, role, status string) error {
	return r.DB.Model(&models.User{}).Where("id = ?", id).Updates(map[string]interface{}{"role": role, "status": status}).Error
}
func (r UserRepository) CountActive() (int64, error) {
	var n int64
	err := r.DB.Model(&models.User{}).Where("status = ?", "active").Count(&n).Error
	return n, err
}
