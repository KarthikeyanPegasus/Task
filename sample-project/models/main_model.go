package models

import (
	"time"
)

type Task struct {
	ID          string    `json:"id"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	Priority    Priority  `json:"priority"`
	DateTime    time.Time `json:"time"`
}

type Priority int

const (
	High Priority = iota
	Medium
	Low
)

type CreateTaskRequest struct {
	ID          string   `json:"id,omitempty"`
	Title       string   `json:"title"`
	Description string   `json:"description"`
	Priority    Priority `json:"priority"`
	DateTime    string   `json:"time"`
}

type CreateTaskResponse struct {
	ID string `json:"id"`
}

type Mask []string

type UpdateTaskRequest struct {
	ID          string   `json:"id"`
	UpdateMask  Mask     `json:"update_mask"`
	Title       string   `json:"title,omitempty"`
	Description string   `json:"description,omitempty"`
	Priority    Priority `json:"priority,omitempty"`
	DateTime    string   `json:"time,omitempty"`
}

type UpdateTaskResponse struct {
	ID string `json:"id"`
}

type DeleteTaskRequest struct {
	ID string `json:"id"`
}

type GetTaskRequest struct {
	ID string `json:"id"`
}

type GetTaskResponse struct {
	Task Task `json:"task"`
}

type ListTaskRequest struct {
	Ids []string `json:"ids"`
}

type ListTaskResponse struct {
	Tasks []*Task `json:"Tasks"`
}

type Tasks []Task

type TaskService interface {
	CreateTask(req CreateTaskRequest) (CreateTaskResponse, error)
	UpdateTask(req UpdateTaskRequest) (UpdateTaskResponse, error)
	DeleteTask(req DeleteTaskRequest) error
	GetTask(req GetTaskRequest) (GetTaskResponse, error)
}
