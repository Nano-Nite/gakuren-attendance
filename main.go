package main

import (
	"log"

	"gakuren-system.com/pkg/auth"
	"gakuren-system.com/pkg/db"
	"gakuren-system.com/pkg/helper"
	v1 "gakuren-system.com/pkg/route/v1"
)

func main() {
	// Connect to the database
	log.Println("Connecting to the database")
	db.ConnectDB()

	// Run JWT Service
	log.Println("Starting JWT Service")
	if err := auth.NewJWTService(); err != nil {
		log.Fatal(err)
	}

	// init variable
	helper.InitVariableDB()

	// Setup and launch the Fiber routes
	log.Println("Launching Fiber Routes")
	if "v1" == helper.API_VERSION {
		v1.SetupRoutes()
	}
}
