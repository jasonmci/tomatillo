package main

import (
	"testing"
    "os"
)

func TestHandleAddCommand(t *testing.T) {
    db := setupTestDB(t)
    defer teardownTestDB(db)

    tests := []struct {
        name        string
        args        []string
        expectError bool
    }{
        {"Valid Task", []string{"--name=Task1", "--estimate=3"}, false},
        {"Missing Name", []string{"--estimate=3"}, true},
        {"Empty Args", []string{}, true},
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            err := handleAddCommand(db, tt.args)
            if tt.expectError {
                if err == nil {
                    t.Errorf("Expected error, got nil")
                }
            } else {
                if err != nil {
                    t.Errorf("Did not expect error, got %v", err)
                }

                // Verify that the task was added to the database
                if !tt.expectError {
                    var count int
                    err = db.QueryRow(`SELECT COUNT(*) FROM tasks WHERE name = ?`, "Task1").Scan(&count)
                    if err != nil {
                        t.Fatalf("Failed to query database: %v", err)
                    }

                    if count != 1 {
                        t.Errorf("Expected 1 task to be added, got %d", count)
                    }
                }
            }
        })
    }
}

func TestHandleLoadTasksCommand(t *testing.T) {
	db := setupTestDB(t)
	defer teardownTestDB(db)

	// Create a temporary file with test data
	fileContent := "Task1,3\nTask2,2\nTask3,1"
	tmpfile, err := os.CreateTemp("", "tasks_test_*.txt")
	if err != nil {
		t.Fatalf("Failed to create temp file: %v", err)
	}
	defer os.Remove(tmpfile.Name())

	if _, err := tmpfile.Write([]byte(fileContent)); err != nil {
		t.Fatalf("Failed to write to temp file: %v", err)
	}
	if err := tmpfile.Close(); err != nil {
		t.Fatalf("Failed to close temp file: %v", err)
	}

	// Run the handleLoadTasksCommand function
	err = handleLoadTasksCommand(db, []string{"--file", tmpfile.Name()})
	if err != nil {
		t.Errorf("handleLoadTasksCommand returned an error: %v", err)
	}

	// Verify that the tasks were added to the database
	var count int
	err = db.QueryRow(`SELECT COUNT(*) FROM tasks`).Scan(&count)
	if err != nil {
		t.Fatalf("Failed to query database: %v", err)
	}

	if count != 3 {
		t.Errorf("Expected 3 tasks to be loaded, got %d", count)
	}

	// Verify that specific tasks were added
	tasks := []struct {
		name     string
		estimate int
	}{
		{"Task1", 3},
		{"Task2", 2},
		{"Task3", 1},
	}

	for _, task := range tasks {
		var estimate int
		err = db.QueryRow(`SELECT estimate FROM tasks WHERE name = ?`, task.name).Scan(&estimate)
		if err != nil {
			t.Fatalf("Failed to query database for task %s: %v", task.name, err)
		}

		if estimate != task.estimate {
			t.Errorf("Expected estimate for task %s to be %d, got %d", task.name, task.estimate, estimate)
		}
	}
}

