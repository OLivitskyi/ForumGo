package main

import (
	"Forum/internal/app/server"
	"log"
)

func main() {
	// Create a new configuration instance
	// Read configuration from file
	config := server.NewConfig()
	if err := config.ReadConfig(); err != nil {
		log.Fatal(err)
	}

	// Start the server with the obtained configuration
	log.Fatal(server.Start(config))
}


/*
To DO:

*login logic

*authentication - sessions and cookies

*post creation

*create categories

*handle UNIQUE constraint failure

*sanitize and validatge input

*multiple hompage logs?

*/
