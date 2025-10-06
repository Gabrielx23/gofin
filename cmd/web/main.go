package main

import (
	"log"
	"net/http"

	"gofin/internal/container"
)

const ServerPort = ":8080"

func main() {
	container, err := container.NewContainerWithDefaultConfig()
	if err != nil {
		log.Fatalf("Failed to initialize container: %v", err)
	}
	defer container.DB.Close()

	mux := http.NewServeMux()

	_, err = NewRouter(container, mux)
	if err != nil {
		log.Fatalf("Failed to initialize router: %v", err)
	}

	Start(ServerPort, mux)
}
