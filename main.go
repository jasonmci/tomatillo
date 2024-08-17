package main

import (
	"database/sql"
	"flag"
	"fmt"
	"os"
)

var db *sql.DB

func main() {
    db = initializeDatabase()
    defer db.Close()

    if len(os.Args) < 2 {
        fmt.Println("expected 'add', 'list', 'update', 'done', 'edit', 'delete', or 'report' subcommands")
        os.Exit(1)
    }

    // Define subcommands
    addTaskFlag := flag.NewFlagSet("add", flag.ExitOnError)
    taskName := addTaskFlag.String("name", "", "Task name")
    taskEstimate := addTaskFlag.Int("estimate", 1, "Pomodoro estimate")

    listTasksFlag := flag.NewFlagSet("list", flag.ExitOnError)
    listDays := listTasksFlag.Int("days", 0, "Number of days' tasks to show")

    updateTaskFlag := flag.NewFlagSet("update", flag.ExitOnError)
    taskId := updateTaskFlag.Int("id", 0, "Task ID to update")

    doneTaskFlag := flag.NewFlagSet("done", flag.ExitOnError)
    doneTaskId := doneTaskFlag.Int("id", 0, "Task ID to mark as done")

    editTaskFlag := flag.NewFlagSet("edit", flag.ExitOnError)
    editTaskId := editTaskFlag.Int("id", 0, "Task ID to edit")
    newEstimate := editTaskFlag.Int("estimate", 1, "New Pomodoro estimate")

    reportFlag := flag.NewFlagSet("report", flag.ExitOnError)
    reportType := reportFlag.String("type", "monthly", "Report type: 'monthly' or 'yearly'")

    // delete flag
    deleteTaskFlag := flag.NewFlagSet("delete", flag.ExitOnError)
    deleteTaskId := deleteTaskFlag.Int("id", 0, "Task ID to delete")

    // Parse based on subcommand
    switch os.Args[1] {
    case "add":
        addTaskFlag.Parse(os.Args[2:])
        addTask(db, *taskName, *taskEstimate)
    case "list":
        listTasksFlag.Parse(os.Args[2:])
        listTasks(*listDays)
    case "update":
        updateTaskFlag.Parse(os.Args[2:])
        if *taskId > 0 {
            updateActual(*taskId)
        } else {
            fmt.Println("Please provide a valid task ID.")
        }
    case "done":
        doneTaskFlag.Parse(os.Args[2:])
        if *doneTaskId > 0 {
            markAsDone(*doneTaskId)
        } else {
            fmt.Println("Please provide a valid task ID.")
        }
    case "edit":
        editTaskFlag.Parse(os.Args[2:])
        if *editTaskId > 0 {
            updateEstimate(*editTaskId, *newEstimate)
        } else {
            fmt.Println("Please provide a valid task ID and new estimate.")
        }
    case "report":
        reportFlag.Parse(os.Args[2:])
        if *reportType == "yearly" {
            generateYearlyCalendarReport(db)
        } else {
            generateMonthlyReport(db)
        }
    // case delete
    case "delete":
        deleteTaskFlag.Parse(os.Args[2:])
        if *deleteTaskId > 0 {
            deleteTask(db, *deleteTaskId)
        } else {
            fmt.Println("Please provide a valid task ID.")
        }
    default:
        fmt.Println("expected 'add', 'list', 'update', 'done', 'edit', or 'report' subcommands")
        os.Exit(1)
    }
}
