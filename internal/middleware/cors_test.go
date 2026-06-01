package middleware

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
)

func TestCORS_noOrigins_passthrough(t *testing.T) {
	gin.SetMode(gin.TestMode)
	r := setupTestRouter(CORS(nil))
	w := performRequestWithOrigin(r, "GET", "/ping", "")
	if w.Header().Get("Access-Control-Allow-Origin") != "" {
		t.Fatal("expected no CORS headers when origins empty")
	}
}

func TestCORS_allowedOrigin(t *testing.T) {
	gin.SetMode(gin.TestMode)
	r := setupTestRouter(CORS([]string{"http://localhost:5173"}))
	w := performRequestWithOrigin(r, "GET", "/ping", "http://localhost:5173")
	if w.Header().Get("Access-Control-Allow-Origin") != "http://localhost:5173" {
		t.Fatalf("unexpected ACAO: %s", w.Header().Get("Access-Control-Allow-Origin"))
	}
}

func TestCORS_preflight(t *testing.T) {
	gin.SetMode(gin.TestMode)
	r := setupTestRouter(CORS([]string{"http://localhost:5173"}))
	w := performRequestWithOrigin(r, http.MethodOptions, "/ping", "http://localhost:5173")
	if w.Code != http.StatusNoContent {
		t.Fatalf("expected 204, got %d", w.Code)
	}
}

func setupTestRouter(mw gin.HandlerFunc) *gin.Engine {
	r := gin.New()
	r.Use(mw)
	r.GET("/ping", func(c *gin.Context) { c.String(http.StatusOK, "ok") })
	return r
}

func performRequestWithOrigin(r *gin.Engine, method, path, origin string) *httptest.ResponseRecorder {
	w := httptest.NewRecorder()
	req := httptest.NewRequest(method, path, nil)
	if origin != "" {
		req.Header.Set("Origin", origin)
	}
	r.ServeHTTP(w, req)
	return w
}
