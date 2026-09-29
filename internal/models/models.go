package models

import "time"

type User struct {
	ID           uint      `gorm:"primaryKey" json:"id"`
	Name         string    `gorm:"not null" json:"name"`
	Email        string    `gorm:"uniqueIndex;not null" json:"email"`
	PasswordHash string    `gorm:"not null" json:"-"`
	Role         string    `gorm:"not null;default:'public'" json:"role"`
	Status       string    `gorm:"not null;default:'active'" json:"status"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}
type Report struct {
	ID             uint      `gorm:"primaryKey" json:"id"`
	UserID         uint      `gorm:"not null;index" json:"user_id"`
	Title          string    `gorm:"not null" json:"title"`
	Description    string    `gorm:"type:text" json:"description"`
	PhotoURL       string    `json:"photo_url"`
	Category       string    `gorm:"not null" json:"category"`
	Severity       int       `gorm:"not null" json:"severity"`
	Latitude       float64   `gorm:"not null" json:"latitude"`
	Longitude      float64   `gorm:"not null" json:"longitude"`
	Status         string    `gorm:"not null;default:'open'" json:"status"`
	LastActivityAt time.Time `gorm:"not null" json:"last_activity_at"`
	CreatedAt      time.Time `json:"created_at"`
	UpdatedAt      time.Time `json:"updated_at"`
	User           User      `json:"user,omitempty" gorm:"foreignKey:UserID"`
	Comments       []Comment `json:"comments,omitempty"`
	Likes          []Like    `json:"likes,omitempty"`
}
type Comment struct {
	ID                 uint      `gorm:"primaryKey" json:"id"`
	ReportID           uint      `gorm:"not null;index" json:"report_id"`
	UserID             uint      `gorm:"not null;index" json:"user_id"`
	Body               string    `gorm:"type:text;not null" json:"body"`
	IsAuthorityComment bool      `gorm:"not null;default:false" json:"is_authority_comment"`
	CreatedAt          time.Time `json:"created_at"`
	User               User      `json:"user,omitempty" gorm:"foreignKey:UserID"`
}
type Like struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	ReportID  uint      `gorm:"not null;uniqueIndex:idx_report_user" json:"report_id"`
	UserID    uint      `gorm:"not null;uniqueIndex:idx_report_user" json:"user_id"`
	CreatedAt time.Time `json:"created_at"`
}
type AuthorityRequest struct {
	ID               uint      `gorm:"primaryKey" json:"id"`
	UserID           uint      `gorm:"not null;index" json:"user_id"`
	OrganizationName string    `gorm:"not null" json:"organization_name"`
	Justification    string    `gorm:"type:text" json:"justification"`
	Status           string    `gorm:"not null;default:'pending'" json:"status"`
	ReviewedBy       *uint     `json:"reviewed_by"`
	CreatedAt        time.Time `json:"created_at"`
	UpdatedAt        time.Time `json:"updated_at"`
}

// ReportPin is a lightweight projection for the map overlay.
type ReportPin struct {
	ID        uint    `json:"id"`
	Category  string  `json:"category"`
	Severity  int     `json:"severity"`
	Status    string  `json:"status"`
	Latitude  float64 `json:"latitude"`
	Longitude float64 `json:"longitude"`
}

type CategoryCount struct {
	Category string `json:"category"`
	Count    int64  `json:"count"`
}

type StatusCount struct {
	Status string `json:"status"`
	Count  int64  `json:"count"`
}
