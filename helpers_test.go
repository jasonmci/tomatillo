package main

import (
	"database/sql"
	"fmt"
	"os"
	"testing"
	"time"

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

func TestGetMonthAbbreviation(t *testing.T) {
    tests := []struct {
        month    time.Month
        expected string
    }{
        {time.January, "Jan"},
        {time.December, "Dec"},
    }

    for _, tt := range tests {
        result := getMonthAbbreviation(tt.month)
        if result != tt.expected {
            t.Errorf("getMonthAbbreviation(%v) = %q; want %q", tt.month, result, tt.expected)
        }
    }
}

func TestHalfHour(t *testing.T) {
    tests := []struct {
        hour    int
        minute  int
        expected int
    }{
        {12, 1, 24},
        {23, 34, 47},
        {9, 45, 19},
        {0, 30, 1},
        {0, 0, 0},
    }

    for _, tt := range tests {
        result := getHalfHour(tt.hour, tt.minute)
        if result != tt.expected {
            t.Errorf("getHalfHour(%d, %d) = %d; want %d", tt.hour, tt.minute, result, tt.expected)
        }
    }
}

func TestGetWeek(t *testing.T) {
    // Mock the current date to be a Wednesday, September 20, 2024
    mockDate := time.Date(2024, time.September, 20, 0, 0, 0, 0, time.Local)
    sunday, saturday := getWeek(mockDate)

    expectedSunday := time.Date(2024, time.September, 15, 0, 0, 0, 0, time.Local)
    expectedSaturday := time.Date(2024, time.September, 21, 0, 0, 0, 0, time.Local)

    // Check if the calculated Sunday and Saturday match the expected values
    if !sunday.Equal(expectedSunday) {
        t.Errorf("expected Sunday to be %v, but got %v", expectedSunday, sunday)
    }
    if !saturday.Equal(expectedSaturday) {
        t.Errorf("expected Saturday to be %v, but got %v", expectedSaturday, saturday)
    }
}

func TestGetMonth(t *testing.T) {
    // Mock the current date to be a Wednesday, September 20, 2024
    mockDate := time.Date(2024, time.September, 20, 0, 0, 0, 0, time.Local)
    firstOfMonth, lastOfMonth := getMonth(mockDate)

    expectedFirstOfMonth := time.Date(2024, time.September, 1, 0, 0, 0, 0, time.Local)
    expectedLastOfMonth := time.Date(2024, time.September, 30, 0, 0, 0, 0, time.Local)

    // Check if the calculated first and last days of the month match the expected values
    if !firstOfMonth.Equal(expectedFirstOfMonth) {
        t.Errorf("expected first day of month to be %v, but got %v", expectedFirstOfMonth, firstOfMonth)
    }
    if !lastOfMonth.Equal(expectedLastOfMonth) {
        t.Errorf("expected last day of month to be %v, but got %v", expectedLastOfMonth, lastOfMonth)
    }
}

func TestGetCurrentYear(t *testing.T) {

}

func TestGetAllDaysOfMonth(t *testing.T) {
    // Mock the current date to be a Wednesday, September 20, 2024
    mockDate := time.Date(2024, time.February, 20, 0, 0, 0, 0, time.Local)
    days := getAllDaysOfMonth(mockDate)

    expectedDays := []string{
        "2024-02-01", "2024-02-02", "2024-02-03", "2024-02-04", "2024-02-05", "2024-02-06", "2024-02-07",
        "2024-02-08", "2024-02-09", "2024-02-10", "2024-02-11", "2024-02-12", "2024-02-13", "2024-02-14",
        "2024-02-15", "2024-02-16", "2024-02-17", "2024-02-18", "2024-02-19", "2024-02-20", "2024-02-21",
        "2024-02-22", "2024-02-23", "2024-02-24", "2024-02-25", "2024-02-26", "2024-02-27", "2024-02-28",
        "2024-02-29",
    }

    // Check if the calculated days of the month match the expected values
    for i, day := range days {
        if day != expectedDays[i] {
            t.Errorf("expected day %d to be %q, but got %q", i+1, expectedDays[i], day)
        }
    }
}