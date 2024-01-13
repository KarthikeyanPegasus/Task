package handler

import (
	"encoding/json"
	"net/http"
	"sample-project/Bloc"
	"sample-project/models"
)

type Handler struct {
	task *Bloc.TaskServer
}

func NewHandler(task *Bloc.TaskServer) *Handler {
	return &Handler{task: task}
}

func (s *Handler) GetTask(writer http.ResponseWriter, request *http.Request) {
	var getTaskRequest models.GetTaskRequest
	if err := json.NewDecoder(request.Body).Decode(&getTaskRequest); err != nil {
		http.Error(writer, "Invalid JSON format", http.StatusBadRequest)
		return
	}

	// Handle BLOC here

	resp, err := s.task.GetTask(&getTaskRequest)
	if err != nil {
		http.Error(writer, "Internal server error", http.StatusInternalServerError)
		return
	}

	if err := json.NewEncoder(writer).Encode(models.GetTaskResponse{Task: resp.Task}); err != nil {
		http.Error(writer, "Internal server error", http.StatusInternalServerError)
		return
	}
}

func (s *Handler) CreateTask(writer http.ResponseWriter, request *http.Request) {
	var createTaskRequest models.CreateTaskRequest
	if err := json.NewDecoder(request.Body).Decode(&createTaskRequest); err != nil {
		http.Error(writer, "Invalid JSON format", http.StatusBadRequest)
		return
	}

	// Handle BLOC here

	resp, err := s.task.CreateTask(&createTaskRequest)
	if err != nil {
		http.Error(writer, "Internal server error", http.StatusInternalServerError)
		return
	}

	if err := json.NewEncoder(writer).Encode(models.CreateTaskResponse{ID: resp.ID}); err != nil {
		http.Error(writer, "Internal server error", http.StatusInternalServerError)
		return
	}
}

func (s *Handler) UpdateTask(writer http.ResponseWriter, request *http.Request) {
	var updateTaskRequest models.UpdateTaskRequest
	if err := json.NewDecoder(request.Body).Decode(&updateTaskRequest); err != nil {
		http.Error(writer, "Invalid JSON format", http.StatusBadRequest)
		return
	}
	resp, err := s.task.UpdateTask(&updateTaskRequest)
	if err != nil {
		json.NewEncoder(writer).Encode(err)
		http.Error(writer, "Internal server error", http.StatusInternalServerError)
		return
	}

	if err := json.NewEncoder(writer).Encode(models.UpdateTaskResponse{ID: resp.ID}); err != nil {
		http.Error(writer, "Internal server error", http.StatusInternalServerError)
		return
	}
}

func (s *Handler) DeleteTask(writer http.ResponseWriter, request *http.Request) {
	var deleteTaskRequest models.DeleteTaskRequest
	if err := json.NewDecoder(request.Body).Decode(&deleteTaskRequest); err != nil {
		http.Error(writer, "Invalid JSON format", http.StatusBadRequest)
		return
	}

	err := s.task.DeleteTask(&deleteTaskRequest)
	if err != nil {
		http.Error(writer, "Internal server error", http.StatusInternalServerError)
		return
	}

	http.Error(writer, "successfully deleted", http.StatusOK)
	return
}
