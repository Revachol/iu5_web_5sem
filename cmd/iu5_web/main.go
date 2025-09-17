package main

import (
	"log"

	app "github.com/Revachol/iu5_web_5sem/internal/api"
)

func main() {
	log.Println("App started on http://127.0.0.1:8080/")
	app.StartServer()
	log.Println("App finished")
}
