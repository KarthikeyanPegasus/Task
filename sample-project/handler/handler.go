package handler

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"sample-project/models"
	"strconv"
	"time"
)

type TaskServer struct {
	db *sql.DB
}

func NewTaskServer(db *sql.DB) *TaskServer {
	return &TaskServer{db: db}
}

func (s *TaskServer) GetTask(writer http.ResponseWriter, request *http.Request) {
	var getTaskRequest models.GetTaskRequest
	if err := json.NewDecoder(request.Body).Decode(&getTaskRequest); err != nil {
		http.Error(writer, "Invalid JSON format", http.StatusBadRequest)
		return
	}

	// Handle BLOC here

	query := `SELECT id, title, description, priority, datetime FROM "Tasks".tasks WHERE id = $1;`

	var task models.Task
	err := s.db.QueryRow(query, getTaskRequest.ID).Scan(
		&task.ID,
		&task.Title,
		&task.Description,
		&task.Priority,
		&task.DateTime,
	)
	if err != nil {
		http.Error(writer, "Internal server error", http.StatusInternalServerError)
		return
	}

	if err := json.NewEncoder(writer).Encode(models.GetTaskResponse{Task: task}); err != nil {
		http.Error(writer, "Internal server error", http.StatusInternalServerError)
		return
	}
}

func (s *TaskServer) CreateTask(writer http.ResponseWriter, request *http.Request) {
	var createTaskRequest models.CreateTaskRequest
	if err := json.NewDecoder(request.Body).Decode(&createTaskRequest); err != nil {
		http.Error(writer, "Invalid JSON format", http.StatusBadRequest)
		return
	}

	// Handle BLOC here

	time, err := time.Parse(time.RFC3339, createTaskRequest.DateTime)
	if err != nil {
		http.Error(writer, "Invalid date time format", http.StatusBadRequest)
		return
	}

	id := generateId()
	_, err = s.db.Exec(
		`INSERT INTO "Tasks".tasks (id, title, description, priority, datetime) VALUES ($1, $2, $3,$4, $5)`,
		id,
		createTaskRequest.Title,
		createTaskRequest.Description,
		createTaskRequest.Priority,
		time,
	)
	if err != nil {
		http.Error(writer, "Internal server error", http.StatusInternalServerError)
		return
	}

	if err := json.NewEncoder(writer).Encode(models.CreateTaskResponse{ID: id}); err != nil {
		http.Error(writer, "Internal server error", http.StatusInternalServerError)
		return
	}
}

func (s *TaskServer) UpdateTask(writer http.ResponseWriter, request *http.Request) {
	var updateTaskRequest models.UpdateTaskRequest
	if err := json.NewDecoder(request.Body).Decode(&updateTaskRequest); err != nil {
		http.Error(writer, "Invalid JSON format", http.StatusBadRequest)
		return
	}
	//Have to handle BLOC here
	params := make([]interface{}, 0, len(updateTaskRequest.UpdateMask)+1)
	updateRequestMap := make(map[string]interface{})
	for _, mask := range updateTaskRequest.UpdateMask {
		if mask == "title" {
			updateRequestMap["title"] = updateTaskRequest.Title
		}
		if mask == "description" {
			updateRequestMap["description"] = updateTaskRequest.Description
		}
		if mask == "priority" {
			updateRequestMap["priority"] = updateTaskRequest.Priority
		}
		if mask == "time" {
			updateRequestMap["time"] = updateTaskRequest.DateTime
		}
	}

	query := `UPDATE "Tasks".tasks SET `
	for i, mask := range updateTaskRequest.UpdateMask {
		if i == len(updateTaskRequest.UpdateMask)-1 {
			query += mask + ` = $` + strconv.Itoa(i+1)
		} else {
			query += mask + ` = $` + strconv.Itoa(i+1) + `, `
		}
		params = append(params, updateRequestMap[mask])
	}
	query += ` WHERE id = $` + strconv.Itoa(len(updateTaskRequest.UpdateMask)+1) + `;`
	params = append(params, updateTaskRequest.ID)

	_, err := s.db.Exec(query, params...)
	if err != nil {
		json.NewEncoder(writer).Encode(err)
		http.Error(writer, "Internal server error", http.StatusInternalServerError)
		return
	}

	if err := json.NewEncoder(writer).Encode(models.UpdateTaskResponse{ID: updateTaskRequest.ID}); err != nil {
		http.Error(writer, "Internal server error", http.StatusInternalServerError)
		return
	}
}

func (s *TaskServer) DeleteTask(writer http.ResponseWriter, request *http.Request) {
	var deleteTaskRequest models.DeleteTaskRequest
	if err := json.NewDecoder(request.Body).Decode(&deleteTaskRequest); err != nil {
		http.Error(writer, "Invalid JSON format", http.StatusBadRequest)
		return
	}

	// Handle BLOC here

	_, err := s.db.Exec(`DELETE FROM "Tasks".tasks WHERE id = $1`, deleteTaskRequest.ID)
	if err != nil {
		http.Error(writer, "Internal server error", http.StatusInternalServerError)
		return
	}

	http.Error(writer, "successfully deleted", http.StatusOK)
	return
}
