package Bloc

import (
	"database/sql"
	"sample-project/cache"
	"sample-project/models"
	"strconv"
	"time"
)

type TaskServer struct {
	db *sql.DB
	c  *cache.Cache
}

func NewTaskServer(db *sql.DB, c *cache.Cache) *TaskServer {
	return &TaskServer{db: db, c: c}
}

func (s *TaskServer) GetTask(request *models.GetTaskRequest) (*models.GetTaskResponse, error) {
	// check cache first
	if task, err := s.c.GetCache(request.ID); err == nil {
		return &models.GetTaskResponse{Task: *task}, nil
	}

	query := `SELECT id, title, description, priority, datetime FROM "Tasks".tasks WHERE id = $1;`

	var task models.Task
	err := s.db.QueryRow(query, request.ID).Scan(
		&task.ID,
		&task.Title,
		&task.Description,
		&task.Priority,
		&task.DateTime,
	)
	if err != nil {
		return &models.GetTaskResponse{}, err
	}

	// set cache
	err = s.c.SetCache(request.ID, &task)
	if err != nil {
		return &models.GetTaskResponse{}, err
	}

	return &models.GetTaskResponse{Task: task}, nil
}

func (s *TaskServer) CreateTask(request *models.CreateTaskRequest) (*models.CreateTaskResponse, error) {
	// Handle validation here

	datetime, err := time.Parse(time.RFC3339, request.DateTime)
	if err != nil {
		return &models.CreateTaskResponse{}, err
	}

	id := generateId()
	_, err = s.db.Exec(
		`INSERT INTO "Tasks".tasks (id, title, description, priority, datetime) VALUES ($1, $2, $3,$4, $5)`,
		id,
		request.Title,
		request.Description,
		request.Priority,
		datetime,
	)
	if err != nil {
		return &models.CreateTaskResponse{}, err
	}

	err = s.c.SetCache(id, &models.Task{
		ID:          id,
		Title:       request.Title,
		Description: request.Description,
		Priority:    request.Priority,
		DateTime:    datetime,
	})

	return &models.CreateTaskResponse{ID: id}, nil
}

func (s *TaskServer) UpdateTask(request *models.UpdateTaskRequest) (*models.UpdateTaskResponse, error) {

	//Have to handle BLOC here
	params := make([]interface{}, 0, len(request.UpdateMask)+1)
	updateRequestMap := make(map[string]interface{})
	for _, mask := range request.UpdateMask {
		if mask == "title" {
			updateRequestMap["title"] = request.Title
		}
		if mask == "description" {
			updateRequestMap["description"] = request.Description
		}
		if mask == "priority" {
			updateRequestMap["priority"] = request.Priority
		}
		if mask == "time" {
			updateRequestMap["time"] = request.DateTime
		}
	}

	query := `UPDATE "Tasks".tasks SET `
	for i, mask := range request.UpdateMask {
		if i == len(request.UpdateMask)-1 {
			query += mask + ` = $` + strconv.Itoa(i+1)
		} else {
			query += mask + ` = $` + strconv.Itoa(i+1) + `, `
		}
		params = append(params, updateRequestMap[mask])
	}
	query += ` WHERE id = $` + strconv.Itoa(len(request.UpdateMask)+1) + `;`
	params = append(params, request.ID)

	_, err := s.db.Exec(query, params...)
	if err != nil {
		return &models.UpdateTaskResponse{}, err
	}

	s.c.DelCache(request.ID)

	return &models.UpdateTaskResponse{ID: request.ID}, nil
}

func (s *TaskServer) DeleteTask(request *models.DeleteTaskRequest) error {

	_, err := s.db.Exec(`DELETE FROM "Tasks".tasks WHERE id = $1`, request.ID)
	if err != nil {
		return err
	}

	return nil
}
