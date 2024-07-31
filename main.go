package main

import (
    "log"
    
    "tomatillo/task" // Corrected import path
    
    "github.com/rivo/tview"
)

func main() {
    taskManager := task.NewTaskManager()

    app := tview.NewApplication()

    // Task List View
    taskList := tview.NewList()

    // Form for adding tasks
    form := tview.NewForm()

    form.AddInputField("Title", "", 20, nil, nil).
        AddButton("Add Task", func() {
            title := form.GetFormItemByLabel("Title").(*tview.InputField).GetText()
            if title != "" {
                taskManager.AddTask(title)
                updateTaskList(taskList, taskManager)
            }
        }).
		AddButton("Quit", func () {
			app.Stop()
		})

    // Layout
    flex := tview.NewFlex().
        AddItem(taskList, 0, 1, true).
        AddItem(form, 0, 1, false)

    if err := app.SetRoot(flex, true).SetFocus(form).Run(); err != nil {
        log.Fatalf("Error launching application: %v\n", err)
    }
}

func updateTaskList(taskList *tview.List, taskManager *task.TaskManager) {
    taskList.Clear()
    for _, task := range taskManager.GetTasks() { // Corrected method name
        taskList.AddItem(task.Title, "", 0, nil)
    }
}
