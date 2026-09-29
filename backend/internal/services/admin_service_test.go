package services

import (
	"res_nam/internal/models"
	"testing"
)

type fakeAdminUsers struct {
	updated      bool
	role, status string
}

func (f *fakeAdminUsers) List() ([]models.User, error) { return nil, nil }
func (f *fakeAdminUsers) UpdateRoleStatus(_ uint, role, status string) error {
	f.updated = true
	f.role = role
	f.status = status
	return nil
}
func (f *fakeAdminUsers) CountActive() (int64, error) { return 4, nil }

type fakeAdminReports struct{}

func (f fakeAdminReports) CountAll() (int64, error) { return 10, nil }
func (f fakeAdminReports) CountByCategory() ([]models.CategoryCount, error) {
	return []models.CategoryCount{{Category: "water", Count: 6}}, nil
}
func (f fakeAdminReports) CountByStatus() ([]models.StatusCount, error) {
	return []models.StatusCount{{Status: "open", Count: 7}}, nil
}

func TestAdminUpdateUserValidates(t *testing.T) {
	f := &fakeAdminUsers{}
	s := AdminService{UserStore: f, ReportStore: fakeAdminReports{}}
	if e := s.UpdateUser(1, "authority", "active"); e != nil || !f.updated {
		t.Fatalf("valid update failed: %v", e)
	}
	if e := s.UpdateUser(1, "root", "active"); e == nil {
		t.Fatal("expected role validation error")
	}
	if e := s.UpdateUser(1, "public", "zombie"); e == nil {
		t.Fatal("expected status validation error")
	}
	if e := s.UpdateUser(0, "public", "active"); e == nil {
		t.Fatal("expected id validation error")
	}
}

func TestAdminAnalyticsAggregates(t *testing.T) {
	a, e := (AdminService{UserStore: &fakeAdminUsers{}, ReportStore: fakeAdminReports{}}).Analytics()
	if e != nil {
		t.Fatal(e)
	}
	if a.TotalReports != 10 || a.ActiveUsers != 4 || len(a.ByCategory) != 1 || len(a.ByStatus) != 1 {
		t.Fatalf("unexpected analytics: %#v", a)
	}
}
