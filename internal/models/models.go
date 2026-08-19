package models

import "time"

type User struct {
	ID           uint   `gorm:"primaryKey"`
	Name         string `gorm:"not null"`
	Email        string `gorm:"uniqueIndex;not null"`
	PasswordHash string `gorm:"not null"`
	Role         string `gorm:"not null;default:'public'"`
	Status       string `gorm:"not null;default:'active'"`
	CreatedAt    time.Time
	UpdatedAt    time.Time
}
type Report struct {
	ID             uint   `gorm:"primaryKey"`
	UserID         uint   `gorm:"not null;index"`
	Title          string `gorm:"not null"`
	Description    string `gorm:"type:text"`
	PhotoURL       string
	Category       string    `gorm:"not null"`
	Severity       int       `gorm:"not null"`
	Latitude       float64   `gorm:"not null"`
	Longitude      float64   `gorm:"not null"`
	Status         string    `gorm:"not null;default:'open'"`
	LastActivityAt time.Time `gorm:"not null"`
	CreatedAt      time.Time
	UpdatedAt      time.Time
}
type Comment struct {
	ID                 uint   `gorm:"primaryKey"`
	ReportID           uint   `gorm:"not null;index"`
	UserID             uint   `gorm:"not null;index"`
	Body               string `gorm:"type:text;not null"`
	IsAuthorityComment bool   `gorm:"not null;default:false"`
	CreatedAt          time.Time
}
type Like struct {
	ID        uint `gorm:"primaryKey"`
	ReportID  uint `gorm:"not null;uniqueIndex:idx_report_user"`
	UserID    uint `gorm:"not null;uniqueIndex:idx_report_user"`
	CreatedAt time.Time
}
type AuthorityRequest struct {
	ID               uint   `gorm:"primaryKey"`
	UserID           uint   `gorm:"not null;index"`
	OrganizationName string `gorm:"not null"`
	Justification    string `gorm:"type:text"`
	Status           string `gorm:"not null;default:'pending'"`
	ReviewedBy       *uint
	CreatedAt        time.Time
	UpdatedAt        time.Time
}
