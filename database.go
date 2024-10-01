package main

import (
	"database/sql"
	"fmt"
	"log"
	"time"

	_ "github.com/mattn/go-sqlite3"
)

type DayAggregate struct {
    Day             string
    TotalEstimate   int
    TotalActual     int
    TotalDone       int
}

type TaskTrackingAggregate struct {
    Year        int
    Month       time.Month
    Day         int
    TaskCount   int
}


type TaskTracking struct {
    TaskID      int
    Date        string
    HalfHour    int
    Status      string // e.g., "in progress", "done", etc.
}

// Struct to hold task data
type Task struct {
    ID        int
    Name      string
    Estimate  int
    Actual    int
    CreatedAt time.Time
    UpdatedAt time.Time
    Done      bool
    Status    string
}

func initializeDatabase() *sql.DB {
    db, err := sql.Open("sqlite3", "./tomatillo.db")
    if err != nil {
        log.Fatal(err)
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
    _, err = db.Exec(createTableQuery)
    if err != nil {
        log.Fatal(err)
    }

    createTasksTableQuery := `
    CREATE TABLE IF NOT EXISTS task_tracking (
        id INTEGER PRIMARY KEY AUTOINCREMENT,
        task_id INTEGER NOT NULL,
        date DATE NOT NULL,
        half_hour INTEGER NOT NULL CHECK (half_hour BETWEEN 0 AND 47),
        task_name TEXT,
        status TEXT DEFAULT 'active',
        FOREIGN KEY (task_id) REFERENCES tasks(id) ON DELETE CASCADE
        UNIQUE(task_id, date, half_hour)
    );`

    _, err = db.Exec(createTasksTableQuery)
    if err != nil {
        log.Fatal(err)
    }

    return db
}

func getDailyTasks(db *sql.DB) ([]Task, error){
	query := `
    SELECT id, name, estimate, actual, created_at, updated_at, done 
    FROM tasks 
    WHERE DATE(datetime(created_at, 'localtime')) = DATE('now', 'localtime')
    ORDER BY created_at;
    `

	rows, err := db.Query(query)
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

    var tasks []Task
    for rows.Next() {
        var id, estimate, actual int
        var name string
        var createdAt, updatedAt time.Time
        var done bool

        err := rows.Scan(&id, &name, &estimate, &actual, &createdAt, &updatedAt, &done)
        if err != nil {
            log.Fatal(err)
        }

        tasks = append(tasks, Task{ ID: id, Name: name, Estimate: estimate, Actual: actual, CreatedAt: createdAt, UpdatedAt: updatedAt, Done: done })
    }
    return tasks, nil
}

func getMonthlyData(db *sql.DB) ([]DayAggregate, error) {
    query := `
    SELECT
        DATE(datetime(created_at, 'localtime')) as day,
        SUM(estimate) as total_estimate,
        SUM(actual) as total_actual,
        COUNT(CASE WHEN done = 1 THEN 1 ELSE NULL END) as total_done
    FROM tasks
    WHERE strftime('%Y-%m', datetime(created_at, 'localtime')) = strftime('%Y-%m', 'now', 'localtime')
    GROUP BY day
    ORDER BY day;
    `
    
    rows, err := db.Query(query)
    if err != nil {
        return nil, err
    }
    defer rows.Close()

    var dayAggregates []DayAggregate
    for rows.Next() {
        var dayAggregate DayAggregate
        err := rows.Scan(&dayAggregate.Day, &dayAggregate.TotalEstimate, &dayAggregate.TotalActual, &dayAggregate.TotalDone)
        if err != nil {
            return nil, err
        }
        dayAggregates = append(dayAggregates, dayAggregate)
    }

    return dayAggregates, nil
}

func getAllDaysOfMonth() []string {
    now := time.Now()
    year, month, _ := now.Date()
    location := now.Location()

    daysInMonth := time.Date(year, month+1, 0, 0, 0, 0, 0, location).Day()
    var days []string

    for i := 1; i <= daysInMonth; i++ {
        day := time.Date(year, month, i, 0, 0, 0, 0, location)
        days = append(days, day.Format("2006-01-02"))
    }

    return days
}

func getYearlyData(db *sql.DB, year int) ([]TaskTrackingAggregate, error) {
    // aggregate half_hours completed for each day for the year
    query := fmt.Sprintf(`
    WITH RECURSIVE all_dates AS (
        SELECT '%d-01-01' as date
        UNION ALL
        SELECT date(date, '+1 day')
        FROM all_dates
        WHERE date < '%d-12-31'
    )
    SELECT a.date, COUNT(t.date) as task_count
    FROM all_dates a
    LEFT JOIN task_tracking t ON a.date = t.date
    GROUP BY a.date
    ORDER BY a.date;`, year, year)

    rows, err := db.Query(query)
    if err != nil {
        return nil, err
    }
    defer rows.Close()

    var taskTrackingAggregates []TaskTrackingAggregate

    for rows.Next() {
        var dateStr string
        var count int

        if err := rows.Scan(&dateStr, &count); err != nil {
            return nil, err
        }

        // parse the date from the string format
        date, err := time.Parse("2006-01-02", dateStr)
        if err != nil {
            return nil, err
        }

        taskTrackingAggregates = append(taskTrackingAggregates, TaskTrackingAggregate{
            Year: date.Year(),
            Month: date.Month(),
            Day: date.Day(),
            TaskCount: count,
        })
    }
    return taskTrackingAggregates, nil
}

// Function to fetch tasks from the database
func getTasks(days int, status string) ([]Task, error) {
    var query string
    var rows *sql.Rows
    var err error

    // Build query based on status
    if status == "all" {
        query = fmt.Sprintf("SELECT * FROM tasks WHERE created_at >= date('now', '-%d days')", days)
    } else if status == "wip" || status == "inprogress" {
        query = fmt.Sprintf("SELECT * FROM tasks WHERE created_at >= date('now', '-%d days') AND done = 0 AND actual > 0", days)
    } else if status == "todo" {
        query = fmt.Sprintf("SELECT * FROM tasks WHERE created_at >= date('now', '-%d days') AND done = 0 AND actual = 0", days)
    } else if status == "done" {
        query = fmt.Sprintf("SELECT * FROM tasks WHERE created_at >= date('now', '-%d days') AND done = 1", days)
    } else {
        return nil, fmt.Errorf("invalid status filter")
    }

    rows, err = db.Query(query)
    if err != nil {
        return nil, err
    }
    defer rows.Close()

    // Slice to store task data
    tasks := []Task{}

    // Fetch data from rows
    for rows.Next() {
        var task Task
        err := rows.Scan(&task.ID, &task.Name, &task.Estimate, &task.Actual, &task.CreatedAt, &task.UpdatedAt, &task.Done)
        if err != nil {
            return nil, err
        }

        // Determine task status
        if task.Actual > 0 {
            task.Status = "In Progress"
        } else {
            task.Status = "To Do"
        }
        if task.Done {
            task.Status = "Done"
        }

        tasks = append(tasks, task)
    }

    return tasks, nil
}

func getTasksForDay(date string) ([]TaskTracking, error) {
    rows, err := db.Query(`SELECT task_id, date, half_hour, status FROM task_tracking WHERE date = ?`, date)
    if err != nil {
        return nil, err
    }
    defer rows.Close()

    var tasks []TaskTracking
    for rows.Next() {
        var task TaskTracking
        err := rows.Scan(&task.TaskID, &task.Date, &task.HalfHour, &task.Status)
        if err != nil {
            return nil, err
        }
        tasks = append(tasks, task)
    }

    return tasks, nil
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

func activateTask(id int) {
    // get current date in yyyy-mm-dd format
    // insert into task_tracking table
    currentDate := time.Now().Format("2006-01-02")
    half_hour := getHalfHour(time.Now().Hour(), time.Now().Minute())
    insertTrackingTask(id, currentDate, half_hour)
}

func insertTrackingTask(id int, currentDate string, halfHour int) {
    query := `
    INSERT INTO task_tracking (task_id, date, half_hour, status) 
    VALUES (?, ?, ?, 'active')
    ON CONFLICT(task_id, date, half_hour)
    DO UPDATE SET status = 'done';
    `
    _, err := db.Exec(query, id, currentDate, halfHour)
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