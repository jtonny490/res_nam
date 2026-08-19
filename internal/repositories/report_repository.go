package repositories

import (
	"gorm.io/gorm"
	"res_nam/internal/models"
)

type ReportRepository struct{ DB *gorm.DB }

func (r ReportRepository) List(cat, status string, page, limit int) ([]models.Report, int64, error) {
	var a []models.Report
	var n int64
	q := r.DB.Model(&models.Report{})
	if cat != "" {
		q = q.Where("category = ?", cat)
	}
	if status != "" {
		q = q.Where("status = ?", status)
	}
	q.Count(&n)
	err := q.Order("created_at desc").Offset((page - 1) * limit).Limit(limit).Find(&a).Error
	return a, n, err
}
func (r ReportRepository) Create(x *models.Report) error { return r.DB.Create(x).Error }
func (r ReportRepository) Get(id uint) (*models.Report, error) {
	var x models.Report
	e := r.DB.First(&x, id).Error
	return &x, e
}
