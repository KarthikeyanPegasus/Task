package main

import (
	"sample-project/database"
	"sample-project/handler"
	"sample-project/route"
)

func main() {
	// init cache and database
	db, err := database.Init()
	if err != nil {
		panic(err)
	}

	server := handler.NewTaskServer(db)

	route.InitRoute(server)
}
