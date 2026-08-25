package services

import (
	"errors"
	"res_nam/internal/models"
	"strings"
	"time"
)

type ReportStore interface {
	List(string, string, int, int) ([]models.Report, int64, error)
	Get(uint) (*models.Report, error)
	Create(*models.Report) error
	AddComment(*models.Comment, time.Time) error
	FindLike(uint, uint) (*models.Like, error)
	CreateLike(*models.Like) error
	DeleteLike(*models.Like) error
	UpdateStatus(uint, string, time.Time) error
}
type ReportService struct {
	Reports ReportStore
	Now     func() time.Time
}
type CreateReportInput struct {
	Title, Description, PhotoURL, Category string
	Severity                               int
	Latitude, Longitude                    float64
}

func (s ReportService) now() time.Time {
	if s.Now != nil {
		return s.Now()
	}
	return time.Now()
}
func (s ReportService) List(c, st string, p, l int) ([]models.Report, int64, error) {
	return s.Reports.List(c, st, p, l)
}
func (s ReportService) Get(id uint) (*models.Report, error) {
	if id == 0 {
		return nil, errors.New("invalid report id")
	}
	return s.Reports.Get(id)
}
func (s ReportService) Create(uid uint, in CreateReportInput) (*models.Report, error) {
	in.Title = strings.TrimSpace(in.Title)
	in.Category = strings.TrimSpace(in.Category)
	if uid == 0 {
		return nil, errors.New("invalid user")
	}
	if in.Title == "" || in.Category == "" {
		return nil, errors.New("title and category are required")
	}
	if in.Severity < 1 || in.Severity > 5 {
		return nil, errors.New("severity must be between 1 and 5")
	}
	if !InLakeVictoriaCoverage(in.Latitude, in.Longitude) {
		return nil, errors.New("location is outside the Lake Victoria pilot coverage")
	}
	r := &models.Report{UserID: uid, Title: in.Title, Description: in.Description, PhotoURL: in.PhotoURL, Category: in.Category, Severity: in.Severity, Latitude: in.Latitude, Longitude: in.Longitude, Status: "open", LastActivityAt: s.now()}
	return r, s.Reports.Create(r)
}
func (s ReportService) Comment(rid, uid uint, role, body string) (*models.Comment, error) {
	body = strings.TrimSpace(body)
	if rid == 0 || uid == 0 {
		return nil, errors.New("invalid report or user")
	}
	if body == "" {
		return nil, errors.New("body is required")
	}
	x := &models.Comment{
		ReportID:           rid,
		UserID:             uid,
		Body:               body,
		IsAuthorityComment: role == "authority",
	}
	return x, s.Reports.AddComment(x, s.now())
}
func (s ReportService) ToggleLike(rid, uid uint) (bool, error) {
	if rid == 0 || uid == 0 {
		return false, errors.New("invalid report or user")
	}
	x, err := s.Reports.FindLike(rid, uid)
	if err == nil {
		return false, s.Reports.DeleteLike(x)
	}
	x = &models.Like{ReportID: rid, UserID: uid}
	return true, s.Reports.CreateLike(x)
}
func (s ReportService) UpdateStatus(id uint, status string) error {
	switch status {
	case "open", "investigating", "stale", "resolved":
	default:
		return errors.New("invalid report status")
	}
	return s.Reports.UpdateStatus(id, status, s.now())
}
