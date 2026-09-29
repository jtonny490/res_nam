package handlers

import (
	"bytes"
	"encoding/json"
	"errors"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	"time"

	"github.com/gin-gonic/gin"

	"res_nam/internal/models"
	"res_nam/internal/services"
)

func TestMain(m *testing.M) {
	gin.SetMode(gin.TestMode)
	os.Exit(m.Run())
}

type fakeReportStore struct {
	report  *models.Report
	created *models.Report
	updated *models.Report
	deleted bool
	pins    []models.ReportPin
}

func (f *fakeReportStore) List(string, string, int, int) ([]models.Report, int64, error) {
	return nil, 0, nil
}
func (f *fakeReportStore) Get(id uint) (*models.Report, error) {
	if f.report != nil {
		return f.report, nil
	}
	return &models.Report{ID: id, UserID: 1}, nil
}
func (f *fakeReportStore) Create(x *models.Report) error { f.created = x; return nil }
func (f *fakeReportStore) Update(x *models.Report) error { f.updated = x; return nil }
func (f *fakeReportStore) Delete(uint) error             { f.deleted = true; return nil }
func (f *fakeReportStore) AddComment(*models.Comment, time.Time) error {
	return nil
}
func (f *fakeReportStore) ListComments(uint) ([]models.Comment, error) { return nil, nil }
func (f *fakeReportStore) FindLike(uint, uint) (*models.Like, error) {
	return nil, errors.New("not found")
}
func (f *fakeReportStore) CreateLike(*models.Like) error              { return nil }
func (f *fakeReportStore) DeleteLike(*models.Like) error              { return nil }
func (f *fakeReportStore) UpdateStatus(uint, string, time.Time) error { return nil }
func (f *fakeReportStore) Pins() ([]models.ReportPin, error)          { return f.pins, nil }

func reportContext(method, target string, body []byte, uid uint, role string) (*gin.Context, *httptest.ResponseRecorder) {
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest(method, target, bytes.NewReader(body))
	c.Request.Header.Set("Content-Type", "application/json")
	if uid != 0 {
		c.Set("user_id", uid)
	}
	if role != "" {
		c.Set("role", role)
	}
	return c, w
}

func TestReportHandlerCreateJSON(t *testing.T) {
	f := &fakeReportStore{}
	h := ReportHandler{S: services.ReportService{Reports: f}}
	body, _ := json.Marshal(map[string]any{"title": "Spill", "category": "water", "severity": 3, "latitude": -1.0, "longitude": 34.0})
	c, w := reportContext(http.MethodPost, "/api/reports", body, 7, "public")

	h.Create(c)

	if w.Code != http.StatusCreated {
		t.Fatalf("status = %d body=%s", w.Code, w.Body.String())
	}
	if f.created == nil || f.created.UserID != 7 {
		t.Fatalf("report not created: %#v", f.created)
	}
}

func TestReportHandlerCreateMultipartWithPhoto(t *testing.T) {
	f := &fakeReportStore{}
	h := ReportHandler{S: services.ReportService{Reports: f}}

	var buf bytes.Buffer
	mw := multipart.NewWriter(&buf)
	for k, v := range map[string]string{"title": "Oil", "category": "water", "severity": "4", "latitude": "-1", "longitude": "34"} {
		_ = mw.WriteField(k, v)
	}
	part, _ := mw.CreateFormFile("photo", "spill.png")
	part.Write([]byte("fake-image"))
	mw.Close()

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest(http.MethodPost, "/api/reports", &buf)
	c.Request.Header.Set("Content-Type", mw.FormDataContentType())
	c.Set("user_id", uint(7))
	c.Set("role", "public")

	h.Create(c)
	defer os.RemoveAll("uploads")

	if w.Code != http.StatusCreated {
		t.Fatalf("status = %d body=%s", w.Code, w.Body.String())
	}
	if f.created == nil || f.created.PhotoURL == "" {
		t.Fatalf("photo not stored: %#v", f.created)
	}
}

func TestReportHandlerUpdateForbidden(t *testing.T) {
	f := &fakeReportStore{report: &models.Report{ID: 5, UserID: 2}}
	h := ReportHandler{S: services.ReportService{Reports: f}}
	body, _ := json.Marshal(map[string]any{"title": "x"})
	c, w := reportContext(http.MethodPatch, "/api/reports/5", body, 99, "public")
	c.Params = gin.Params{{Key: "id", Value: "5"}}

	h.Update(c)

	if w.Code != http.StatusForbidden {
		t.Fatalf("status = %d body=%s", w.Code, w.Body.String())
	}
}

func TestReportHandlerDeleteNoContent(t *testing.T) {
	f := &fakeReportStore{report: &models.Report{ID: 5, UserID: 2}}
	h := ReportHandler{S: services.ReportService{Reports: f}}
	c, _ := reportContext(http.MethodDelete, "/api/reports/5", nil, 2, "public")
	c.Params = gin.Params{{Key: "id", Value: "5"}}

	h.Delete(c)

	if c.Writer.Status() != http.StatusNoContent || !f.deleted {
		t.Fatalf("status = %d deleted=%v", c.Writer.Status(), f.deleted)
	}
}

type fakeUserStore struct{ user *models.User }

func (f *fakeUserStore) Create(u *models.User) error { u.ID = 1; f.user = u; return nil }
func (f *fakeUserStore) FindByEmail(string) (*models.User, error) {
	if f.user == nil {
		return nil, errors.New("not found")
	}
	return f.user, nil
}

func TestAuthHandlerRegisterAndLogin(t *testing.T) {
	store := &fakeUserStore{}
	h := AuthHandler{S: &services.AuthService{Users: store, Secret: "test"}}
	body, _ := json.Marshal(map[string]string{"name": "A User", "email": "a@example.com", "password": "password"})
	c, w := reportContext(http.MethodPost, "/api/auth/register", body, 0, "")

	h.Register(c)
	if w.Code != http.StatusCreated {
		t.Fatalf("register status = %d body=%s", w.Code, w.Body.String())
	}

	bad, _ := json.Marshal(map[string]string{"email": "a@example.com", "password": "wrong"})
	c, w = reportContext(http.MethodPost, "/api/auth/login", bad, 0, "")
	h.Login(c)
	if w.Code != http.StatusUnauthorized {
		t.Fatalf("login status = %d, want 401", w.Code)
	}
}

type fakeAdminUserStore struct{ updated bool }

func (f *fakeAdminUserStore) List() ([]models.User, error) {
	return []models.User{{ID: 1, Name: "A"}}, nil
}
func (f *fakeAdminUserStore) UpdateRoleStatus(uint, string, string) error {
	f.updated = true
	return nil
}
func (f *fakeAdminUserStore) CountActive() (int64, error) { return 2, nil }

type fakeAdminReportStore struct{}

func (fakeAdminReportStore) CountAll() (int64, error) { return 3, nil }
func (fakeAdminReportStore) CountByCategory() ([]models.CategoryCount, error) {
	return nil, nil
}
func (fakeAdminReportStore) CountByStatus() ([]models.StatusCount, error) { return nil, nil }

func TestAdminHandlerUsersAndUpdate(t *testing.T) {
	users := &fakeAdminUserStore{}
	h := AdminHandler{S: services.AdminService{UserStore: users, ReportStore: fakeAdminReportStore{}}}

	c, w := reportContext(http.MethodGet, "/api/admin/users", nil, 1, "admin")
	h.Users(c)
	if w.Code != http.StatusOK {
		t.Fatalf("users status = %d", w.Code)
	}

	body, _ := json.Marshal(map[string]string{"role": "authority", "status": "active"})
	c, w = reportContext(http.MethodPatch, "/api/admin/users/1", body, 1, "admin")
	c.Params = gin.Params{{Key: "id", Value: "1"}}
	h.UpdateUser(c)
	if w.Code != http.StatusOK || !users.updated {
		t.Fatalf("update status = %d updated=%v", w.Code, users.updated)
	}
}

func TestAdminHandlerAnalytics(t *testing.T) {
	h := AdminHandler{S: services.AdminService{UserStore: &fakeAdminUserStore{}, ReportStore: fakeAdminReportStore{}}}
	c, w := reportContext(http.MethodGet, "/api/admin/analytics", nil, 1, "admin")
	h.Analytics(c)
	if w.Code != http.StatusOK {
		t.Fatalf("status = %d body=%s", w.Code, w.Body.String())
	}
}
