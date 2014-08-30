package main

import (
	godis "github.com/simonz05/godis/redis"
	"os"
	"time"
)

// Gets redis instance from configuration parameters.
// TODO: override config if specific ENV variable is in place, eg: heroku
func redis_instance() (redis *godis.Client) {
	host := config.GetStringDefault("redis.netaddress", "tcp:localhost:6379")
	db := config.GetIntDefault("redis.database", 0)
	passwd := config.GetStringDefault("redis.password", "")

	redis = godis.New(host, db, passwd)

	// Try to set a "dummy" value, panic if redis is not accessible.
	var error_on_set = redis.Set("INIT", time.Now().UTC())
	if error_on_set != nil {
		log.Critical("Cannot access Redis server")
		os.Exit(2)
	}

	return
}
