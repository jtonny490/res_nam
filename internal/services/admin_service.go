package services

import (
	"errors"
	"res_nam/internal/models"
)

var validRoles = map[string]bool{"public": true, "authority": true, "admin": true}
var validUserStatuses = map[string]bool{"active": true, "pending": true, "banned": true}

func ValidRole(r string) bool       { return validRoles[r] }
func ValidUserStatus(s string) bool { return validUserStatuses[s] }
func Roles() []string               { return []string{"public", "authority", "admin"} }
func UserStatuses() []string        { return []string{"active", "pending", "banned"} }

type AdminUserStore interface {
	List() ([]models.User, error)
	UpdateRoleStatus(uint, string, string) error
	CountActive() (int64, error)
}

type AdminReportStore interface {
	CountAll() (int64, error)
	CountByCategory() ([]models.CategoryCount, error)
	CountByStatus() ([]models.StatusCount, error)
}

type AdminService struct {
	UserStore   AdminUserStore
	ReportStore AdminReportStore
}

// Analytics is the aggregate view served to the admin dashboard.
type Analytics struct {
	TotalReports int64                  `json:"total_reports"`
	ActiveUsers  int64                  `json:"active_users"`
	ByCategory   []models.CategoryCount `json:"by_category"`
	ByStatus     []models.StatusCount   `json:"by_status"`
}

func (s AdminService) ListUsers() ([]models.User, error) {
	return s.UserStore.List()
}

func (s AdminService) UpdateUser(id uint, role, status string) error {
	if id == 0 {
		return errors.New("invalid user id")
	}
	if !ValidRole(role) {
		return errors.New("role must be one of: public, authority, admin")
	}
	if !ValidUserStatus(status) {
		return errors.New("status must be one of: active, pending, banned")
	}
	return s.UserStore.UpdateRoleStatus(id, role, status)
}

func (s AdminService) Analytics() (Analytics, error) {
	var a Analytics
	total, err := s.ReportStore.CountAll()
	if err != nil {
		return a, err
	}
	users, err := s.UserStore.CountActive()
	if err != nil {
		return a, err
	}
	byCategory, err := s.ReportStore.CountByCategory()
	if err != nil {
		return a, err
	}
	byStatus, err := s.ReportStore.CountByStatus()
	if err != nil {
		return a, err
	}
	a.TotalReports = total
	a.ActiveUsers = users
	a.ByCategory = byCategory
	a.ByStatus = byStatus
	return a, nil
}
