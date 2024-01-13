package Bloc

import (
	"math/rand"
	"strconv"
	"time"
)

func generateId() string {
	sec := time.Now().Nanosecond()
	TaskId := randStringRunes(2) + strconv.Itoa(sec)
	return TaskId
}

var letterRunes = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")

func randStringRunes(n int) string {
	b := make([]rune, n)
	for i := range b {
		b[i] = letterRunes[rand.Intn(len(letterRunes))]
	}
	return string(b)
}
