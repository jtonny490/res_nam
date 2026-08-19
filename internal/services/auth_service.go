package services

import (
	"errors"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
	"regexp"
	"res_nam/internal/models"
	"time"
)

type UserStore interface {
	Create(*models.User) error
	FindByEmail(string) (*models.User, error)
}
type AuthService struct {
	Users  UserStore
	Secret string
}

func (s *AuthService) Register(n, e, p string) (*models.User, string, error) {
	if !regexp.MustCompile(`^[^@]+@[^@]+\.[^@]+$`).MatchString(e) || len(p) < 8 {
		return nil, "", errors.New("valid email and password of at least 8 characters required")
	}
	if _, err := s.Users.FindByEmail(e); err == nil {
		return nil, "", errors.New("email already registered")
	}
	h, err := bcrypt.GenerateFromPassword([]byte(p), 12)
	if err != nil {
		return nil, "", err
	}
	u := &models.User{Name: n, Email: e, PasswordHash: string(h), Role: "public", Status: "active"}
	if err = s.Users.Create(u); err != nil {
		return nil, "", err
	}
	return u, s.token(u)
}
func (s *AuthService) Login(e, p string) (*models.User, string, error) {
	u, err := s.Users.FindByEmail(e)
	if err != nil || bcrypt.CompareHashAndPassword([]byte(u.PasswordHash), []byte(p)) != nil {
		return nil, "", errors.New("invalid credentials")
	}
	t, err := s.token(u)
	return u, t, err
}
func (s *AuthService) token(u *models.User) (string, error) {
	return jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{"user_id": u.ID, "role": u.Role, "exp": time.Now().Add(24 * time.Hour).Unix()}).SignedString([]byte(s.Secret))
}
