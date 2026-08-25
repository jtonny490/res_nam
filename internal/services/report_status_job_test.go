package services

import (
	"testing"
	"time"
)

type fakeReportStatusStore struct {
	cutoff time.Time
	count  int64
	err    error
}

func (f *fakeReportStatusStore) MarkStaleBefore(cutoff time.Time) (int64, error) {
	f.cutoff = cutoff
	return f.count, f.err
}

func TestReportStatusJobMarksReportsInactiveForSevenDays(t *testing.T) {
	now := time.Date(2026, 8, 25, 12, 0, 0, 0, time.UTC)
	store := &fakeReportStatusStore{count: 3}
	job := ReportStatusJob{Reports: store, Now: func() time.Time { return now }}

	count, err := job.Run()
	if err != nil {
		t.Fatal(err)
	}
	if count != 3 {
		t.Fatalf("count = %d, want 3", count)
	}
	want := now.Add(-7 * 24 * time.Hour)
	if !store.cutoff.Equal(want) {
		t.Fatalf("cutoff = %v, want %v", store.cutoff, want)
	}
}
