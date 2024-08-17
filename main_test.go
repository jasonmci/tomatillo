package main

import (
	"testing"
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

