package main

import (
	godis "github.com/simonz05/godis/redis"
	"os"
	"time"
  "strconv"
)

// Gets redis instance from configuration parameters.
func redis_instance() (redis *godis.Client) {
	host := os.Getenv("REDIS_HOST")
	db, _ := strconv.Atoi(os.Getenv("REDIS_DB"))
	passwd := os.Getenv("REDIS_PASS")

	redis = godis.New(host, db, passwd)

	var error_on_set = redis.Set("INIT", time.Now().UTC())
	if error_on_set != nil {
		log.Critical("Cannot access Redis server")
		os.Exit(2)
	}

	return
}
