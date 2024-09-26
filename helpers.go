package main

import (
	"time"
)

// Get start and end dates for the current week (Monday to Sunday)
func getCurrentWeek() (time.Time, time.Time) {
    // Get the current date
    now := time.Now().Local()
    // Find the Monday of the current week
    sunday := now.AddDate(0, 0, -int(now.Weekday()))
    // Get the Sunday of the current week
    saturday := sunday.AddDate(0, 0, 6)

    // Format as YYYY-MM-DD
    return sunday, saturday
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

// Helper function to get the first three letters of the month
func getMonthAbbreviation(month time.Month) string {
    // Get the full month name and return the first three letters
    return month.String()[:3]
}

func getCurrentYear() (time.Time, time.Time) {
    // Get the current date
    now := time.Now().Local()

    // Find the first day of the year
    firstOfYear := time.Date(now.Year(), 1, 1, 0, 0, 0, 0, now.Location())

    // Find the last day of the year by going to the next year and subtracting one day
    nextYear := firstOfYear.AddDate(1, 0, 0)
    lastOfYear := nextYear.AddDate(0, 0, -1)

    return firstOfYear, lastOfYear
}

func generateEmojis(count int, emoji string) string {
    result := ""
    for i := 0; i < count; i++ {
        result += emoji
    }
    return result
}
