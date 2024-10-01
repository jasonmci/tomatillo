package main

import (
	"database/sql"
	"fmt"
	"log"
	"strings"
	"time"

)

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

func formatDate(t time.Time) string {
    return t.Format("2006-01-02") // Go uses a reference date to specify the format
}

func getHalfHour(hour int, minute int) int {

    if minute >= 30 {
        return hour * 2 + 1
    } else {
        return hour * 2
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
    fmt.Println("════════════|══════|════════════════")

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

func generateDailyBlock(date string) {
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
        generateDailyBlock(day)  // Reuse your daily report generation
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
        generateDailyBlock(dayStr)  // Reuse your daily report generation
    }

    fmt.Println("╚════════════════════════════════════════════════════════════════════════════════════╝ ")
}

// generate a report for yearly data of tasks completed. each row is a month and each column is a day
func generateYearlyCountReport() {
    
    t := time.Now().Local()
    year := t.Year()

    reports, _ := getYearlyData(db, year)
    var currentMonth time.Month
    var lastDay int
    startOfYear, endOfYear := getCurrentYear()
    fmt.Println("╔═══════════════════════════════════════════╗ ")
    fmt.Printf( "║ Yearly Report (%s to %s)  ║  \n", startOfYear.Format("2006-01-02"), endOfYear.Format("2006-01-02")) 
    fmt.Println("╠═══════════════════════════════════════════╩═══════════════════════════════════════════════════════╗ ")
    fmt.Print("║       01 02 03 04 05 06 07 08 09 10 11 12 13 14 15 16 17 18 19 20 21 22 23 24 25 26 27 28 29 30 31║ ")

    for _, report := range reports {
        // Check if the month changes, and print the header for a new month
        if report.Month != currentMonth {

            currentMonth = report.Month
            monthAbbrev := getMonthAbbreviation(currentMonth)
            if currentMonth == 3 && lastDay == 29 {
                fmt.Print("      ║")
            } else if currentMonth == 3 && lastDay == 28 {
                fmt.Print("         ║")    
            } else if currentMonth == 5 || currentMonth == 7 || currentMonth == 10 || currentMonth == 12 {
                fmt.Print("   ║")             
            } else if currentMonth.String() != "January" {
                fmt.Print("║")
            }
            fmt.Printf("\n║ %s  ", monthAbbrev)
        }
            
        // Print each day's task count
        if report.TaskCount == 0 {
            fmt.Printf("%3s", "··")
        } else {
            fmt.Printf(" %3s", colorize(fmt.Sprintf("%2d", report.TaskCount), "33"))
        }
        lastDay = report.Day
    }

    fmt.Println("║\n╚═══════════════════════════════════════════════════════════════════════════════════════════════════╝ ")
    fmt.Println()
}

// Function to generate report from the data
func generateTaskReport(tasks []Task) {
    fmt.Printf("%-3s   %-46s   %-12s   %-12s\n", "ID", "Name", "Created", "Updated")
    fmt.Println(strings.Repeat("═", 80))

    for _, task := range tasks {
        estimateSprouts := generateEmojis(task.Estimate, "🌱")
        actualTomatoes := generateEmojis(task.Actual, "🍅")

        fmt.Printf("%-3d   %-46s   %-12s   %-12s\n", task.ID, task.Name, formatDate(task.CreatedAt), formatDate(task.UpdatedAt))
        fmt.Printf("      %s\n", task.Status)
        fmt.Printf("      Estimate: %s Actual: %s\n", estimateSprouts, actualTomatoes)
        fmt.Println(strings.Repeat("═", 80))
    }
}

func generateTodayReport() {
    tasks, err := getDailyTasks(db)
    if err != nil {
        log.Fatal(err)
    }

    fmt.Println("╔════════════════════════════════════════════════════════════════════════════════════╗ ")
	fmt.Printf( "║ %-3s   %-5s   %-54s   %-4s   %-4s ║\n", "ID", "Done?", "Task", "Est.", "Act.")
	fmt.Println("╠════════════════════════════════════════════════════════════════════════════════════╣ ")

    for _, task := range tasks {
        id := task.ID
        name := task.Name
        estimate := task.Estimate
        actual := task.Actual
        done := task.Done
        completedTasks := 0

        status := "No"
        if done {
            status = "Yes"
            completedTasks++
        } 
        fmt.Printf("║ %-3d   %-5s   %-54s   %-4d   %-4d ║\n", id, status, name, estimate, actual)
    }
    fmt.Println("╚════════════════════════════════════════════════════════════════════════════════════╝ ")
}


func listTasks(days int, status string) {
    tasks, err := getTasks(days, status)
    if err != nil {
        log.Fatal(err)
    }

    generateTaskReport(tasks)
}