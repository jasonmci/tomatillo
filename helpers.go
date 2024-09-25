package main

import (
	"database/sql"
	"fmt"
	"log"
	"strings"
	"time"
)

// Get start and end dates for the current week (Monday to Sunday)
func getCurrentWeek() (string, string) {
    // Get the current date
    now := time.Now().Local()
    //fmt.Println(now)
    // Find the Monday of the current week
    sunday := now.AddDate(0, 0, -int(now.Weekday()))
    // Get the Sunday of the current week
    saturday := sunday.AddDate(0, 0, 6)
    //fmt.Println("---", monday)

    // Format as YYYY-MM-DD
    return sunday.Format("2006-01-02"), saturday.Format("2006-01-02")
}

func getCurrentMonth() (time.Time, time.Time) {
    // Get the current date
    now := time.Now().Local()

    // Find the first day of the month
    firstOfMonth := time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, now.Location())

    // Find the last day of the month by going to the next month and subtracting one day
    nextMonth := firstOfMonth.AddDate(0, 1, 0)
    lastOfMonth := nextMonth.AddDate(0, 0, -1)

    return firstOfMonth, lastOfMonth
}


func generateEmojis(count int, emoji string) string {
    result := ""
    for i := 0; i < count; i++ {
        result += emoji
    }
    return result
}

func addTask(db *sql.DB, name string, estimate int) error {
    if name == "" {
        return fmt.Errorf("task name cannot be empty")
    }

    now := time.Now().Local()

    query := `INSERT INTO tasks (name, estimate, actual, created_at, updated_at, done) 
    VALUES (?, ?, 0, ?, ?, 0)`
    
    result, err := db.Exec(query, name, estimate, now.Format("2006-01-02 15:04:05"), now.Format("2006-01-02 15:04:05"))
    if err != nil {
        return fmt.Errorf("failed to add task: %v", err)
    }
    
    id, err := result.LastInsertId()
    if err != nil {
        return fmt.Errorf("failed to get the ID of the inserted task: %v", err)
    }
    
    estimateSprouts := generateEmojis(estimate, "🌱")
    fmt.Printf("Added task: %s\nID: %d\nEstimate: %d %s\n", name, id, estimate, estimateSprouts)
    return nil
}


func getHalfHour(hour int, minute int) int {

    if minute >= 30 {
        return hour * 2 + 1
    } else {
        return hour * 2
    }
}

func activateTask(id int) {
    // get current date in yyyy-mm-dd format
    // insert into task_tracking table

    currentDate := time.Now().Format("2006-01-02")
    half_hour := getHalfHour(time.Now().Hour(), time.Now().Minute())

    query := `
    INSERT INTO task_tracking (task_id, date, half_hour, status) 
    VALUES (?, ?, ?, 'active')
    ON CONFLICT(task_id, date, half_hour)
    DO UPDATE SET status = 'done';
    `
    _, err := db.Exec(query, id, currentDate, half_hour)
    if err != nil {
        log.Fatal(err)
    }
}

func updateActual(id int) {

    query := `UPDATE tasks SET actual = actual + 1, updated_at = datetime('now', 'localtime') WHERE id = ?`
    result, err := db.Exec(query, id)
    if err != nil {
        log.Fatal(err)
    }

      
    rowsAffected, err := result.RowsAffected()
    if err != nil {
        log.Fatal(err)
    }

    if rowsAffected == 0 {
        fmt.Printf("No task found with ID: %d\n", id)
    } else {
        fmt.Printf("Updated task with ID: %d, increased 'actual' count by 1\n", id)
    }
}

func updateEstimate(id int, newEstimate int) {
    query := `UPDATE tasks SET estimate = ?, updated_at = datetime('now', 'localtime') WHERE id = ?`
    result, err := db.Exec(query, newEstimate, id)
    if err != nil {
        log.Fatal(err)
    }

    rowsAffected, err := result.RowsAffected()
    if err != nil {
        log.Fatal(err)
    }

    if rowsAffected == 0 {
        fmt.Printf("No task found with ID: %d\n", id)
    } else {
        fmt.Printf("Task with ID: %d has been updated with new estimate: %d 🌱\n", id, newEstimate)
    }
}

func markAsDone(id int) {
    query := `UPDATE tasks SET done = 1, updated_at =  datetime('now', 'localtime') WHERE id = ?`
    result, err := db.Exec(query, id)
    if err != nil {
        log.Fatal(err)
    }

    rowsAffected, err := result.RowsAffected()
    if err != nil {
        log.Fatal(err)
    }

    if rowsAffected == 0 {
        fmt.Printf("No task found with ID: %d\n", id)
    } else {
        fmt.Printf("Task with ID: %d has been marked as done\n", id)
    }
}

func listTasks(days int) {
    var query string
    var rows *sql.Rows
    var err error

    if days == 1 {
        query := `
        SELECT id, name, estimate, actual, DATE(created_at), DATE(updated_at), done 
        FROM tasks
        WHERE DATE(datetime(created_at, 'localtime')) = DATE('now', 'localtime')
        ORDER BY created_at;
        `
        rows, err = db.Query(query)
    } else if days > 1 {
        query = `
        SELECT id, name, estimate, actual, DATE(created_at), DATE(updated_at), done 
        FROM tasks 
        WHERE DATE(datetime(created_at, 'localtime')) >= DATE('now', 'localtime', ? || ' days')
        ORDER BY created_at DESC;
        `
        rows, err = db.Query(query, fmt.Sprintf("-%d", days))
    } else {
        query = "SELECT id, name, estimate, actual, DATE(created_at), DATE(updated_at), done FROM tasks ORDER BY created_at DESC;"
        rows, err = db.Query(query)
    }

    if err != nil {
        log.Fatal(err)
    }
    defer rows.Close()

    fmt.Printf("%-3s | %-50s | %-12s | %-12s\n", "ID", "Name", "Created", "Updated")
    fmt.Println(strings.Repeat("-", 84))

    for rows.Next() {
        var id int
        var name string
        var estimate int
        var actual int
        var createdAt string
        var updatedAt string
        var done bool
        err := rows.Scan(&id, &name, &estimate, &actual, &createdAt, &updatedAt, &done)
        if err != nil {
            log.Fatal(err)
        }
        estimateSprouts := generateEmojis(estimate, "🌱")
        actualTomatoes := generateEmojis(actual, "🍅")
        doneStatus := "To Do"
        // if actual is greater than 0 then it's in progress
        if actual > 0 {
            doneStatus = "In Progress"
        }
        if done {
            doneStatus = "Done"
        }
        fmt.Printf("%-3d |  %-1s %-46s | %-12s | %-12s\n", id, "📋", name, createdAt, updatedAt)
        fmt.Printf("    |   Status:   %s\n", doneStatus)
        fmt.Printf("    |   Estimate: %s\n", estimateSprouts)
        fmt.Printf("    |   Actual:   %s\n", actualTomatoes)
        fmt.Println(strings.Repeat("-", 84))
    }
}

// delete a task
func deleteTask(db *sql.DB, id int) {
    query := `DELETE FROM tasks WHERE id = ?`
    result, err := db.Exec(query, id)
    if err != nil {
        log.Fatal(err)
    }

    rowsAffected, err := result.RowsAffected()
    if err != nil {
        log.Fatal(err)
    }

    if rowsAffected == 0 {
        fmt.Printf("No task found with ID: %d\n", id)
    } else {
        fmt.Printf("Task with ID: %d has been deleted\n", id)
    }
}