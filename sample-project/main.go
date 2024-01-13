package main

import (
	"sample-project/Bloc"
	"sample-project/Bloc/cron"
	"sample-project/cache"
	"sample-project/database"
	"sample-project/handler"
	"sample-project/route"
)

func main() {
	redis := cache.RedisInit()
	db, err := database.Init()
	if err != nil {
		panic(err)
	}

	cronServer := cron.NewCronServer(db)
	go cronServer.NewCronJob()

	cache := cache.NewCache(redis)
	BlocServer := Bloc.NewTaskServer(db, cache)
	server := handler.NewHandler(BlocServer)

	route.InitRoute(server)
}
