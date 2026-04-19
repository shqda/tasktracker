package handler

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"tasktracker/internal/errs"
	"tasktracker/internal/model"
	"tasktracker/internal/server/handler/mocks"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func init() {
	gin.SetMode(gin.TestMode)
}

var ErrDatabaseDown = errors.New("db down")

func TestTaskHandler_GetLastTask(t *testing.T) {

	method := "GetLastTask"
	url := "/tasks/last"

	tests := []struct {
		name      string
		setupMock func(m *handler.MockTaskServiceInterface)
		wantCode  int
		wantBody  string
	}{
		{
			name: "valid",
			setupMock: func(m *handler.MockTaskServiceInterface) {
				m.
					On(method).
					Once().
					Return(&model.Task{ID: 10}, nil)
			},
			wantCode: http.StatusOK,
			wantBody: `{"id":10,"title":""}`,
		},
		{
			name: "task not found",
			setupMock: func(m *handler.MockTaskServiceInterface) {
				m.
					On(method).
					Once().
					Return(nil, errs.ErrTaskNotFound)
			},
			wantCode: http.StatusNotFound,
			wantBody: `{"error":"task not found"}`,
		},
		{
			name: "service error",
			setupMock: func(m *handler.MockTaskServiceInterface) {
				m.
					On(method).
					Once().
					Return(nil, ErrDatabaseDown)
			},
			wantCode: http.StatusInternalServerError,
			wantBody: `{"error":"internal error"}`,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			mockService := handler.NewMockTaskServiceInterface(t)
			tt.setupMock(mockService)

			h := NewTaskHandler(mockService)

			rec := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(rec)
			c.Request = httptest.NewRequest(http.MethodGet, url, nil)

			h.GetLastTask(c)

			require.Equal(t, tt.wantCode, rec.Code)
			assert.JSONEq(t, tt.wantBody, rec.Body.String())

			mockService.AssertExpectations(t)
		})
	}
}

func TestTaskHandler_PostTask(t *testing.T) {
	method := "CreateTask"
	url := "/tasks"

	tests := []struct {
		name      string
		body      string
		setupMock func(m *handler.MockTaskServiceInterface)
		wantCode  int
		wantBody  string
	}{
		{
			name: "valid",
			body: `{"task":"blabla"}`,
			setupMock: func(m *handler.MockTaskServiceInterface) {
				m.
					On(method, "blabla").
					Once().
					Return(&model.Task{ID: 5, Title: "blabla"}, nil)
			},
			wantCode: http.StatusCreated,
			wantBody: `{"id":5,"title":"blabla"}`,
		},
		{
			name:      "invalid JSON",
			body:      `{"field":"blabla"}`,
			setupMock: func(m *handler.MockTaskServiceInterface) {},
			wantCode:  http.StatusBadRequest,
			wantBody:  `{"error":"invalid JSON"}`,
		},
		{
			name: "service error",
			body: `{"task":"blabla"}`,
			setupMock: func(m *handler.MockTaskServiceInterface) {
				m.
					On(method, "blabla").
					Once().
					Return(nil, ErrDatabaseDown)
			},
			wantCode: http.StatusInternalServerError,
			wantBody: `{"error":"internal error"}`,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			mockService := handler.NewMockTaskServiceInterface(t)
			tt.setupMock(mockService)

			h := NewTaskHandler(mockService)

			rec := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(rec)
			c.Request = httptest.NewRequest(http.MethodPost, url, strings.NewReader(tt.body))
			c.Request.Header.Set("Content-Type", "application/json")

			h.PostTask(c)

			require.Equal(t, tt.wantCode, rec.Code)
			assert.JSONEq(t, tt.wantBody, rec.Body.String())

			mockService.AssertExpectations(t)
		})
	}
}

func TestTaskHandler_GetTaskByID(t *testing.T) {
	method := "GetTaskByID"
	url := "/tasks/:id"

	tests := []struct {
		name      string
		id        string
		setupMock func(m *handler.MockTaskServiceInterface)
		wantCode  int
		wantBody  string
	}{
		{
			name: "valid",
			id:   "10",
			setupMock: func(m *handler.MockTaskServiceInterface) {
				m.
					On(method, 10).
					Once().
					Return(&model.Task{ID: 10}, nil)
			},
			wantCode: http.StatusOK,
			wantBody: `{"id":10,"title":""}`,
		},
		{
			name:      "invalid id",
			id:        "abc",
			setupMock: func(m *handler.MockTaskServiceInterface) {},
			wantCode:  http.StatusBadRequest,
			wantBody:  `{"error":"invalid id"}`,
		},
		{
			name: "task not found",
			id:   "10",
			setupMock: func(m *handler.MockTaskServiceInterface) {
				m.
					On(method, 10).
					Once().
					Return(nil, errs.ErrTaskNotFound)
			},
			wantCode: http.StatusNotFound,
			wantBody: `{"error":"task not found"}`,
		},
		{
			name: "service error",
			id:   "10",
			setupMock: func(m *handler.MockTaskServiceInterface) {
				m.
					On(method, 10).
					Once().
					Return(nil, ErrDatabaseDown)
			},
			wantCode: http.StatusInternalServerError,
			wantBody: `{"error":"internal error"}`,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			mockService := handler.NewMockTaskServiceInterface(t)
			tt.setupMock(mockService)

			h := NewTaskHandler(mockService)

			rec := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(rec)
			c.Request = httptest.NewRequest(http.MethodGet, url, nil)
			c.Params = gin.Params{{Key: "id", Value: tt.id}}

			h.GetTaskByID(c)

			require.Equal(t, tt.wantCode, rec.Code)
			assert.JSONEq(t, tt.wantBody, rec.Body.String())

			mockService.AssertExpectations(t)
		})
	}
}

func TestTaskHandler_GetAllTasks(t *testing.T) {
	method := "GetAllTasks"
	url := "/tasks/"

	tests := []struct {
		name      string
		setupMock func(m *handler.MockTaskServiceInterface)
		wantCode  int
		wantBody  string
	}{
		{
			name: "valid",
			setupMock: func(m *handler.MockTaskServiceInterface) {
				m.
					On(method).
					Once().
					Return([]model.Task{
						{Title: "task1", ID: 0},
						{Title: "task2", ID: 1},
						{Title: "task3", ID: 2},
					}, nil)
			},
			wantCode: http.StatusOK,
			wantBody: `[{"id":0,"title":"task1"},{"id":1,"title":"task2"},{"id":2,"title":"task3"}]`,
		},
		{
			name: "empty task list",
			setupMock: func(m *handler.MockTaskServiceInterface) {
				m.
					On(method).
					Once().
					Return([]model.Task{}, nil)
			},
			wantCode: http.StatusOK,
			wantBody: `[]`,
		},
		{
			name: "nil tasks",
			setupMock: func(m *handler.MockTaskServiceInterface) {
				m.
					On(method).
					Once().
					Return(nil, nil)
			},
			wantCode: http.StatusOK,
			wantBody: `null`,
		},
		{
			name: "task not found",
			setupMock: func(m *handler.MockTaskServiceInterface) {
				m.
					On(method).
					Once().
					Return(nil, errs.ErrTaskNotFound)
			},
			wantCode: http.StatusNotFound,
			wantBody: `{"error":"task not found"}`,
		},
		{
			name: "service error",
			setupMock: func(m *handler.MockTaskServiceInterface) {
				m.
					On(method).
					Once().
					Return(nil, ErrDatabaseDown)
			},
			wantCode: http.StatusInternalServerError,
			wantBody: `{"error":"internal error"}`,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			mockService := handler.NewMockTaskServiceInterface(t)
			tt.setupMock(mockService)

			h := NewTaskHandler(mockService)

			rec := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(rec)
			c.Request = httptest.NewRequest(http.MethodGet, url, nil)

			h.GetAllTasks(c)

			require.Equal(t, tt.wantCode, rec.Code)
			assert.JSONEq(t, tt.wantBody, rec.Body.String())

			mockService.AssertExpectations(t)
		})
	}
}

func TestTaskHandler_DeleteTask(t *testing.T) {
	method := "DeleteTask"
	url := "/tasks/:id"

	tests := []struct {
		name      string
		deleteId  string
		setupMock func(m *handler.MockTaskServiceInterface)
		wantCode  int
		wantBody  string
	}{
		{
			name:     "valid",
			deleteId: "1",
			setupMock: func(m *handler.MockTaskServiceInterface) {
				m.
					On(method, 1).
					Once().
					Return(nil)
			},
			wantCode: http.StatusNoContent,
		},
		{
			name:      "invalid id",
			deleteId:  "abc",
			setupMock: func(m *handler.MockTaskServiceInterface) {},
			wantCode:  http.StatusBadRequest,
			wantBody:  `{"error":"invalid id"}`,
		},
		{
			name:     "task not found",
			deleteId: "1",
			setupMock: func(m *handler.MockTaskServiceInterface) {
				m.
					On(method, 1).
					Once().
					Return(errs.ErrTaskNotFound)
			},
			wantCode: http.StatusNotFound,
			wantBody: `{"error":"task not found"}`,
		},
		{
			name:     "service error",
			deleteId: "1",
			setupMock: func(m *handler.MockTaskServiceInterface) {
				m.
					On(method, 1).
					Once().
					Return(ErrDatabaseDown)
			},
			wantCode: http.StatusInternalServerError,
			wantBody: `{"error":"internal error"}`,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			mockService := handler.NewMockTaskServiceInterface(t)
			tt.setupMock(mockService)

			h := NewTaskHandler(mockService)

			rec := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(rec)
			c.Request = httptest.NewRequest(http.MethodDelete, url, nil)
			c.Params = gin.Params{{Key: "id", Value: tt.deleteId}}

			h.DeleteTask(c)

			require.Equal(t, tt.wantCode, rec.Code)
			if tt.wantBody != "" {
				assert.JSONEq(t, tt.wantBody, rec.Body.String())
			}

			mockService.AssertExpectations(t)
		})
	}
}

func TestTaskHandler_RenameTask(t *testing.T) {
	method := "RenameTask"
	url := "/tasks/:id"

	tests := []struct {
		name      string
		id        string
		body      string
		setupMock func(m *handler.MockTaskServiceInterface)
		wantCode  int
		wantBody  string
	}{
		{
			name: "valid",
			id:   "1",
			body: `{"title":"new name"}`,
			setupMock: func(m *handler.MockTaskServiceInterface) {
				m.
					On(method, 1, "new name").
					Once().
					Return(nil)
			},
			wantCode: http.StatusOK,
			wantBody: `{"title":"new name"}`,
		},
		{
			name:      "invalid id",
			id:        "abc",
			body:      `{"title":"new name"}`,
			setupMock: func(m *handler.MockTaskServiceInterface) {},
			wantCode:  http.StatusBadRequest,
			wantBody:  `{"error":"invalid id"}`,
		},
		{
			name:      "invalid body",
			id:        "1",
			body:      `{}`,
			setupMock: func(m *handler.MockTaskServiceInterface) {},
			wantCode:  http.StatusBadRequest,
			wantBody:  `{"error":"invalid JSON"}`,
		},
		{
			name: "task not found",
			id:   "1",
			body: `{"title":"new name"}`,
			setupMock: func(m *handler.MockTaskServiceInterface) {
				m.
					On(method, 1, "new name").
					Once().
					Return(errs.ErrTaskNotFound)
			},
			wantCode: http.StatusNotFound,
			wantBody: `{"error":"task not found"}`,
		},
		{
			name: "service error",
			id:   "1",
			body: `{"title":"new name"}`,
			setupMock: func(m *handler.MockTaskServiceInterface) {
				m.
					On(method, 1, "new name").
					Once().
					Return(ErrDatabaseDown)
			},
			wantCode: http.StatusInternalServerError,
			wantBody: `{"error":"internal error"}`,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			mockService := handler.NewMockTaskServiceInterface(t)
			tt.setupMock(mockService)

			h := NewTaskHandler(mockService)

			rec := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(rec)
			c.Request = httptest.NewRequest(http.MethodPut, url, strings.NewReader(tt.body))
			c.Request.Header.Set("Content-Type", "application/json")
			c.Params = gin.Params{{Key: "id", Value: tt.id}}

			h.RenameTask(c)

			require.Equal(t, tt.wantCode, rec.Code)
			assert.JSONEq(t, tt.wantBody, rec.Body.String())

			mockService.AssertExpectations(t)
		})
	}
}
