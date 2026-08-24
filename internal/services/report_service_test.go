package services

import (
	"errors"
	"res_nam/internal/models"
	"testing"
	"time"
)

type fakeReports struct {
	created *models.Report
	comment *models.Comment
	like    *models.Like
	deleted bool
	status  string
}

func (f *fakeReports) List(string, string, int, int) ([]models.Report, int64, error) {
	return nil, 0, nil
}
func (f *fakeReports) Get(uint) (*models.Report, error)   { return &models.Report{}, nil }
func (f *fakeReports) Create(x *models.Report) error      { f.created = x; return nil }
func (f *fakeReports) AddComment(x *models.Comment) error { f.comment = x; return nil }
func (f *fakeReports) FindLike(uint, uint) (*models.Like, error) {
	if f.like != nil {
		return f.like, nil
	}
	return nil, errors.New("not found")
}
func (f *fakeReports) CreateLike(x *models.Like) error                  { f.like = x; return nil }
func (f *fakeReports) DeleteLike(x *models.Like) error                  { f.deleted = true; return nil }
func (f *fakeReports) UpdateStatus(_ uint, s string, _ time.Time) error { f.status = s; return nil }
func TestReportServiceCreateValidatesAndSetsDefaults(t *testing.T) {
	f := &fakeReports{}
	now := time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC)
	s := ReportService{Reports: f, Now: func() time.Time { return now }}
	r, e := s.Create(7, CreateReportInput{Title: "  Spill ", Category: "water", Severity: 4, Latitude: -1, Longitude: 34})
	if e != nil || r.Title != "Spill" || r.Status != "open" || !r.LastActivityAt.Equal(now) {
		t.Fatalf("unexpected result: %#v %v", r, e)
	}
}
func TestReportServiceRejectsOutOfBounds(t *testing.T) {
	_, e := (ReportService{Reports: &fakeReports{}}).Create(1, CreateReportInput{Title: "x", Category: "water", Severity: 2, Latitude: 10, Longitude: 34})
	if e == nil {
		t.Fatal("expected coverage error")
	}
}
func TestReportServiceLikeToggles(t *testing.T) {
	f := &fakeReports{}
	s := ReportService{Reports: f}
	liked, e := s.ToggleLike(2, 3)
	if e != nil || !liked {
		t.Fatal("expected like")
	}
	liked, e = s.ToggleLike(2, 3)
	if e != nil || liked || !f.deleted {
		t.Fatal("expected unlike")
	}
}
