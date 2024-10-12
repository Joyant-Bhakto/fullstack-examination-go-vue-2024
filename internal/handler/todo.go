package handler

import (
    "net/http"
    "github.com/labstack/echo/v4" 
    "github.com/zuu-development/fullstack-examination-2024/internal/errors"
    "github.com/zuu-development/fullstack-examination-2024/internal/model"
    "github.com/zuu-development/fullstack-examination-2024/internal/service"
)

// TodoHandler is the request handler for the todo endpoint.
type TodoHandler interface {
	Create(c echo.Context) error
	Update(c echo.Context) error
	Delete(c echo.Context) error
	Find(c echo.Context) error
	FindAll(c echo.Context) error
}

type todoHandler struct {
	service service.Todo
}

// NewTodo returns a new instance of the todo handler.
func NewTodo(s service.Todo) TodoHandler {
	return &todoHandler{service: s}
}

// CreateRequest is the request parameter for creating a new todo
type CreateRequest struct {
	Task string `json:"task" validate:"required"`
}

// FindAllRequest represents the query parameters for searching todos
type FindAllRequest struct {
	Task   string       `query:"task"`   // Query parameter for task
	Status model.Status `query:"status"` // Query parameter for status
}

// Create creates a new todo
func (t *todoHandler) Create(c echo.Context) error {
	var req CreateRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest,
			ResponseError{Errors: []Error{{Code: errors.CodeBadRequest, Message: err.Error()}}})
	}

	// Check if the Task field is empty
	if req.Task == "" {
		return c.JSON(http.StatusBadRequest,
			ResponseError{Errors: []Error{{Code: errors.CodeBadRequest, Message: "task is required"}}})
	}

	// Keeping someIntValue as per your requirement
	someIntValue := 1 // Replace with actual logic to get the value
	todo, err := t.service.Create(req.Task, someIntValue) // Pass the required int value
	if err != nil {
		return c.JSON(http.StatusInternalServerError,
			ResponseError{Errors: []Error{{Code: errors.CodeInternalServerError, Message: err.Error()}}})
	}

	return c.JSON(http.StatusCreated, ResponseData{Data: todo})
}
// UpdateRequest is the request parameter for updating a todo
type UpdateRequest struct {
	UpdateRequestBody
	UpdateRequestPath
}

// UpdateRequestBody is the request body for updating a todo
type UpdateRequestBody struct {
	Task   string       `json:"task,omitempty"`
	Status model.Status `json:"status,omitempty"`
}

// UpdateRequestPath is the request parameter for updating a todo
type UpdateRequestPath struct {
	ID int `param:"id" validate:"required"`
}

// Update updates an existing todo
func (t *todoHandler) Update(c echo.Context) error {
	var req UpdateRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest,
			ResponseError{Errors: []Error{{Code: errors.CodeBadRequest, Message: err.Error()}}})
	}

	// Keeping someIntValue as per your requirement
	someIntValue := 1 // Replace with actual logic to get the value
	todo, err := t.service.Update(req.ID, req.Task, req.Status, someIntValue) // Pass the required int value
	if err != nil {
		if err == model.ErrNotFound {
			return c.JSON(http.StatusNotFound,
				ResponseError{Errors: []Error{{Code: errors.CodeNotFound, Message: "todo not found"}}})
		}
		return c.JSON(http.StatusInternalServerError,
			ResponseError{Errors: []Error{{Code: errors.CodeInternalServerError, Message: err.Error()}}})
	}

	return c.JSON(http.StatusOK, ResponseData{Data: todo})
}
// DeleteRequest is the request parameter for deleting a todo
type DeleteRequest struct {
	ID int `param:"id" validate:"required"`
}

// Delete deletes a todo
func (t *todoHandler) Delete(c echo.Context) error {
	var req DeleteRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest,
			ResponseError{Errors: []Error{{Code: errors.CodeBadRequest, Message: err.Error()}}})
	}

	if err := t.service.Delete(req.ID); err != nil {
		if err == model.ErrNotFound {
			return c.JSON(http.StatusNotFound,
				ResponseError{Errors: []Error{{Code: errors.CodeNotFound, Message: "todo not found"}}})
		}
		return c.JSON(http.StatusInternalServerError,
			ResponseError{Errors: []Error{{Code: errors.CodeInternalServerError, Message: err.Error()}}})
	}
	return c.NoContent(http.StatusNoContent)
}

// FindRequest is the request parameter for finding a todo
type FindRequest struct {
	ID int `param:"id" validate:"required"`
}

// Find finds a todo
func (t *todoHandler) Find(c echo.Context) error {
	var req FindRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest,
			ResponseError{Errors: []Error{{Code: errors.CodeBadRequest, Message: err.Error()}}})
	}

	res, err := t.service.Find(req.ID)
	if err != nil {
		if err == model.ErrNotFound {
			return c.JSON(http.StatusNotFound,
				ResponseError{Errors: []Error{{Code: errors.CodeNotFound, Message: "todo not found"}}})
		}
		return c.JSON(http.StatusInternalServerError,
			ResponseError{Errors: []Error{{Code: errors.CodeInternalServerError, Message: err.Error()}}})
	}
	return c.JSON(http.StatusOK, ResponseData{Data: res})
}


func (t *todoHandler) FindAll(c echo.Context) error {
	var req FindAllRequest
	// Bind query parameters
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest,
			ResponseError{Errors: []Error{{Code: errors.CodeBadRequest, Message: err.Error()}}})
	}

	// Call the service to filter todos based on task and status, and sort by priority
	res, err := t.service.FindAllFiltered(req.Task, string(req.Status), "priority")
	if err != nil {
		return c.JSON(http.StatusInternalServerError,
			ResponseError{Errors: []Error{{Code: errors.CodeInternalServerError, Message: err.Error()}}})
	}
	return c.JSON(http.StatusOK, ResponseData{Data: res})
}
