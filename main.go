package main

import (
	"go-server/config"
	"go-server/router"
)

func main() {
	config.InitDB()
	r := router.InitRouter()
	r.Run(":5200")
}