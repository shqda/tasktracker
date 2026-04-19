package postgres

import (
	"fmt"
	"os"
	"tasktracker/internal/errs"
	"testing"
	"time"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"github.com/pressly/goose/v3"
	"github.com/stretchr/testify/require"
)

const migrationDir = "../../../migrations"

func setupDB(t *testing.T) *sqlx.DB {
	t.Helper()

	db := newTestDB(t)
	err := goose.SetDialect("postgres")
	if err != nil {
		t.Fatal("failed to setup db: ", err)
	}
	err = goose.Up(db.DB, migrationDir)
	if err != nil {
		t.Fatal("failed to apply migrations: ", err)
	}
	return db
}

func newTestDB(t *testing.T) *sqlx.DB {
	t.Helper()

	host := os.Getenv("POSTGRES_HOST")
	if host == "" {
		host = "localhost"
	}

	port := os.Getenv("POSTGRES_PORT")
	if port == "" {
		port = "5432"
	}

	user := os.Getenv("POSTGRES_USER")
	if user == "" {
		user = "postgres"
	}

	password := os.Getenv("POSTGRES_PASSWORD")
	if password == "" {
		password = "postgres"
	}

	adminDSN := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=postgres sslmode=disable",
		host, port, user, password,
	)
	adminDB, err := sqlx.Connect("postgres", adminDSN)
	require.NoError(t, err, "failed to connect to admin DB")

	dbName := fmt.Sprintf("tasktracker_test_%d", time.Now().UnixNano())

	_, err = adminDB.Exec("CREATE DATABASE " + dbName)
	require.NoError(t, err, "failed to create test DB")

	testDSN := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbName,
	)
	db, err := sqlx.Connect("postgres", testDSN)
	require.NoError(t, err, "failed to connect to test DB")

	t.Cleanup(func() {
		db.Close()
		adminDB.Exec("DROP DATABASE " + dbName)
		adminDB.Close()
	})
	return db
}

func insertTask(t *testing.T, s *PostgresDB, title string) {
	t.Helper()
	if _, err := s.InsertTask(title); err != nil {
		t.Fatal("failed to insert task: ", err)
	}
}

func TestPostgresDB_GetLastTask_EmptyTable(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test")
	}
	db := setupDB(t)
	defer db.Close()
	storage := &PostgresDB{DB: db}
	_, err := storage.GetLastTask()

	require.ErrorIs(t, err, errs.ErrTaskNotFound)
}

func TestPostgresDB_GetLastTask_WithData(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test")
	}
	db := setupDB(t)
	defer db.Close()
	storage := &PostgresDB{DB: db}

	insertTask(t, storage, "first task")
	insertTask(t, storage, "second task")
	task, err := storage.GetLastTask()

	require.NoError(t, err)
	require.Equal(t, "second task", task.Title)
}

func TestPostgresDB_InsertTask(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test")
	}
	db := setupDB(t)
	defer db.Close()
	storage := &PostgresDB{DB: db}

	id, err := storage.InsertTask("test task")

	require.NoError(t, err)
	require.Equal(t, 1, id)
}

func TestPostgresDB_GetTaskByID(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test")
	}
	db := setupDB(t)
	defer db.Close()
	storage := &PostgresDB{DB: db}

	insertTask(t, storage, "test task")
	task, err := storage.GetTaskByID(1)

	require.NoError(t, err)
	require.Equal(t, "test task", task.Title)
}

func TestPostgresDB_GetAllTasks_WithData(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test")
	}
	db := setupDB(t)
	defer db.Close()
	storage := &PostgresDB{DB: db}

	insertTask(t, storage, "first task")
	insertTask(t, storage, "second task")
	tasks, err := storage.GetAllTasks()

	require.NoError(t, err)
	require.Equal(t, 2, len(tasks))
}

func TestPostgresDB_GetAllTasks_EmptyTable(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test")
	}
	db := setupDB(t)
	defer db.Close()
	storage := &PostgresDB{DB: db}

	tasks, err := storage.GetAllTasks()

	require.NoError(t, err)
	require.Empty(t, tasks)
}

func TestPostgresDB_GetTaskByID_InvalidID(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test")
	}
	db := setupDB(t)
	defer db.Close()
	storage := &PostgresDB{DB: db}

	_, err := storage.GetTaskByID(999)

	require.ErrorIs(t, err, errs.ErrTaskNotFound)
}

func TestPostgresDB_DeleteTask_ValidID(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test")
	}
	db := setupDB(t)
	defer db.Close()
	storage := &PostgresDB{DB: db}

	insertTask(t, storage, "test task")
	err := storage.DeleteTask(1)
	_, err2 := storage.GetTaskByID(1)

	require.NoError(t, err)
	require.ErrorIs(t, err2, errs.ErrTaskNotFound)
}

func TestPostgresDB_DeleteTask_InvalidID(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test")
	}
	db := setupDB(t)
	defer db.Close()
	storage := &PostgresDB{DB: db}

	insertTask(t, storage, "test task")
	err := storage.DeleteTask(2)

	require.ErrorIs(t, err, errs.ErrTaskNotFound)
}

func TestPostgresDB_UpdateTask_ValidID(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test")
	}
	db := setupDB(t)
	defer db.Close()
	storage := &PostgresDB{DB: db}

	insertTask(t, storage, "test task")
	err := storage.UpdateTask(1, "new test task")
	task, _ := storage.GetTaskByID(1)

	require.NoError(t, err)
	require.Equal(t, "new test task", task.Title)
}

func TestPostgresDB_UpdateTask_InvalidID(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test")
	}
	db := setupDB(t)
	defer db.Close()
	storage := &PostgresDB{DB: db}

	insertTask(t, storage, "test task")
	err := storage.UpdateTask(2, "new test task")

	require.ErrorIs(t, err, errs.ErrTaskNotFound)
}
