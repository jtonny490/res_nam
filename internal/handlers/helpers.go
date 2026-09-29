package handlers

import (
	"errors"
	"fmt"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/gin-gonic/gin"

	"res_nam/internal/services"
)

func errJSON(c *gin.Context, s int, e error) { c.AbortWithStatusJSON(s, gin.H{"error": e.Error()}) }

// domainErr maps service errors onto HTTP status codes.
func domainErr(c *gin.Context, e error) {
	if errors.Is(e, services.ErrForbidden) {
		errJSON(c, http.StatusForbidden, e)
		return
	}
	errJSON(c, http.StatusBadRequest, e)
}

var allowedImageExt = map[string]bool{".jpg": true, ".jpeg": true, ".png": true, ".webp": true, ".gif": true}

func saveUpload(c *gin.Context, f *multipart.FileHeader) (string, error) {
	ext := strings.ToLower(filepath.Ext(f.Filename))
	if !allowedImageExt[ext] {
		return "", errors.New("unsupported image type")
	}
	if err := os.MkdirAll("uploads", 0755); err != nil {
		return "", err
	}
	name := fmt.Sprintf("%d%s", time.Now().UnixNano(), ext)
	if err := c.SaveUploadedFile(f, filepath.Join("uploads", name)); err != nil {
		return "", err
	}
	return "/uploads/" + name, nil
}
