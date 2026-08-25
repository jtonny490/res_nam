package repositories

import (
	"gorm.io/gorm"
	"res_nam/internal/models"
	"time"
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
	e := r.DB.Preload("User").Preload("Comments.User").Preload("Likes").First(&x, id).Error
	return &x, e
}
func (r ReportRepository) AddComment(x *models.Comment, at time.Time) error {
	return r.DB.Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(x).Error; err != nil {
			return err
		}
		updates := map[string]interface{}{"last_activity_at": at}
		if x.IsAuthorityComment {
			updates["status"] = gorm.Expr("CASE WHEN status = ? THEN ? ELSE status END", "open", "investigating")
		}
		return tx.Model(&models.Report{}).Where("id = ?", x.ReportID).Updates(updates).Error
	})
}
func (r ReportRepository) FindLike(reportID, userID uint) (*models.Like, error) {
	var x models.Like
	err := r.DB.Where("report_id = ? AND user_id = ?", reportID, userID).First(&x).Error
	return &x, err
}
func (r ReportRepository) CreateLike(x *models.Like) error { return r.DB.Create(x).Error }
func (r ReportRepository) DeleteLike(x *models.Like) error { return r.DB.Delete(x).Error }
func (r ReportRepository) UpdateStatus(id uint, status string, at time.Time) error {
	return r.DB.Model(&models.Report{}).Where("id = ?", id).Updates(map[string]interface{}{"status": status, "last_activity_at": at}).Error
}
