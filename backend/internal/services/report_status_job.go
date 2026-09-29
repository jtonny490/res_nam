package services

import (
	"context"
	"time"
)

const staleReportAge = 7 * 24 * time.Hour

type ReportStatusStore interface {
	MarkStaleBefore(time.Time) (int64, error)
}

type ReportStatusJob struct {
	Reports ReportStatusStore
	Now     func() time.Time
}

func (j ReportStatusJob) Run() (int64, error) {
	now := time.Now()
	if j.Now != nil {
		now = j.Now()
	}
	return j.Reports.MarkStaleBefore(now.Add(-staleReportAge))
}

func (j ReportStatusJob) RunPeriodically(ctx context.Context, interval time.Duration, onError func(error)) {
	if _, err := j.Run(); err != nil && onError != nil {
		onError(err)
	}
	ticker := time.NewTicker(interval)
	defer ticker.Stop()
	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			if _, err := j.Run(); err != nil && onError != nil {
				onError(err)
			}
		}
	}
}
