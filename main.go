package main

import (
	api "bcraft/api"
	"bcraft/db"
)

func init() {
	InitConfig()

	db.Init()
}

func main() {
	api.Init()
}
