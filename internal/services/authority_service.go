package services

import (
	"errors"
	"res_nam/internal/models"
	"strings"
)

type AuthorityStore interface {
	Create(*models.AuthorityRequest) error
	PendingByUser(uint) (*models.AuthorityRequest, error)
	ListPending() ([]models.AuthorityRequest, error)
	Review(uint, uint, string) error
}
type AuthorityService struct{ Requests AuthorityStore }

func (s AuthorityService) Apply(uid uint, org, just string) (*models.AuthorityRequest, error) {
	org = strings.TrimSpace(org)
	if uid == 0 || org == "" {
		return nil, errors.New("organization name is required")
	}
	if x, e := s.Requests.PendingByUser(uid); e == nil && x.ID > 0 {
		return nil, errors.New("pending authority request already exists")
	}
	x := &models.AuthorityRequest{UserID: uid, OrganizationName: org, Justification: strings.TrimSpace(just), Status: "pending"}
	return x, s.Requests.Create(x)
}
func (s AuthorityService) Pending() ([]models.AuthorityRequest, error) {
	return s.Requests.ListPending()
}
func (s AuthorityService) Review(id, admin uint, status string) error {
	if id == 0 || admin == 0 {
		return errors.New("invalid request or reviewer")
	}
	if status != "approved" && status != "rejected" {
		return errors.New("status must be approved or rejected")
	}
	return s.Requests.Review(id, admin, status)
}
