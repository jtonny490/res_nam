package services

import (
	"errors"
	"res_nam/internal/models"
	"testing"
)

type fakeUsers struct{ u *models.User }

func (f *fakeUsers) Create(u *models.User) error { u.ID = 1; f.u = u; return nil }
func (f *fakeUsers) FindByEmail(string) (*models.User, error) {
	if f.u == nil {
		return nil, errors.New("not found")
	}
	return f.u, nil
}
func TestAuthRegisterAndLogin(t *testing.T) {
	f := &fakeUsers{}
	s := AuthService{Users: f, Secret: "test-secret"}
	u, token, e := s.Register("A User", "a@example.com", "password")
	if e != nil || u.Role != "public" || token == "" {
		t.Fatalf("register failed: %v", e)
	}
	if _, _, e = s.Login("a@example.com", "password"); e != nil {
		t.Fatalf("login failed: %v", e)
	}
}
func TestAuthRejectsWeakPassword(t *testing.T) {
	s := AuthService{Users: &fakeUsers{}}
	_, _, e := s.Register("A", "a@example.com", "short")
	if e == nil {
		t.Fatal("expected validation error")
	}
}
