package main

import (
	"github.com/fs111/simpleconfig"
	godis "github.com/simonz05/godis/redis"
)

// Gets redis instance from configuration parameters.
// TODO: panic and quit if server is not accessible
// TODO: override config if specific ENV variable is in place, eg: heroku
func redis_instance(*simpleconfig.Config) *godis.Client {
	host := config.GetStringDefault("redis.netaddress", "tcp:localhost:6379")
	db := config.GetIntDefault("redis.database", 0)
	passwd := config.GetStringDefault("redis.password", "")

	return godis.New(host, db, passwd)
}
