package main

import (
	"cms-api/internal/app"
)

// version is set at build time
var version = "dev"

func main() {
	app.Run(version)
}
