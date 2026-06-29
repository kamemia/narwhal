package main

import (
	"log"

	"narwhal/internal/server"
	"narwhal/internal/upload"
)

func main() {
	uploadManager := upload.NewManager("uploads")

	router := server.NewRouter(uploadManager)
	server := server.NewServer(8080, router)

	if err := server.Run(); err != nil {
		log.Fatal(err)
	}
}
