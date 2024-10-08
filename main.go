package main

import (
	"bufio"
	"database/sql"
	"flag"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
	"time"
)

var db *sql.DB

func main() {
	db = initializeDatabase("./tomatillo.db")
	defer db.Close()

	if len(os.Args) < 2 {
		fmt.Println("expected 'add', 'activate', 'simple', 'list', 'update', 'done', 'edit', 'delete', 'load', or 'report' subcommands")
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
	case "activate":
		handleActivateCommand(os.Args[2:])
	case "backfill":
		handleBackfillCommand(os.Args[2:])
	case "load":
		handleLoadTasksCommand(db, os.Args[2:])
	case "today":
		// use the handle report command with the --type flag set to today
		os.Args = append(os.Args, "--type", "today")
		handleReportCommand(os.Args[2:])
	case "version":
		fmt.Println("tomatillo v0.1")
	case "help":
		handleHelpCommand()
	default:
		fmt.Println("expected 'add', 'activate', 'today', 'list', 'update', 'done', 'edit', 'delete', 'load', version', or 'report' subcommands")
		os.Exit(1)
	}
}

func handleHelpCommand() {
	fmt.Println("Usage: task [command] [arguments]")
	fmt.Println("\nCommands:")
	fmt.Println("  add     Add a new task")
	fmt.Println("  list    List tasks")
	fmt.Println("  update  Update the actual pomodoros of a task")
	fmt.Println("  done    Mark a task as done")
	fmt.Println("  edit    Edit the estimate of a task")
	fmt.Println("  report  Generate a report")
	fmt.Println("  delete  Delete a task")
	fmt.Println("  load    Load tasks from a file")
	fmt.Println("  version Print the version of the application")
}

// Helper function to handle the 'add' command
func handleAddCommand(db *sql.DB, args []string) error {
	addTaskFlag := flag.NewFlagSet("add", flag.ExitOnError)
	taskName := addTaskFlag.String("name", "", "Task name (or use -n)")
	taskEstimate := addTaskFlag.Int("estimate", 1, "Pomodoro estimate (or use -e)")
	addTaskFlag.StringVar(taskName, "n", "", "Task name (short version)")
	addTaskFlag.IntVar(taskEstimate, "e", 1, "Pomodoro estimate (short version)")

	addTaskFlag.Parse(args)

	if *taskName == "" {
		return fmt.Errorf("task name is required")
	}
	addTask(db, *taskName, *taskEstimate)
	return nil
}

func handleLoadTasksCommand(db *sql.DB, args []string) error {
	loadTasksFlag := flag.NewFlagSet("load", flag.ExitOnError)
	filePath := loadTasksFlag.String("file", "", "Path to the file containing tasks and estimates")
	// add a short version of the flag
	loadTasksFlag.StringVar(filePath, "f", "", "Path to the file containing tasks and estimates")

	loadTasksFlag.Parse(args)

	if *filePath == "" {
		return fmt.Errorf("file path is required")
	}

	file, err := os.Open(*filePath)
	if err != nil {
		return fmt.Errorf("failed to open file: %v", err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		parts := strings.Split(line, ",")
		if len(parts) != 2 {
			return fmt.Errorf("invalid format in line: %s", line)
		}
		taskName := strings.TrimSpace(parts[0])
		estimate, err := strconv.Atoi(strings.TrimSpace(parts[1]))
		if err != nil {
			return fmt.Errorf("invalid estimate in line: %s", line)
		}
		err = addTask(db, taskName, estimate)
		if err != nil {
			return fmt.Errorf("failed to add task: %v", err)
		}
	}

	if err := scanner.Err(); err != nil {
		return fmt.Errorf("error reading file: %v", err)
	}

	fmt.Println("Tasks successfully loaded from file.")
	return nil
}


// Helper function to handle the 'list' command
func handleListCommand(args []string) {
	listTasksFlag := flag.NewFlagSet("list", flag.ExitOnError)
	listDays := listTasksFlag.Int("days", 0, "Number of days' tasks to show")

	// add a short version of the flag
	listTasksFlag.IntVar(listDays, "d", 0, "Number of days' tasks to show")

	status := listTasksFlag.String("status", "all", "Status of tasks to show: 'all', 'done', 'todo', 'wip'")
	listTasksFlag.StringVar(status, "s", "all", "Short version of status filter: active, completed, or all")

	listTasksFlag.Parse(args)
    listTasks(*listDays, strings.ToLower(*status))
}


// helper function to handle activating a current task
func handleActivateCommand(args []string) {
	activateTaskFlag := flag.NewFlagSet("activate", flag.ExitOnError)
	activateTaskId := activateTaskFlag.Int("id", 0, "Task ID to activate")
	activateTaskFlag.Parse(args)

	if *activateTaskId <= 0 {
		log.Println("Please provide a valid task ID.")
		os.Exit(1)
	}

	// insert into task_tracking table
	currentDate := time.Now().Format("2006-01-02")
	half_hour := getHalfHour(time.Now().Hour(), time.Now().Minute())
	insertTrackingTask(*activateTaskId, currentDate, half_hour)
}

func handleBackfillCommand(args []string) {
	backfillFlag := flag.NewFlagSet("backfill", flag.ExitOnError)
	backfillTaskId := backfillFlag.Int("id", 0, "Task ID to backfill")
	backfillTaskDate := backfillFlag.String("date", "", "Date to backfill the task")
	backfillTaskHalfHour := backfillFlag.Int("halfhour", 0, "Half hour to backfill the task")
	backfillFlag.Parse(args)

	if *backfillTaskId <= 0 {
		log.Println("Please provide a valid task ID.")
		os.Exit(1)
	}
	insertTrackingTask(*backfillTaskId, *backfillTaskDate, *backfillTaskHalfHour)
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
	// add a short version of the flag
	editTaskFlag.IntVar(newEstimate, "e", 1, "New Pomodoro estimate")
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
	reportType := reportFlag.String("type", "monthly", "Report type: 'monthly','yearly' or 'weekly'")
	// add a short version of the flag
	reportFlag.StringVar(reportType, "t", "weekly", "Report type: 'monthly' or 'yearly' or 'weekly'")
	reportFlag.Parse(args)

	if *reportType == "yearly" {
		generateYearlyCountReport()
	} else if *reportType == "blockmonth" {
			generateMonthlyBlockReport()
	} else if *reportType == "today" {
		generateTodayReport()
	}  else if *reportType == "blockweek" {
		generateWeeklyBlockReport()
	} else {
		generateTodayReport()
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
