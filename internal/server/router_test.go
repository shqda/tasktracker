package server

import (
	"net/http"
	"net/http/httptest"
	"tasktracker/internal/server/mocks"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestRouter_RegisterRoutes(t *testing.T) {
	gin.SetMode(gin.TestMode)

	tests := []struct {
		name      string
		method    string
		url       string
		setupMock func(m *server.MockTaskHandlerInterface)
		wantCode  int
	}{
		{
			name:   "Get last task",
			method: http.MethodGet,
			url:    "/tasks/last",
			setupMock: func(m *server.MockTaskHandlerInterface) {
				m.
					On("GetLastTask", mock.Anything).
					Once().
					Return()
			},
			wantCode: http.StatusOK,
		},
		{
			name:   "Post new task",
			method: http.MethodPost,
			url:    "/tasks",
			setupMock: func(m *server.MockTaskHandlerInterface) {
				m.
					On("PostTask", mock.Anything).
					Once().
					Return()
			},
			wantCode: http.StatusOK,
		},
		{
			name:      "Unknown route",
			method:    http.MethodPost,
			url:       "/tasks/unknown",
			setupMock: func(m *server.MockTaskHandlerInterface) {},
			wantCode:  http.StatusNotFound,
		},
	}
	for _, tc := range tests {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			mockHandler := server.NewMockTaskHandlerInterface(t)
			tc.setupMock(mockHandler)

			r := NewRouter(nil, mockHandler)
			r.RegisterRoutes()

			req := httptest.NewRequest(tc.method, tc.url, nil)
			w := httptest.NewRecorder()

			r.Engine.ServeHTTP(w, req)

			assert.Equal(t, tc.wantCode, w.Code)
			mockHandler.AssertExpectations(t)
		})
	}
}
