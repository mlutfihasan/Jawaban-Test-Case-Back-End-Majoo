package main

import (
	"GOMISPLUS/controllers"
)

var server = controllers.Server{}

func main() {
	server.Initialize()

	server.Run(":8089") // kalo ubah jangan dicommit
}
