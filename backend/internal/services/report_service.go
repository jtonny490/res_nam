package services

import (
	"errors"
	"res_nam/internal/models"
	"strings"
	"time"
)

var ErrForbidden = errors.New("forbidden")

var validCategories = map[string]bool{
	"water": true, "air": true, "waste": true, "deforestation": true, "other": true,
}

// Categories returns the accepted report categories in display order.
func Categories() []string {
	return []string{"water", "air", "waste", "deforestation", "other"}
}

func ValidCategory(c string) bool { return validCategories[c] }

type ReportStore interface {
	List(string, string, int, int) ([]models.Report, int64, error)
	Get(uint) (*models.Report, error)
	Create(*models.Report) error
	Update(*models.Report) error
	Delete(uint) error
	AddComment(*models.Comment, time.Time) error
	ListComments(uint) ([]models.Comment, error)
	FindLike(uint, uint) (*models.Like, error)
	CreateLike(*models.Like) error
	DeleteLike(*models.Like) error
	UpdateStatus(uint, string, time.Time) error
	Pins() ([]models.ReportPin, error)
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

// UpdateReportInput uses pointers so callers can patch individual fields.
type UpdateReportInput struct {
	Title       *string
	Description *string
	PhotoURL    *string
	Category    *string
	Severity    *int
	Latitude    *float64
	Longitude   *float64
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
	if !ValidCategory(in.Category) {
		return nil, errors.New("category must be one of: " + strings.Join(Categories(), ", "))
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

// Update patches a report. Only the owner or an admin may edit it.
func (s ReportService) Update(id, uid uint, role string, in UpdateReportInput) (*models.Report, error) {
	if id == 0 || uid == 0 {
		return nil, errors.New("invalid report or user")
	}
	r, err := s.Reports.Get(id)
	if err != nil {
		return nil, err
	}
	if r.UserID != uid && role != "admin" {
		return nil, ErrForbidden
	}
	if in.Title != nil {
		t := strings.TrimSpace(*in.Title)
		if t == "" {
			return nil, errors.New("title cannot be empty")
		}
		r.Title = t
	}
	if in.Description != nil {
		r.Description = *in.Description
	}
	if in.PhotoURL != nil {
		r.PhotoURL = *in.PhotoURL
	}
	if in.Category != nil {
		c := strings.TrimSpace(*in.Category)
		if !ValidCategory(c) {
			return nil, errors.New("category must be one of: " + strings.Join(Categories(), ", "))
		}
		r.Category = c
	}
	if in.Severity != nil {
		if *in.Severity < 1 || *in.Severity > 5 {
			return nil, errors.New("severity must be between 1 and 5")
		}
		r.Severity = *in.Severity
	}
	if in.Latitude != nil {
		r.Latitude = *in.Latitude
	}
	if in.Longitude != nil {
		r.Longitude = *in.Longitude
	}
	if !InLakeVictoriaCoverage(r.Latitude, r.Longitude) {
		return nil, errors.New("location is outside the Lake Victoria pilot coverage")
	}
	return r, s.Reports.Update(r)
}

// Delete removes a report. Only the owner or an admin may delete it.
func (s ReportService) Delete(id, uid uint, role string) error {
	if id == 0 || uid == 0 {
		return errors.New("invalid report or user")
	}
	r, err := s.Reports.Get(id)
	if err != nil {
		return err
	}
	if r.UserID != uid && role != "admin" {
		return ErrForbidden
	}
	return s.Reports.Delete(id)
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

// Comments returns the comment thread for a report, oldest first.
func (s ReportService) Comments(rid uint) ([]models.Comment, error) {
	if rid == 0 {
		return nil, errors.New("invalid report id")
	}
	return s.Reports.ListComments(rid)
}

// Like is idempotent: liking an already-liked report succeeds without duplicating.
func (s ReportService) Like(rid, uid uint) error {
	if rid == 0 || uid == 0 {
		return errors.New("invalid report or user")
	}
	if _, err := s.Reports.FindLike(rid, uid); err == nil {
		return nil
	}
	return s.Reports.CreateLike(&models.Like{ReportID: rid, UserID: uid})
}

// Unlike is idempotent: unliking a report that is not liked succeeds.
func (s ReportService) Unlike(rid, uid uint) error {
	if rid == 0 || uid == 0 {
		return errors.New("invalid report or user")
	}
	x, err := s.Reports.FindLike(rid, uid)
	if err != nil {
		return nil
	}
	return s.Reports.DeleteLike(x)
}

func (s ReportService) UpdateStatus(id uint, status string) error {
	switch status {
	case "open", "investigating", "stale", "resolved":
	default:
		return errors.New("invalid report status")
	}
	return s.Reports.UpdateStatus(id, status, s.now())
}

func (s ReportService) Resolve(id uint) error {
	if id == 0 {
		return errors.New("invalid report id")
	}
	return s.Reports.UpdateStatus(id, "resolved", s.now())
}

// Pins returns the minimal report data needed to draw map markers.
func (s ReportService) Pins() ([]models.ReportPin, error) {
	return s.Reports.Pins()
}
