package server

import (
	"github.com/gin-gonic/gin"
	"github.com/go-playground/assert/v2"
	"net/http"
	"net/http/httptest"
	"testing"
)

type MockTaskHandler struct{}

func (m *MockTaskHandler) GetLastTask(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"task": "last"})
}

func (m *MockTaskHandler) PostTask(c *gin.Context) {
	c.JSON(http.StatusCreated, gin.H{"status": "created"})
}

func TestRouter_RegisterRoutes(t *testing.T) {
	gin.SetMode(gin.TestMode)

	tests := []struct {
		name     string
		method   string
		url      string
		wantCode int
	}{
		{
			name:     "Get last task",
			method:   http.MethodGet,
			url:      "/tasks/last",
			wantCode: http.StatusOK,
		},
		{
			name:     "Post new task",
			method:   http.MethodPost,
			url:      "/tasks",
			wantCode: http.StatusCreated,
		},
		{
			name:     "Unknown route",
			method:   http.MethodPost,
			url:      "/tasks/unknown",
			wantCode: http.StatusNotFound,
		},
	}
	for _, tc := range tests {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			mockHandler := new(MockTaskHandler)
			r := NewRouter(nil, mockHandler)
			r.RegisterRoutes()

			req := httptest.NewRequest(tc.method, tc.url, nil)
			w := httptest.NewRecorder()

			r.Engine.ServeHTTP(w, req)

			assert.Equal(t, tc.wantCode, w.Code)
		})
	}
}
