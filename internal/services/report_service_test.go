package services

import (
	"errors"
	"res_nam/internal/models"
	"testing"
	"time"
)

type fakeReports struct {
	report    *models.Report
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
func (f *fakeReports) Get(id uint) (*models.Report, error) {
	if f.report != nil {
		return f.report, nil
	}
	return &models.Report{ID: id}, nil
}
func (f *fakeReports) Create(x *models.Report) error { f.created = x; return nil }
func (f *fakeReports) Update(x *models.Report) error { f.created = x; return nil }
func (f *fakeReports) Delete(uint) error             { f.deleted = true; return nil }
func (f *fakeReports) AddComment(x *models.Comment, at time.Time) error {
	f.comment = x
	f.commentAt = at
	return nil
}
func (f *fakeReports) ListComments(uint) ([]models.Comment, error) {
	return nil, nil
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
func (f *fakeReports) Pins() ([]models.ReportPin, error)                { return nil, nil }
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
func TestReportServiceRejectsUnknownCategory(t *testing.T) {
	_, e := (ReportService{Reports: &fakeReports{}}).Create(1, CreateReportInput{Title: "x", Category: "noise", Severity: 2, Latitude: -1, Longitude: 34})
	if e == nil {
		t.Fatal("expected category error")
	}
}
func TestReportServiceLikeIsIdempotentAndUnlikeRemoves(t *testing.T) {
	f := &fakeReports{}
	s := ReportService{Reports: f}
	if e := s.Like(2, 3); e != nil || f.like == nil {
		t.Fatal("expected like")
	}
	if e := s.Like(2, 3); e != nil {
		t.Fatal("expected idempotent like")
	}
	if e := s.Unlike(2, 3); e != nil || !f.deleted {
		t.Fatal("expected unlike")
	}
}
func TestReportServiceUpdateRequiresOwnership(t *testing.T) {
	f := &fakeReports{report: &models.Report{ID: 5, UserID: 2, Latitude: -1, Longitude: 34}}
	s := ReportService{Reports: f}
	title := "updated"
	if _, e := s.Update(5, 2, "public", UpdateReportInput{Title: &title}); e != nil {
		t.Fatalf("owner should update: %v", e)
	}
	if _, e := s.Update(5, 99, "public", UpdateReportInput{Title: &title}); !errors.Is(e, ErrForbidden) {
		t.Fatalf("non-owner = %v, want forbidden", e)
	}
}
func TestReportServiceDeleteAllowsAdmin(t *testing.T) {
	f := &fakeReports{report: &models.Report{ID: 5, UserID: 2}}
	s := ReportService{Reports: f}
	if e := s.Delete(5, 99, "admin"); e != nil {
		t.Fatalf("admin should delete: %v", e)
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
