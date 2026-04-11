package handlers

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"tasktracker/internal/models"
	"tasktracker/internal/services/mocks"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestTaskHandler_GetLastTask(t *testing.T) {
	gin.SetMode(gin.TestMode)

	method := "LastTask"
	url := "/tasks/last"

	tests := []struct {
		name      string
		setupMock func(m *mocks.MockTaskService)
		wantCode  int
		wantId    int32
	}{
		{
			name: "success",
			setupMock: func(m *mocks.MockTaskService) {
				m.
					On(method, mock.Anything).
					Once().
					Return(&models.Task{ID: 10}, nil)
			},
			wantCode: http.StatusOK,
			wantId:   10,
		},
		{
			name: "service error",
			setupMock: func(m *mocks.MockTaskService) {
				m.
					On(method, mock.Anything).
					Once().
					Return(&models.Task{}, errors.New("service error"))
			},
			wantCode: http.StatusNotFound,
		},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			mockService := mocks.NewMockTaskService(t)
			tc.setupMock(mockService)

			h := NewTaskHandler(mockService)

			rec := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(rec)
			c.Request = httptest.NewRequest(http.MethodGet, url, nil)

			h.GetLastTask(c)

			require.Equal(t, tc.wantCode, rec.Code)

			var resp models.Task
			err := json.Unmarshal(rec.Body.Bytes(), &resp)
			require.NoError(t, err)

			assert.Equal(t, tc.wantId, resp.ID)

			mockService.AssertExpectations(t)
		})
	}
}

func TestTaskHandler_PostTask(t *testing.T) {
	gin.SetMode(gin.TestMode)

	method := "CreateTask"
	url := "/tasks"

	tests := []struct {
		name      string
		body      string
		setupMock func(m *mocks.MockTaskService)
		wantCode  int
		wantBody  string
	}{
		{
			name: "success",
			body: `{"task":"blabla"}`,
			setupMock: func(m *mocks.MockTaskService) {
				m.
					On(method, mock.AnythingOfType("string")).
					Once().
					Return(&models.Task{ID: 5, Title: "blabla"}, nil)
			},
			wantCode: http.StatusCreated,
			wantBody: `{"id":5,"title":"blabla"}`,
		},
		{
			name:      "invalid JSON",
			body:      `{"field":"blabla"}`,
			setupMock: func(m *mocks.MockTaskService) {},
			wantCode:  http.StatusBadRequest,
			wantBody:  `{"error":"invalid JSON"}`,
		},
	}
	for _, tc := range tests {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			mockService := mocks.NewMockTaskService(t)
			tc.setupMock(mockService)

			h := NewTaskHandler(mockService)

			rec := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(rec)
			c.Request = httptest.NewRequest(http.MethodPost, url, strings.NewReader(tc.body))
			c.Request.Header.Set("Content-Type", "application/json")

			h.PostTask(c)

			require.Equal(t, tc.wantCode, rec.Code)
			assert.JSONEq(t, tc.wantBody, rec.Body.String())

			mockService.AssertExpectations(t)
		})
	}
}
