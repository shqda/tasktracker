package PostgreSQL

import (
	"database/sql"
	"fmt"
	"github.com/jmoiron/sqlx"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
	"github.com/stretchr/testify/require"
	"os"
	"testing"
	"time"
)

func setupDB(t *testing.T) *sqlx.DB {
	t.Helper()

	db := setupTestDB(t)

	execSQLFile(t, db, "./testdata/schema.sql")
	execSQLFile(t, db, "./testdata/cleanup.sql")

	return db
}

func setupTestDB(t *testing.T) *sqlx.DB {
	t.Helper()

	err := godotenv.Load("testdata/.env")
	require.NoError(t, err)

	host := os.Getenv("POSTGRES_HOST")
	port := os.Getenv("POSTGRES_PORT")
	user := os.Getenv("POSTGRES_USER")
	password := os.Getenv("POSTGRES_PASSWORD")

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

func execSQLFile(t *testing.T, db *sqlx.DB, path string) {
	t.Helper()

	content, err := os.ReadFile(path)

	if err != nil {
		t.Fatal("failed to read sql file: ", err)
	}

	if _, err := db.Exec(string(content)); err != nil {
		t.Fatal("failed to exec sql file: ", err)
	}
}

func insertTask(t *testing.T, s *PostgresDB, title string) {
	t.Helper()

	if _, err := s.InsertTask(title); err != nil {
		t.Fatal("failed to insert task: ", err)
	}
}

func TestPostgresDB_GetLastTask_EmptyTable(t *testing.T) {
	db := setupDB(t)
	defer db.Close()

	storage := &PostgresDB{DB: db}

	_, err := storage.GetLastTask()

	require.ErrorIs(t, err, sql.ErrNoRows)
}

func TestPostgresDB_GetLastTask_WithData(t *testing.T) {
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
	db := setupDB(t)
	defer db.Close()

	storage := &PostgresDB{DB: db}

	id, err := storage.InsertTask("test task")

	require.NoError(t, err)
	require.Equal(t, 1, id)
}
