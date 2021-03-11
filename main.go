package main

import (
	"github.com/99-66/go-gin-token-server/config"
	"github.com/99-66/go-gin-token-server/routes"
	"log"
)

var err error

func main() {
	config.REDIS, err = config.InitRedis()
	if err != nil {
		panic(err)
	}
	defer config.REDIS.Close()

	r := routes.InitRouter()
	log.Fatal(r.Run())
}