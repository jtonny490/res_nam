package repositories

import (
	"gorm.io/gorm"
	"res_nam/internal/models"
)

type AuthorityRepository struct{ DB *gorm.DB }

func (r AuthorityRepository) Create(x *models.AuthorityRequest) error { return r.DB.Create(x).Error }
func (r AuthorityRepository) PendingByUser(id uint) (*models.AuthorityRequest, error) {
	var x models.AuthorityRequest
	e := r.DB.Where("user_id = ? AND status = ?", id, "pending").First(&x).Error
	return &x, e
}
func (r AuthorityRepository) ListPending() ([]models.AuthorityRequest, error) {
	var x []models.AuthorityRequest
	e := r.DB.Where("status = ?", "pending").Order("created_at").Find(&x).Error
	return x, e
}
func (r AuthorityRepository) Review(id, admin uint, status string) error {
	return r.DB.Transaction(func(tx *gorm.DB) error {
		var x models.AuthorityRequest
		if e := tx.First(&x, id).Error; e != nil {
			return e
		}
		if x.Status != "pending" {
			return gorm.ErrInvalidData
		}
		if e := tx.Model(&x).Updates(map[string]interface{}{"status": status, "reviewed_by": admin}).Error; e != nil {
			return e
		}
		if status == "approved" {
			return tx.Model(&models.User{}).Where("id = ?", x.UserID).Update("role", "authority").Error
		}
		return nil
	})
}
