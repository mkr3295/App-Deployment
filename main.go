package main

import (
	"encoding/json"
	"fmt"
	"os"
	"strconv"
)

type Task struct {
	ID        int    `json:"id"`
	Title     string `json:"title"`
	Completed bool   `json:"completed"`
}

const dataFile = "tasks.json"

func loadTasks() ([]Task, error) {
	var tasks []Task

	file, err := os.ReadFile(dataFile)
	if err != nil {
		if os.IsNotExist(err) {
			return tasks, nil
		}
		return nil, err
	}

	err = json.Unmarshal(file, &tasks)
	return tasks, err
}

func saveTasks(tasks []Task) error {
	data, err := json.MarshalIndent(tasks, "", "  ")
	if err != nil {
		return err
	}

	return os.WriteFile(dataFile, data, 0644)
}

func addTask(title string) error {
	tasks, err := loadTasks()
	if err != nil {
		return err
	}

	id := 1
	if len(tasks) > 0 {
		id = tasks[len(tasks)-1].ID + 1
	}

	task := Task{
		ID:        id,
		Title:     title,
		Completed: false,
	}

	tasks = append(tasks, task)
	return saveTasks(tasks)
}

func listTasks() error {
	tasks, err := loadTasks()
	if err != nil {
		return err
	}

	if len(tasks) == 0 {
		fmt.Println("No tasks found.")
		return nil
	}

	for _, task := range tasks {
		status := " "
		if task.Completed {
			status = "✓"
		}

		fmt.Printf("[%s] %d - %s\n", status, task.ID, task.Title)
	}

	return nil
}

func completeTask(id int) error {
	tasks, err := loadTasks()
	if err != nil {
		return err
	}

	found := false

	for i := range tasks {
		if tasks[i].ID == id {
			tasks[i].Completed = true
			found = true
			break
		}
	}

	if !found {
		return fmt.Errorf("task %d not found", id)
	}

	return saveTasks(tasks)
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage:")
		fmt.Println("  go run main.go add \"Task name\"")
		fmt.Println("  go run main.go list")
		fmt.Println("  go run main.go done <id>")
		return
	}

	command := os.Args[1]

	switch command {
	case "add":
		if len(os.Args) < 3 {
			fmt.Println("Please provide a task title.")
			return
		}

		if err := addTask(os.Args[2]); err != nil {
			fmt.Println("Error:", err)
			return
		}

		fmt.Println("Task added!")

	case "list":
		if err := listTasks(); err != nil {
			fmt.Println("Error:", err)
		}

	case "done":
		if len(os.Args) < 3 {
			fmt.Println("Please provide a task ID.")
			return
		}

		id, err := strconv.Atoi(os.Args[2])
		if err != nil {
			fmt.Println("Invalid task ID.")
			return
		}

		if err := completeTask(id); err != nil {
			fmt.Println("Error:", err)
			return
		}

		fmt.Println("Task completed!")

	default:
		fmt.Println("Unknown command.")
	}
}
