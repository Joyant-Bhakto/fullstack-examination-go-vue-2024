package model

import "time"

// Todo is the model for the todo endpoint.
type Todo struct {
	ID        int `gorm:"primaryKey"`
	Task      string
	Status    Status
	Priority  int
	CreatedAt time.Time `gorm:"autoCreateTime"`
	UpdatedAt time.Time `gorm:"autoUpdateTime"`
}

// NewTodo returns a new instance of the todo model.
func NewTodo(task string, priority int) *Todo {
	return &Todo{
		Task:     task,
		Status:   Created,
		Priority: priority,
	}
}

// NewUpdateTodo returns a new instance of the todo model for updating.
func NewUpdateTodo(id int, task string, status Status, priority int) *Todo {
	return &Todo{
		ID:       id,
		Task:     task,
		Status:   status,
		Priority: priority,
	}
}

// Status is the status of the task.
type Status string

const (
	// Created is the status for a created task.
	Created = Status("created")
	// Processing is the status for a processing task.
	Processing = Status("processing")
	// Done is the status for a done task.
	Done = Status("done")
)

// StatusMap is a map of task status.
var StatusMap = map[Status]bool{
	Created:    true,
	Processing: true,
	Done:       true,
}
