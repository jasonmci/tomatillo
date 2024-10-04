package main

import (
	"time"
)

func getWeek(mytime time.Time) (time.Time, time.Time) {
    sunday := mytime.AddDate(0, 0, -int(mytime.Weekday()))
    saturday := sunday.AddDate(0, 0, 6)
    return sunday, saturday
}

func getMonth(mytime time.Time) (time.Time, time.Time) {

    // Find the first day of the month
    firstOfMonth := time.Date(mytime.Year(), mytime.Month(), 1, 0, 0, 0, 0, mytime.Location())

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

func generateEmojis(count int, emoji string) string {
    result := ""
    for i := 0; i < count; i++ {
        result += emoji
    }
    return result
}
