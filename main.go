package main

import (
	"database/sql"
	"flag"
	"fmt"
	"log"
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

	// Subcommand handling
	command := os.Args[1]
	switch command {
	case "add":
		handleAddCommand(db, os.Args[2:])
	case "list":
		handleListCommand(os.Args[2:])
	case "update":
		handleUpdateCommand(os.Args[2:])
	case "done":
		handleDoneCommand(os.Args[2:])
	case "edit":
		handleEditCommand(os.Args[2:])
	case "report":
		handleReportCommand(os.Args[2:])
	case "delete":
		handleDeleteCommand(os.Args[2:])
	default:
		fmt.Println("expected 'add', 'list', 'update', 'done', 'edit', 'delete', or 'report' subcommands")
		os.Exit(1)
	}
}

// Helper function to handle the 'add' command
func handleAddCommand(db *sql.DB, args []string) error {
	addTaskFlag := flag.NewFlagSet("add", flag.ExitOnError)
	taskName := addTaskFlag.String("name", "", "Task name")
	taskEstimate := addTaskFlag.Int("estimate", 1, "Pomodoro estimate")
	addTaskFlag.Parse(args)

	if *taskName == "" {
		return fmt.Errorf("task name is required")
	}
	addTask(db, *taskName, *taskEstimate)
	return nil
}


// Helper function to handle the 'list' command
func handleListCommand(args []string) {
	listTasksFlag := flag.NewFlagSet("list", flag.ExitOnError)
	listDays := listTasksFlag.Int("days", 0, "Number of days' tasks to show")
	listTasksFlag.Parse(args)
	listTasks(*listDays)
}

// Helper function to handle the 'update' command
func handleUpdateCommand(args []string) {
	updateTaskFlag := flag.NewFlagSet("update", flag.ExitOnError)
	taskId := updateTaskFlag.Int("id", 0, "Task ID to update")
	updateTaskFlag.Parse(args)

	if *taskId <= 0 {
		log.Println("Please provide a valid task ID.")
		os.Exit(1)
	}
	updateActual(*taskId)
}

// Helper function to handle the 'done' command
func handleDoneCommand(args []string) {
	doneTaskFlag := flag.NewFlagSet("done", flag.ExitOnError)
	doneTaskId := doneTaskFlag.Int("id", 0, "Task ID to mark as done")
	doneTaskFlag.Parse(args)

	if *doneTaskId <= 0 {
		log.Println("Please provide a valid task ID.")
		os.Exit(1)
	}
	markAsDone(*doneTaskId)
}

// Helper function to handle the 'edit' command
func handleEditCommand(args []string) {
	editTaskFlag := flag.NewFlagSet("edit", flag.ExitOnError)
	editTaskId := editTaskFlag.Int("id", 0, "Task ID to edit")
	newEstimate := editTaskFlag.Int("estimate", 1, "New Pomodoro estimate")
	editTaskFlag.Parse(args)

	if *editTaskId <= 0 {
		log.Println("Please provide a valid task ID.")
		os.Exit(1)
	}
	updateEstimate(*editTaskId, *newEstimate)
}

// Helper function to handle the 'report' command
func handleReportCommand(args []string) {
	reportFlag := flag.NewFlagSet("report", flag.ExitOnError)
	reportType := reportFlag.String("type", "monthly", "Report type: 'monthly' or 'yearly'")
	reportFlag.Parse(args)

	if *reportType == "yearly" {
		generateYearlyCalendarReport(db)
	} else {
		generateMonthlyReport(db)
	}
}

// Helper function to handle the 'delete' command
func handleDeleteCommand(args []string) {
	deleteTaskFlag := flag.NewFlagSet("delete", flag.ExitOnError)
	deleteTaskId := deleteTaskFlag.Int("id", 0, "Task ID to delete")
	deleteTaskFlag.Parse(args)

	if *deleteTaskId <= 0 {
		log.Println("Please provide a valid task ID.")
		os.Exit(1)
	}
	deleteTask(db, *deleteTaskId)
}
