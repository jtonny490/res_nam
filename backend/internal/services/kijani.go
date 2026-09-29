package services

import (
	"context"
	"errors"
	"math"
	"res_nam/internal/models"
	"time"
)

type EnvironmentalAssessment struct {
	Risk       string    `json:"risk"`
	Summary    string    `json:"summary"`
	Source     string    `json:"source"`
	AssessedAt time.Time `json:"assessed_at"`
}
type KijaniProvider interface {
	Assess(context.Context, models.Report) (EnvironmentalAssessment, error)
}
type MockKijaniProvider struct{}

func (MockKijaniProvider) Assess(_ context.Context, r models.Report) (EnvironmentalAssessment, error) {
	risk := "low"
	if r.Severity >= 4 {
		risk = "high"
	} else if r.Severity >= 2 {
		risk = "moderate"
	}
	return EnvironmentalAssessment{Risk: risk, Summary: "Mock environmental assessment based on the submitted severity and Lake Victoria location.", Source: "mock", AssessedAt: time.Now()}, nil
}

type LiveKijaniProvider struct{ BaseURL, APIKey string }

func (p LiveKijaniProvider) Assess(context.Context, models.Report) (EnvironmentalAssessment, error) {
	if p.BaseURL == "" || p.APIKey == "" {
		return EnvironmentalAssessment{}, errors.New("KijaniSpace live integration is not configured")
	}
	return EnvironmentalAssessment{}, errors.New("KijaniSpace API contract has not yet been provided")
}
func InLakeVictoriaCoverage(lat, lng float64) bool {
	return lat >= -1.6 && lat <= 0.7 && lng >= 31.5 && lng <= 35.0 && !math.IsNaN(lat) && !math.IsNaN(lng)
}
