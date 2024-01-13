package route

import (
	"fmt"
	"github.com/go-chi/chi"
	"net/http"
	"sample-project/handler"
)

func InitRoute(server *handler.TaskServer) {
	r := chi.NewRouter()
	r.Get("/task", server.GetTask)
	r.Get("/tasks", server.ListTasks)
	r.Post("/tasks", server.CreateTask)
	r.Put("/tasks", server.UpdateTask)
	r.Delete("/tasks", server.DeleteTask)

	serverAddr := "localhost:8080"
	fmt.Printf("Server is running on http://%s\n", serverAddr)
	err := http.ListenAndServe(serverAddr, r)
	if err != nil {
		return
	}
}
