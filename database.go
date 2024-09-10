package main

import (
	"database/sql"
	"log"
	"time"

	_ "github.com/mattn/go-sqlite3"
)

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

    return db
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

type DayAggregate struct {
    Day          string
    TotalEstimate int
    TotalActual   int
    TotalDone     int
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

func getAllDaysOfWeek() []string {
    var days []string
    now := time.Now()

    for i := 0; i < 7; i++ {
        day := now.AddDate(0, 0, -int(now.Weekday())+i)
        days = append(days, day.Format("2006-01-02"))
    }

    return days
}

func getYearlyData(db *sql.DB) ([]DayAggregate, error) {
    query := `
    SELECT
        DATE(datetime(created_at, 'localtime')) as day,
        SUM(actual) as total_actual,
        SUM(estimate) as total_estimate,
        SUM(done) as total_done
    FROM tasks
    WHERE strftime('%Y', datetime(created_at, 'localtime')) = strftime('%Y', 'now', 'localtime')
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
        err := rows.Scan(&dayAggregate.Day, &dayAggregate.TotalActual, &dayAggregate.TotalEstimate, &dayAggregate.TotalDone)
        if err != nil {
            return nil, err
        }
        dayAggregates = append(dayAggregates, dayAggregate)
    }

    return dayAggregates, nil
}
