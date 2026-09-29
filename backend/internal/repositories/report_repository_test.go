package repositories

import (
	"testing"
	"time"

	"github.com/glebarez/sqlite"
	"gorm.io/gorm"

	"res_nam/internal/models"
)

func testDB(t *testing.T) *gorm.DB {
	t.Helper()
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		t.Fatalf("open db: %v", err)
	}
	sqlDB, err := db.DB()
	if err != nil {
		t.Fatalf("sql db: %v", err)
	}
	sqlDB.SetMaxOpenConns(1)
	if err := db.AutoMigrate(&models.User{}, &models.Report{}, &models.Comment{}, &models.Like{}, &models.AuthorityRequest{}); err != nil {
		t.Fatalf("migrate: %v", err)
	}
	return db
}

func seedUser(t *testing.T, db *gorm.DB) models.User {
	t.Helper()
	u := models.User{Name: "Tester", Email: "t@example.com", PasswordHash: "x", Role: "public", Status: "active"}
	if err := db.Create(&u).Error; err != nil {
		t.Fatalf("seed user: %v", err)
	}
	return u
}

func TestReportRepositoryCreateGetList(t *testing.T) {
	db := testDB(t)
	u := seedUser(t, db)
	r := ReportRepository{DB: db}
	now := time.Now()

	for _, cat := range []string{"water", "waste"} {
		report := &models.Report{UserID: u.ID, Title: cat, Category: cat, Severity: 3, Latitude: -1, Longitude: 34, Status: "open", LastActivityAt: now}
		if err := r.Create(report); err != nil {
			t.Fatalf("create: %v", err)
		}
	}

	got, err := r.Get(1)
	if err != nil || got.Title != "water" {
		t.Fatalf("get = %#v, %v", got, err)
	}

	list, total, err := r.List("waste", "", 1, 20)
	if err != nil || total != 1 || len(list) != 1 || list[0].Category != "waste" {
		t.Fatalf("list = %#v total=%d err=%v", list, total, err)
	}

	_, total, _ = r.List("", "", 1, 1)
	if total != 2 {
		t.Fatalf("total = %d, want 2", total)
	}
}

func TestReportRepositoryUpdateAndDeleteCascades(t *testing.T) {
	db := testDB(t)
	u := seedUser(t, db)
	r := ReportRepository{DB: db}
	report := &models.Report{UserID: u.ID, Title: "old", Category: "water", Severity: 3, Latitude: -1, Longitude: 34, Status: "open", LastActivityAt: time.Now()}
	if err := r.Create(report); err != nil {
		t.Fatal(err)
	}
	report.Title = "new"
	if err := r.Update(report); err != nil {
		t.Fatal(err)
	}
	if err := db.Create(&models.Comment{ReportID: report.ID, UserID: u.ID, Body: "hi"}).Error; err != nil {
		t.Fatal(err)
	}
	if err := db.Create(&models.Like{ReportID: report.ID, UserID: u.ID}).Error; err != nil {
		t.Fatal(err)
	}
	if err := r.Delete(report.ID); err != nil {
		t.Fatal(err)
	}
	var comments, likes int64
	db.Model(&models.Comment{}).Where("report_id = ?", report.ID).Count(&comments)
	db.Model(&models.Like{}).Where("report_id = ?", report.ID).Count(&likes)
	if comments != 0 || likes != 0 {
		t.Fatalf("cascade failed comments=%d likes=%d", comments, likes)
	}
}

func TestReportRepositoryAuthorityCommentTransitionsOnce(t *testing.T) {
	db := testDB(t)
	u := seedUser(t, db)
	r := ReportRepository{DB: db}
	report := &models.Report{UserID: u.ID, Title: "spill", Category: "water", Severity: 4, Latitude: -1, Longitude: 34, Status: "open", LastActivityAt: time.Now()}
	if err := r.Create(report); err != nil {
		t.Fatal(err)
	}
	if err := r.AddComment(&models.Comment{ReportID: report.ID, UserID: u.ID, Body: "public", IsAuthorityComment: false}, time.Now()); err != nil {
		t.Fatal(err)
	}
	var after models.Report
	db.First(&after, report.ID)
	if after.Status != "open" {
		t.Fatalf("public comment changed status to %q", after.Status)
	}
	if err := r.AddComment(&models.Comment{ReportID: report.ID, UserID: u.ID, Body: "official", IsAuthorityComment: true}, time.Now()); err != nil {
		t.Fatal(err)
	}
	db.First(&after, report.ID)
	if after.Status != "investigating" {
		t.Fatalf("status = %q, want investigating", after.Status)
	}
	if err := r.UpdateStatus(report.ID, "resolved", time.Now()); err != nil {
		t.Fatal(err)
	}
	if err := r.AddComment(&models.Comment{ReportID: report.ID, UserID: u.ID, Body: "again", IsAuthorityComment: true}, time.Now()); err != nil {
		t.Fatal(err)
	}
	db.First(&after, report.ID)
	if after.Status != "resolved" {
		t.Fatalf("transition should not override resolved, got %q", after.Status)
	}
}

func TestReportRepositoryMarkStaleBefore(t *testing.T) {
	db := testDB(t)
	u := seedUser(t, db)
	r := ReportRepository{DB: db}
	old := &models.Report{UserID: u.ID, Title: "old", Category: "water", Severity: 2, Latitude: -1, Longitude: 34, Status: "open", LastActivityAt: time.Now().Add(-10 * 24 * time.Hour)}
	fresh := &models.Report{UserID: u.ID, Title: "fresh", Category: "water", Severity: 2, Latitude: -1, Longitude: 34, Status: "open", LastActivityAt: time.Now()}
	resolved := &models.Report{UserID: u.ID, Title: "done", Category: "water", Severity: 2, Latitude: -1, Longitude: 34, Status: "resolved", LastActivityAt: time.Now().Add(-10 * 24 * time.Hour)}
	for _, x := range []*models.Report{old, fresh, resolved} {
		if err := r.Create(x); err != nil {
			t.Fatal(err)
		}
	}
	n, err := r.MarkStaleBefore(time.Now().Add(-7 * 24 * time.Hour))
	if err != nil || n != 1 {
		t.Fatalf("affected = %d err=%v, want 1", n, err)
	}
	var after models.Report
	db.First(&after, old.ID)
	if after.Status != "stale" {
		t.Fatalf("old status = %q", after.Status)
	}
}

func TestReportRepositoryPinsAndCounts(t *testing.T) {
	db := testDB(t)
	u := seedUser(t, db)
	r := ReportRepository{DB: db}
	for _, cat := range []string{"water", "water", "air"} {
		if err := r.Create(&models.Report{UserID: u.ID, Title: cat, Category: cat, Severity: 2, Latitude: -1, Longitude: 34, Status: "open", LastActivityAt: time.Now()}); err != nil {
			t.Fatal(err)
		}
	}
	pins, err := r.Pins()
	if err != nil || len(pins) != 3 {
		t.Fatalf("pins = %#v err=%v", pins, err)
	}
	byCat, err := r.CountByCategory()
	if err != nil || len(byCat) != 2 || byCat[1].Category != "water" || byCat[1].Count != 2 {
		t.Fatalf("byCat = %#v err=%v", byCat, err)
	}
	total, err := r.CountAll()
	if err != nil || total != 3 {
		t.Fatalf("total = %d err=%v", total, err)
	}
}

func TestUserRepositoryListAndUpdate(t *testing.T) {
	db := testDB(t)
	u := seedUser(t, db)
	r := UserRepository{DB: db}
	users, err := r.List()
	if err != nil || len(users) != 1 {
		t.Fatalf("list = %#v err=%v", users, err)
	}
	if err := r.UpdateRoleStatus(u.ID, "authority", "banned"); err != nil {
		t.Fatal(err)
	}
	var after models.User
	db.First(&after, u.ID)
	if after.Role != "authority" || after.Status != "banned" {
		t.Fatalf("user = %#v", after)
	}
	active, err := r.CountActive()
	if err != nil || active != 0 {
		t.Fatalf("active = %d err=%v", active, err)
	}
}

func TestAuthorityRepositoryReviewApprovesRole(t *testing.T) {
	db := testDB(t)
	u := seedUser(t, db)
	repo := AuthorityRepository{DB: db}
	req := &models.AuthorityRequest{UserID: u.ID, OrganizationName: "Org", Status: "pending"}
	if err := repo.Create(req); err != nil {
		t.Fatal(err)
	}
	pending, err := repo.ListPending()
	if err != nil || len(pending) != 1 {
		t.Fatalf("pending = %#v err=%v", pending, err)
	}
	if err := repo.Review(req.ID, u.ID, "approved"); err != nil {
		t.Fatal(err)
	}
	var after models.User
	db.First(&after, u.ID)
	if after.Role != "authority" {
		t.Fatalf("role = %q, want authority", after.Role)
	}
	if err := repo.Review(req.ID, u.ID, "approved"); err == nil {
		t.Fatal("expected non-pending review to fail")
	}
}
