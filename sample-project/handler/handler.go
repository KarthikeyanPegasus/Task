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

	query := `SELECT id, title, description, priority, datetime FROM tasks WHERE id = $1;`

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

func (s *TaskServer) ListTasks(writer http.ResponseWriter, request *http.Request) {
	var listTaskRequest models.ListTaskRequest
	if err := json.NewDecoder(request.Body).Decode(&listTaskRequest); err != nil {
		http.Error(writer, "Invalid JSON format", http.StatusBadRequest)
		return
	}

	// Handle BLOC here
	ids := ""
	for i, id := range listTaskRequest.Ids {
		if i == len(listTaskRequest.Ids)-1 {
			ids += id
			break
		}
		ids += id + `, `
	}

	query := `SELECT id, title, description, priority, datetime FROM tasks WHERE ID in ($1);`

	rows, err := s.db.Query(query, ids)
	if err != nil {
		http.Error(writer, "Internal server error", http.StatusInternalServerError)
		return
	}

	var tasks []*models.Task
	for rows.Next() {
		var task *models.Task
		err := rows.Scan(
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
		tasks = append(tasks, task)
	}

	if err := json.NewEncoder(writer).Encode(models.ListTaskResponse{Tasks: tasks}); err != nil {
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
		//"INSERT INTO Tasks.tasks (id, title, description, priority, datetime) VALUES ($1, $2, $3,$4, $5)",
		"select * from tasks",
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

	query := `UPDATE tasks SET `
	for i, mask := range updateTaskRequest.UpdateMask {
		if i == len(updateTaskRequest.UpdateMask)-1 {
			query += mask + ` = $` + strconv.Itoa(i+1)
			break
		}
		query += mask + ` = $` + strconv.Itoa(i+1) + `, `
	}
	query += ` WHERE id = $` + updateTaskRequest.ID + `;`

	_, err := s.db.Exec(
		query,
		updateTaskRequest.Title,
		updateTaskRequest.Description,
		updateTaskRequest.Priority,
		updateTaskRequest.DateTime,
		updateTaskRequest.ID,
	)
	if err != nil {
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

	_, err := s.db.Exec("DELETE FROM tasks WHERE id = $1", deleteTaskRequest.ID)
	if err != nil {
		http.Error(writer, "Internal server error", http.StatusInternalServerError)
		return
	}

	http.Error(writer, "successfully deleted", http.StatusOK)
	return
}
