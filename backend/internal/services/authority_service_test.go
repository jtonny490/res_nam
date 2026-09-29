package services

import (
	"errors"
	"res_nam/internal/models"
	"testing"
)

type fakeAuthority struct {
	pending bool
	created *models.AuthorityRequest
	review  string
}

func (f *fakeAuthority) Create(x *models.AuthorityRequest) error { f.created = x; return nil }
func (f *fakeAuthority) PendingByUser(uint) (*models.AuthorityRequest, error) {
	if f.pending {
		return &models.AuthorityRequest{ID: 1}, nil
	}
	return nil, errors.New("not found")
}
func (f *fakeAuthority) ListPending() ([]models.AuthorityRequest, error) { return nil, nil }
func (f *fakeAuthority) Review(_, _ uint, s string) error                { f.review = s; return nil }
func TestAuthorityApply(t *testing.T) {
	f := &fakeAuthority{}
	x, e := (AuthorityService{Requests: f}).Apply(2, "  Maritime Authority ", "reason")
	if e != nil || x.Status != "pending" || x.OrganizationName != "Maritime Authority" {
		t.Fatalf("unexpected: %#v %v", x, e)
	}
}
func TestAuthorityRejectsDuplicate(t *testing.T) {
	_, e := (AuthorityService{Requests: &fakeAuthority{pending: true}}).Apply(2, "Org", "")
	if e == nil {
		t.Fatal("expected duplicate error")
	}
}
func TestAuthorityReviewStatus(t *testing.T) {
	f := &fakeAuthority{}
	if e := (AuthorityService{Requests: f}).Review(1, 9, "approved"); e != nil || f.review != "approved" {
		t.Fatal("review failed")
	}
}
