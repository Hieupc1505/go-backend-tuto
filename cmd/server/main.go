package main

import (
	"hieupc05.github/backend-server/internal/initialize"
)

func main() {
	// r := routers.NewRoute()
	// r.Run() // listen and serve on 0.0.0.0:8080 (for windows "localhost:8080")

	initialize.Run()

}
