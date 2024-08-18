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
)

var db *sql.DB

func main() {
	db = initializeDatabase()
	defer db.Close()

	if len(os.Args) < 2 {
		fmt.Println("expected 'add', 'list', 'update', 'done', 'edit', 'delete', 'load', or 'report' subcommands")
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
	case "load":
		handleLoadTasksCommand(db, os.Args[2:])
	case "concise":
		handleConciseReport(db)
	default:
		fmt.Println("expected 'add', 'concise', 'list', 'update', 'done', 'edit', 'delete', 'load', or 'report' subcommands")
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

func handleConciseReport(db *sql.DB) {
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

	var totalEstimate, totalActual, totalTasks, completedTasks int

	fmt.Printf("| %-3s | %-5s | %-54s | %-4s | %-4s |\n", "ID", "Done?", "Task", "Est.", "Act.")
	fmt.Println("| --- | ----- | ------------------------------------------------------ | ---- | ---- |")
	//fmt.Println(strings.Repeat("-", 80))

	for rows.Next() {
		var id, estimate, actual int
		var name, createdAt, updatedAt string
		var done bool

		err := rows.Scan(&id, &name, &estimate, &actual, &createdAt, &updatedAt, &done)
		if err != nil {
			log.Fatal(err)
		}

		status := "No"
		if done {
			status = "Yes"
			completedTasks++
		}

		totalEstimate += estimate
		totalActual += actual
		totalTasks++

		fmt.Printf("| %-3d | %-5s | %-54s | %-4d | %-4d |\n", id, status, name, estimate, actual)
	}

	fmt.Println("\n**Summary:**")
	fmt.Printf("\n- Estimated: %d\n- Actual:    %d\n",
		totalEstimate, totalActual)
	fmt.Printf("- Tasks Completed/Total: %d of %d\n", completedTasks, totalTasks)
}


func handleLoadTasksCommand(db *sql.DB, args []string) error {
	loadTasksFlag := flag.NewFlagSet("load", flag.ExitOnError)
	filePath := loadTasksFlag.String("file", "", "Path to the file containing tasks and estimates")
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
