package main

import (
	"database/sql"
	"fmt"
	"log"
	"strings"
	"time"
)

func generateMonthlyReport(db *sql.DB) {
    data, err := getMonthlyData(db)
    if err != nil {
        log.Fatal(err)
    }

    allDays := getAllDaysOfMonth()
    dataMap := make(map[string]DayAggregate)

    // Map data to dates
    for _, dayAggregate := range data {
        dataMap[dayAggregate.Day] = dayAggregate
    }

    fmt.Println("Monthly Report")
    fmt.Println("Day         | Done | Actual         ")
    fmt.Println("------------|------|----------------")

    for _, day := range allDays {
        if aggregate, found := dataMap[day]; found {
            //estimateSprouts := generateEmojis(aggregate.TotalEstimate, "🌱")
            actualTomatoes := generateEmojis(aggregate.TotalActual, "🍅")
            fmt.Printf("%-11s | %-4d | %-7s\n", day, aggregate.TotalDone, actualTomatoes)
        } else {
            fmt.Printf("%-11s | %-4d | %-7s\n", day, 0, "")
        }
    }
}

func generateYearlyCalendarReport(db *sql.DB) {
    data, err := getYearlyData(db)
    if err != nil {
        log.Fatal(err)
    }

    dataMap := make(map[string]int)
    for _, dayAggregate := range data {
        dataMap[dayAggregate.Day] = dayAggregate.TotalActual
    }

    now := time.Now()
    year := now.Year()
    currentMonth := int(now.Month())

    fmt.Printf("Tomatillo Completion Calendar (Last 12 Months) ending %s %d\n", now.Format("January"), year)

    fmt.Println(strings.Repeat("-", 90))

    for offset := 11; offset >= 0; offset-- {

        month := currentMonth - offset
        yearOffset := 0 

        if month <= 0 {
            month += 12
            yearOffset = -1
        }

        firstDay := time.Date(year + yearOffset, time.Month(month), 1, 0, 0, 0, 0, time.Local)
        daysInMonth := firstDay.AddDate(0, 1, -1).Day()

        fmt.Printf("\n\n%s\n", firstDay.Format("January"))
        fmt.Println("Su Mo Tu We Th Fr Sa")

        startOffset := int(firstDay.Weekday())
        if startOffset > 0 {
            fmt.Print(strings.Repeat("   ", startOffset))
        }

        for day := 1; day <= daysInMonth; day++ {
            currentDay := firstDay.AddDate(0, 0, day-1).Format("2006-01-02")
            tomatoes := dataMap[currentDay]
            if tomatoes > 9 {
                fmt.Printf("\033[32m%2d \033[0m", tomatoes) // Green text for counts > 9
            } else if tomatoes < 3 && tomatoes > 0{
                // print in yellow text
                fmt.Printf("\033[33m%2d \033[0m", tomatoes)

            } else if tomatoes >= 3 && tomatoes <= 9 {
                fmt.Printf("%2d ", tomatoes)
            } else {
                fmt.Printf("-- ")
            }

            if (day+startOffset)%7 == 0 {
                fmt.Println()
            }
        }
    }
}
