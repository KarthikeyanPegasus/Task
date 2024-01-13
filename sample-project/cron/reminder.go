package cron

import (
	"database/sql"
	"fmt"
	"sample-project/models"
	"time"
)

type CronServer struct {
	db *sql.DB
}

func NewCronServer(db *sql.DB) *CronServer {
	return &CronServer{db: db}
}

func (s *CronServer) NewCronJob() {
	for {
		if err := s.DoReminder(); err != nil {
			break
		}
		time.Sleep(5 * time.Minute)
	}
}

func (s *CronServer) DoReminder() error {
	// list all tasks that need to be reminded
	query := `SELECT id, title, description, priority, datetime FROM tasks WHERE datetime < $1;`

	rows, err := s.db.Query(query, time.Now().Add(5*time.Minute).Format(time.RFC3339))
	if err != nil {
		return err
	}

	for rows.Next() {
		var task models.Task
		if err := rows.Scan(
			&task.ID,
			&task.Title,
			&task.Description,
			&task.Priority,
			&task.DateTime,
		); err != nil {
			return err
		}
		fmt.Print(fmt.Sprintf("Task %s is due in 5 minutes", task.Title),
			"\n",
			fmt.Sprintf("Title: %s", task.Title))
		// we can send email to user
	}
	return nil
}
