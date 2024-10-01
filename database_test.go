package main

import (
	"database/sql"
	"fmt"
	"time"
	"testing"

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

// Helper function to insert test data
func insertTestTask(db *sql.DB, name string, estimate, actual int, done bool) {
    query := `
    INSERT INTO tasks (name, estimate, actual, done) 
    VALUES (?, ?, ?, ?)`
    _, err := db.Exec(query, name, estimate, actual, done)
    if err != nil {
        fmt.Printf("failed to insert test task: %v", err)
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

func TestActivateTask(t *testing.T) {
    setupTestDB2(t) // Initialize the test database

    db = testDB

    mockDate := time.Now().Format("2006-01-02")
    mockHalfHour := getHalfHour(time.Now().Hour(), time.Now().Minute())

    fmt.Println(mockDate)
    fmt.Println(mockHalfHour)

    insertTrackingTask(1, mockDate, mockHalfHour)

    if !taskExists(t, 1, mockDate, mockHalfHour, "active") {
        t.Errorf("expected task to be 'active', but it wasn't found")
    }

    insertTrackingTask(1, mockDate, mockHalfHour)

    if !taskExists(t, 1, mockDate, mockHalfHour, "done") {
        t.Errorf("expected task to be updated to 'done', but it wasn't")
    }
}

func TestInsertTrackingTask(t *testing.T) {
    setupTestDB2(t) // Initialize the test database

    db = testDB

    mockDate := time.Now().Format("2006-01-02")
    mockHalfHour := getHalfHour(time.Now().Hour(), time.Now().Minute())

    insertTrackingTask(1, mockDate, mockHalfHour)

    if !taskExists(t, 1, mockDate, mockHalfHour, "active") {
        t.Errorf("expected task to be 'active', but it wasn't found")
    }

    insertTrackingTask(1, mockDate, mockHalfHour)

    if !taskExists(t, 1, mockDate, mockHalfHour, "done") {
        t.Errorf("expected task to be updated to 'done', but it wasn't")
    }
}

// add test for getDailyTasks
func TestGetDailyTasks(t *testing.T) {
    setupTestDB2(t) // Initialize the test database

    db = testDB

    err := addTask(db, "Task 1", 1)
    if err != nil {
        t.Fatalf("failed to add task: %v", err)
    }

    tasks, _ := getDailyTasks(db)
    if len(tasks) != 1 {
        t.Errorf("expected 1 task, got %d", len(tasks))
    }
}

func TestMarkAsDone(t *testing.T) {
    setupTestDB2(t) // Initialize the test database
    defer db.Close()
    db = testDB

    err := addTask(db, "Task 1", 1)
    if err != nil {
        t.Fatalf("failed to add task: %v", err)
    }

    err = markAsDone(1)
    if err != nil {
        t.Fatalf("failed to mark task as done: %v", err)
    }


    // Verify the task status
    var done bool
    query := `SELECT done FROM tasks WHERE id = ?`
    err = db.QueryRow(query, 1).Scan(&done)
    if err != nil {
        t.Fatalf("failed to query task: %v", err)
    }

    if !done {
        t.Errorf("expected task to be marked as done, but it wasn't")
    }
}

func TestGetTasks(t *testing.T) {
    setupTestDB2(t) // Initialize the test database
    db = testDB

    insertTestTask(db, "Task 1", 5, 2, false) // WIP task
    insertTestTask(db, "Task 2", 3, 0, false) // To Do task
    insertTestTask(db, "Task 3", 4, 4, true)  // Done task

    // Test case 1: Status "all"
    tasks, err := getTasks(7, "all")
    if err != nil {
        t.Fatalf("unexpected error: %v", err)
    }
    if len(tasks) != 3 {
        t.Errorf("expected 3 tasks, got %d", len(tasks))
    }

    // Test case 2: Status "wip"
    tasks, err = getTasks(7, "wip")
    if err != nil {
        t.Fatalf("unexpected error: %v", err)
    }
    if len(tasks) != 1 {
        t.Errorf("expected 1 WIP task, got %d", len(tasks))
    }
    if tasks[0].Name != "Task 1" {
        t.Errorf("expected Task 1 to be WIP, got %s", tasks[0].Name)
    }

    // Test case 3: Status "todo"
    tasks, err = getTasks(7, "todo")
    if err != nil {
        t.Fatalf("unexpected error: %v", err)
    }
    if len(tasks) != 1 {
        t.Errorf("expected 1 To Do task, got %d", len(tasks))
    }
    if tasks[0].Name != "Task 2" {
        t.Errorf("expected Task 2 to be To Do, got %s", tasks[0].Name)
    }

    // Test case 4: Status "done"
    tasks, err = getTasks(7, "done")
    if err != nil {
        t.Fatalf("unexpected error: %v", err)
    }
    if len(tasks) != 1 {
        t.Errorf("expected 1 Done task, got %d", len(tasks))
    }
    if tasks[0].Name != "Task 3" {
        t.Errorf("expected Task 3 to be Done, got %s", tasks[0].Name)
    }

    // Test case 5: Invalid status
    _, err = getTasks(7, "invalid")
    if err == nil {
        t.Error("expected error for invalid status, but got none")
    }

}