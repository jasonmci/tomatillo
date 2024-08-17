package main

import (
	"database/sql"
	"fmt"
	"os"
	"testing"

	_ "github.com/mattn/go-sqlite3"
)

func TestGenerateEmojis(t *testing.T) {
    tests := []struct {
        count     int
        emoji     string
        expected  string
    }{
        {0, "🍅", ""},
        {1, "🍅", "🍅"},
        {3, "🍅", "🍅🍅🍅"},
        {2, "🌱", "🌱🌱"},
        {5, "🔥", "🔥🔥🔥🔥🔥"},
    }

    for _, tt := range tests {
        result := generateEmojis(tt.count, tt.emoji)
        if result != tt.expected {
            t.Errorf("generateEmojis(%d, %q) = %q; want %q", tt.count, tt.emoji, result, tt.expected)
        }
    }
}

func setupTestDB(t *testing.T) *sql.DB {
    // Create a temporary database for testing
    db, err := sql.Open("sqlite3", "./test_tomatillo.db")
    if err != nil {
        t.Fatalf("Failed to create test database: %v", err)
    }

    // Create the tasks table
    createTableQuery := `
    CREATE TABLE IF NOT EXISTS tasks (
        id INTEGER PRIMARY KEY AUTOINCREMENT,
        name TEXT NOT NULL,
        estimate INTEGER NOT NULL,
        actual INTEGER DEFAULT 0,
        created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
        updated_at DATETIME DEFAULT CURRENT_TIMESTAMP,
        done BOOLEAN DEFAULT 0
    );`
    _, err = db.Exec(createTableQuery)
    if err != nil {
        t.Fatalf("Failed to create tasks table: %v", err)
    }

    return db
}

func teardownTestDB(db *sql.DB) {
    db.Close()
    os.Remove("./test_tomatillo.db") // Clean up by deleting the test database file
}

func TestAddTask(t *testing.T) {
    db := setupTestDB(t)
    defer teardownTestDB(db)

    // Test cases
    tests := []struct {
        name        string
        estimate    int
        expectError bool
    }{
        {"Task 1", 3, false},
        {"", 3, true}, // This should trigger an error because the name is empty
        {"Task 2", 0, false},
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            err := addTask(db, tt.name, tt.estimate)
            if tt.expectError {
                if err == nil {
                    t.Errorf("Expected an error when adding a task with an empty name, but got none")
                }
            } else {
                if err != nil {
                    t.Errorf("Did not expect an error, but got: %v", err)
                }

                // Verify that the task was added if no error
                var count int
                err := db.QueryRow(`SELECT COUNT(*) FROM tasks WHERE name = ?`, tt.name).Scan(&count)
                if err != nil {
                    t.Fatalf("Failed to query database: %v", err)
                }

                if count != 1 {
                    t.Errorf("Expected 1 task to be added, got %d", count)
                }
            }
        })
    }
}

func TestDeleteTask(t *testing.T) {
    db := setupTestDB(t)
    defer teardownTestDB(db)

    // Add some tasks to the database for testing
    tasks := []struct {
        name     string
        estimate int
    }{
        {"Task A", 3},
        {"Task B", 2},
        {"Task C", 1},
    }

    for _, task := range tasks {
        err := addTask(db, task.name, task.estimate)
        if err != nil {
            t.Fatalf("Failed to add task: %v", err)
        }
    }

    // Test cases
    tests := []struct {
        id           int
        expectError  bool
        expectedRows int
    }{
        {1, false, 2},  // Deleting an existing task (ID: 1)
        {999, true, 2}, // Attempting to delete a non-existent task (ID: 999)
        {2, false, 1},  // Deleting another existing task (ID: 2)
    }

    for _, tt := range tests {
        t.Run(fmt.Sprintf("Delete ID %d", tt.id), func(t *testing.T) {
            deleteTask(db, tt.id)

            // Verify that the task count is as expected
            var count int
            err := db.QueryRow(`SELECT COUNT(*) FROM tasks`).Scan(&count)
            if err != nil {
                t.Fatalf("Failed to query database: %v", err)
            }

            if count != tt.expectedRows {
                t.Errorf("Expected %d tasks to remain, got %d", tt.expectedRows, count)
            }

            // If expecting an error (i.e., task was not deleted), verify the message
            if tt.expectError {
                var taskName string
                err := db.QueryRow(`SELECT name FROM tasks WHERE id = ?`, tt.id).Scan(&taskName)
                if err == nil {
                    t.Errorf("Expected no task to be found with ID %d, but found one", tt.id)
                }
            }
        })
    }
}


