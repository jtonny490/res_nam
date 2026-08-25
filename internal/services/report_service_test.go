package services

import (
	"errors"
	"res_nam/internal/models"
	"testing"
	"time"
)

type fakeReports struct {
	created   *models.Report
	comment   *models.Comment
	commentAt time.Time
	like      *models.Like
	deleted   bool
	status    string
}

func (f *fakeReports) List(string, string, int, int) ([]models.Report, int64, error) {
	return nil, 0, nil
}
func (f *fakeReports) Get(uint) (*models.Report, error) { return &models.Report{}, nil }
func (f *fakeReports) Create(x *models.Report) error    { f.created = x; return nil }
func (f *fakeReports) AddComment(x *models.Comment, at time.Time) error {
	f.comment = x
	f.commentAt = at
	return nil
}
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

func TestReportServiceCommentSetsAuthorityFlagAndActivityTime(t *testing.T) {
	f := &fakeReports{}
	now := time.Date(2026, 8, 25, 10, 0, 0, 0, time.UTC)
	s := ReportService{Reports: f, Now: func() time.Time { return now }}

	comment, err := s.Comment(2, 3, "authority", "  Testing the water  ")
	if err != nil {
		t.Fatal(err)
	}
	if !comment.IsAuthorityComment || comment.Body != "Testing the water" {
		t.Fatalf("unexpected comment: %#v", comment)
	}
	if !f.commentAt.Equal(now) {
		t.Fatalf("activity time = %v, want %v", f.commentAt, now)
	}
}

func TestReportServiceCommentDoesNotFlagOtherRoles(t *testing.T) {
	for _, role := range []string{"public", "admin", ""} {
		t.Run(role, func(t *testing.T) {
			comment, err := (ReportService{Reports: &fakeReports{}}).Comment(2, 3, role, "Observation")
			if err != nil {
				t.Fatal(err)
			}
			if comment.IsAuthorityComment {
				t.Fatalf("role %q was marked as an authority comment", role)
			}
		})
	}
}

func TestReportServiceResolve(t *testing.T) {
	f := &fakeReports{}
	now := time.Date(2026, 8, 25, 10, 0, 0, 0, time.UTC)
	s := ReportService{Reports: f, Now: func() time.Time { return now }}

	if err := s.Resolve(9); err != nil {
		t.Fatal(err)
	}
	if f.status != "resolved" {
		t.Fatalf("status = %q, want resolved", f.status)
	}
}

func TestReportServiceResolveRejectsInvalidID(t *testing.T) {
	if err := (ReportService{Reports: &fakeReports{}}).Resolve(0); err == nil {
		t.Fatal("expected invalid report id error")
	}
}
