package main

import (
    "database/sql"
    "fmt"
    "log"
    "math/rand"
    "time"

    _ "github.com/mattn/go-sqlite3"
)

func main() {
    db, err := sql.Open("sqlite3", ":memory:")
    if err != nil {
        log.Fatal(err)
    }
    defer db.Close()

    // Define the number of tasks and the date range
    numberOfTasks := 100
    monthsBack := 3
    seed := time.Now().UnixNano()
    rand.Seed(seed)

    // Generate random tasks
    for i := 0; i < numberOfTasks; i++ {
        name := fmt.Sprintf("Task %d", i+1)
        estimate := rand.Intn(5) + 1        // Estimate between 1 and 5
        actual := rand.Intn(estimate + 1)   // Actual between 0 and estimate
        daysAgo := rand.Intn(30 * monthsBack)
        createdAt := time.Now().AddDate(0, 0, -daysAgo)
        updatedAt := createdAt.Add(time.Duration(rand.Intn(86400)) * time.Second) // Random time on the same day or later
        done := 0
        if rand.Float32() > 0.5 {
            done = 1
        }

        // Insert the task into the database
        _, err := db.Exec(`INSERT INTO tasks (name, estimate, actual, created_at, updated_at, done) 
                           VALUES (?, ?, ?, ?, ?, ?)`,
            name, estimate, actual, createdAt.Format("2006-01-02 15:04:05"),
            updatedAt.Format("2006-01-02 15:04:05"), done)
        if err != nil {
            log.Fatal(err)
        }
    }

    fmt.Printf("Generated %d tasks spanning the last %d months.\n", numberOfTasks, monthsBack)
}
