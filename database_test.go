package main

import (
	"database/sql"
	"fmt"
	"time"

	//"log"
	"testing"
	//"time"

	_ "github.com/mattn/go-sqlite3" // SQLite driver
)

// Mock database connection for tests
var testDB *sql.DB

// Setup test database (in-memory SQLite)
func setupTestDB2(t *testing.T) {
    var err error
    testDB, err = sql.Open("sqlite3", ":memory:")
    if err != nil {
        t.Fatalf("failed to open test database: %v", err)
    }

    // Create tasks table if it doesn't exist
    createTableQuery := `
    CREATE TABLE IF NOT EXISTS tasks (
        id INTEGER PRIMARY KEY AUTOINCREMENT,
        name TEXT NOT NULL,
        estimate INTEGER NOT NULL,
        actual INTEGER DEFAULT 0,
        created_at DATETIME DEFAULT (datetime('now', 'localtime')),
        updated_at DATETIME DEFAULT (datetime('now', 'localtime')),
        done BOOLEAN DEFAULT 0
    );`
    _, err = testDB.Exec(createTableQuery)
    if err != nil {
        t.Fatalf("failed to create task_tracking table: %v", err)
    }

    // Create the task_tracking table schema for testing
    query := `
    CREATE TABLE task_tracking (
        id INTEGER PRIMARY KEY AUTOINCREMENT,
        task_id INTEGER NOT NULL,
        date TEXT NOT NULL,
        half_hour INTEGER NOT NULL,
        status TEXT NOT NULL,
        UNIQUE(task_id, date, half_hour)
    );`
    _, err = testDB.Exec(query)
    if err != nil {
        t.Fatalf("failed to create task_tracking table: %v", err)
    }
}

// Helper function to check if a task exists in the database
func taskExists(t *testing.T, taskID int, date string, halfHour int, status string) bool {
    var count int
    query := `
    SELECT COUNT(*) FROM task_tracking 
    WHERE task_id = ? AND date = ? AND half_hour = ? AND status = ?`

    err := testDB.QueryRow(query, taskID, date, halfHour, status).Scan(&count)
    if err != nil {
        t.Fatalf("failed to query task: %v", err)
    }
    return count > 0
}

// Function to test activateTask
func TestActivateTask(t *testing.T) {
    setupTestDB2(t) // Initialize the test database

    db = testDB

    mockDate := time.Now().Format("2006-01-02")
    mockHalfHour := getHalfHour(time.Now().Hour(), time.Now().Minute())

    fmt.Println(mockDate)
    fmt.Println(mockHalfHour)

    insertTrackingTask(1, mockDate, mockHalfHour)

    // Check if the task was inserted with 'active' status
    if !taskExists(t, 1, mockDate, mockHalfHour, "active") {
        t.Errorf("expected task to be 'active', but it wasn't found")
    }

    // Test: Update the task to 'done' status
    insertTrackingTask(1, mockDate, mockHalfHour)

    // Check if the task was updated to 'done' status
    if !taskExists(t, 1, mockDate, mockHalfHour, "done") {
        t.Errorf("expected task to be updated to 'done', but it wasn't")
    }
}

func TestInsertTrackingTask(t *testing.T) {
    setupTestDB2(t) // Initialize the test database

    db = testDB

    mockDate := time.Now().Format("2006-01-02")
    mockHalfHour := getHalfHour(time.Now().Hour(), time.Now().Minute())

    // Test: Insert a new task with 'active' status
    insertTrackingTask(1, mockDate, mockHalfHour)

    // Check if the task was inserted with 'active' status
    if !taskExists(t, 1, mockDate, mockHalfHour, "active") {
        t.Errorf("expected task to be 'active', but it wasn't found")
    }

    // Test: Update the task to 'done' status
    insertTrackingTask(1, mockDate, mockHalfHour)

    // Check if the task was updated to 'done' status
    if !taskExists(t, 1, mockDate, mockHalfHour, "done") {
        t.Errorf("expected task to be updated to 'done', but it wasn't")
    }
}

// add test for getDailyTasks
func TestGetDailyTasks(t *testing.T) {
    setupTestDB2(t) // Initialize the test database

    db = testDB

    // Test: Insert a new task with 'active' status
    err := addTask(db, "Task 1", 1)
    if err != nil {
        t.Fatalf("failed to add task: %v", err)
    }

    // Test: Get all tasks for the day
    tasks, _ := getDailyTasks(db)
    if len(tasks) != 1 {
        t.Errorf("expected 1 task, got %d", len(tasks))
    }
}