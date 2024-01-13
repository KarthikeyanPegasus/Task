package main

import (
	"sample-project/cron"
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

	cronServer := cron.NewCronServer(db)
	go cronServer.NewCronJob()

	server := handler.NewTaskServer(db)

	route.InitRoute(server)
}
