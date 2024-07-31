package task

import (
	"testing"
)

func TestAddTask(t *testing.T) {
	tm := NewTaskManager()

	tm.AddTask("Test Task 1")

	if len(tm.tasks) != 1 {
		t.Fatalf("Expected 1 task, got %d", len(tm.GetTasks()))
	}
}

func TestGetTasks(t *testing.T) {
	tm := NewTaskManager()
	tm.AddTask("Test Task2")
	tm.AddTask("Test Task3")

	tasks := tm.GetTasks()
	if len(tasks) != 2 {
		t.Fatalf("Expected 2 tasks, got %d", len(tasks))
	}

	if tasks[0].Title != "Test Task2" || tasks[1].Title != "Test Task3" {
		t.Errorf("Task titles do not match")

	}
}

func TestGetTaskIDs(t *testing.T) {
	tm := NewTaskManager()
	tm.AddTask("Test Task4")
	tm.AddTask("Test Task5")

	tasks := tm.GetTasks()
	
	if tasks[0].ID != 1 || tasks[1].ID != 2 {
		t.Errorf("Task IDs do not match")
	}

}