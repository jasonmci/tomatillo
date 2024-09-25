package main

import (
	"database/sql"
	"fmt"
	"log"
	"strings"
	"time"
)

func generateWeeklyReport(db *sql.DB) {
    data, err := getYearlyData(db)
    if err != nil {
        log.Fatal(err)
    }

    // print the weekly report from Sunday to Saturday
    fmt.Println("Weekly Report")
    fmt.Println("Day         | Done | Actual         ")
    fmt.Println("------------|------|----------------")

    allDays := getAllDaysOfWeek()
    dataMap := make(map[string]DayAggregate)

    // Map data to dates
    for _, dayAggregate := range data {
        dataMap[dayAggregate.Day] = dayAggregate
    }

    for _, day := range allDays {
        if aggregate, found := dataMap[day]; found {
            actualTomatoes := generateEmojis(aggregate.TotalActual, "🍅")
            fmt.Printf("%-11s | %-4d | %-7s\n", day, aggregate.TotalDone, actualTomatoes)
        } else {
            fmt.Printf("%-11s | %-4d | %-7s\n", day, 0, "")
        }
    }

}

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

// Function to wrap text in color
func colorize(text, color string) string {
    return fmt.Sprintf("\033[%sm%s\033[0m", color, text)
}


func checkTaskForHalfHour(halfHour int, taskMap map[int]string) string {
    // Check if a task exists in the specific half-hour slot
    if status, exists := taskMap[halfHour]; exists {
        return status // Return emoji if a task is found
    }
    return "·" // No task found, this is a middle dot, not a period.
}

func generateDailyReport(date string) {
    tasks, err := getTasksForDay(date)
    if err != nil {
        fmt.Println("Error fetching tasks:", err)
        return
    }
    
    // Initialize a map to track task status for each half-hour
    taskMap := make(map[int]string)
    for _, task := range tasks {
        taskMap[task.HalfHour] = colorize("▓", "32") // Use a tomato emoji for completed Pomodoros
    }


    // Print the header for hours
   
    fmt.Printf("║ %s ", date)

    // Loop through the 48 half-hour slots (0 to 47)
    for i := 0; i < 48; i += 2 {
        // Check if tasks exist in each half-hour slot
        firstHalf := checkTaskForHalfHour(i, taskMap)
        secondHalf := checkTaskForHalfHour(i + 1, taskMap)

        // Print the status for the two half-hour slots in each hour
        fmt.Printf("%s%s ", firstHalf, secondHalf)
        
    }
    fmt.Println("║")
}


// Generate a weekly block report
func generateWeeklyBlockReport() {
    startOfWeek, endOfWeek := getCurrentWeek()
    fmt.Println("╔══════════════════════════════════════════╗ ")
    fmt.Printf( "║ Weekly Report (%s to %s) ║ \n", startOfWeek.Format("2006-01-02"), endOfWeek.Format("2006-01-02"))
    fmt.Println("╠══════════════════════════════════════════╩═════════════════════════════════════════╗ ")
    fmt.Println("║            00|01|02|03|04|05|06|07|08|09|10|11|12|13|14|15|16|17|18|19|20|21|22|23 ║ ")
    fmt.Println("╠════════════════════════════════════════════════════════════════════════════════════╣ ")

    // Iterate through each day of the week
    for i := 0; i < 7; i++ {
        day := time.Now().Local().AddDate(0, 0, -int(time.Now().Weekday())+i).Format("2006-01-02")
        //fmt.Printf("\n%s\n", day)
        generateDailyReport(day)  // Reuse your daily report generation
    }

    fmt.Println("╚════════════════════════════════════════════════════════════════════════════════════╝ ")
}

// Generate a monthly block report
func generateMonthlyBlockReport() {
    startOfMonth, endOfMonth := getCurrentMonth()
    fmt.Println("╔═══════════════════════════════════════════╗ ")
    fmt.Printf( "║ Monthly Report (%s to %s) ║ \n", startOfMonth.Format("2006-01-02"), endOfMonth.Format("2006-01-02"))
    fmt.Println("╠═══════════════════════════════════════════╩════════════════════════════════════════╗ ")
    fmt.Println("║            00|01|02|03|04|05|06|07|08|09|10|11|12|13|14|15|16|17|18|19|20|21|22|23 ║ ")
    fmt.Println("╠════════════════════════════════════════════════════════════════════════════════════╣ ")

    // Iterate through each day of the month
    for day := startOfMonth; !day.After(endOfMonth); day = day.AddDate(0, 0, 1) {
        dayStr := day.Format("2006-01-02")
        //fmt.Printf("\n%s\n", day)
        generateDailyReport(dayStr)  // Reuse your daily report generation
    }

    fmt.Println("╚════════════════════════════════════════════════════════════════════════════════════╝ ")
}