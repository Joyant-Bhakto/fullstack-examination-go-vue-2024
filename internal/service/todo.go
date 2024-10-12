// Package service provides the business logic for the todo endpoint.
package service

import (
	"github.com/zuu-development/fullstack-examination-2024/internal/model"
	"github.com/zuu-development/fullstack-examination-2024/internal/repository"
)

// Todo is the service for the todo endpoint.
type Todo interface {
	Create(task string, priority int) (*model.Todo, error)
	Update(id int, task string, status model.Status, priority int) (*model.Todo, error)
	Delete(id int) error
	Find(id int) (*model.Todo, error)
	FindAll() ([]*model.Todo, error)
	FindAllFiltered(task, status string, priority string) (map[string][]*model.Todo, error)
}

type todo struct {
	todoRepository repository.Todo
}

// NewTodo creates a new Todo service.
func NewTodo(r repository.Todo) Todo {
	return &todo{r}
}

func (t *todo) Create(task string, priority int) (*model.Todo, error) {
	todo := model.NewTodo(task, priority)
	if err := t.todoRepository.Create(todo); err != nil {
		return nil, err
	}
	return todo, nil
}

func (t *todo) Update(id int, task string, status model.Status, priority int) (*model.Todo, error) {
	todo := model.NewUpdateTodo(id, task, status, priority) // Pass priority here

	// Fetch the current todo values
	currentTodo, err := t.Find(id)
	if err != nil {
		return nil, err
	}

	// Use existing values if new values are empty
	if todo.Task == "" {
		todo.Task = currentTodo.Task
	}
	if todo.Status == "" {
		todo.Status = currentTodo.Status
	}
	if todo.Priority == 0 { // If priority not set, use the current one
		todo.Priority = currentTodo.Priority
	}

	if err := t.todoRepository.Update(todo); err != nil {
		return nil, err
	}
	return todo, nil
}

func (t *todo) Delete(id int) error {
	if err := t.todoRepository.Delete(id); err != nil {
		return err
	}
	return nil
}

func (t *todo) Find(id int) (*model.Todo, error) {
	todo, err := t.todoRepository.Find(id)
	if err != nil {
		return nil, err
	}
	return todo, nil
}

func (t *todo) FindAll() ([]*model.Todo, error) {
	todo, err := t.todoRepository.FindAll()
	if err != nil {
		return nil, err
	}
	return todo, nil
}

func (t *todo) FindAllFiltered(task, status, priority string) (map[string][]*model.Todo, error) {
	// Call the repository layer to apply filtering
	todos, err := t.todoRepository.FindAllFiltered(task, status, priority)
	if err != nil {
		return nil, err
	}

	// Separate tasks into completed and incomplete lists
	incompleteTasks := []*model.Todo{}
	completedTasks := []*model.Todo{}

	for _, todo := range todos {
		if todo.Status == model.Done {
			completedTasks = append(completedTasks, todo)
		} else {
			incompleteTasks = append(incompleteTasks, todo)
		}
	}

	// Return both lists in a map
	return map[string][]*model.Todo{
		"incompleteTasks": incompleteTasks,
		"completedTasks":  completedTasks,
	}, nil
}
