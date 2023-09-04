package main

import (
	"log"

	"crs_v2/actions"
)

func main() {
	actions.OpenDatabaseConnection()
	defer actions.DatabaseConnection.Close()

	cronJobs := actions.StartCronScheduler()
	defer cronJobs.Stop()
	//start the app
	app := actions.App()

	if err := app.Serve(); err != nil {
		log.Fatal(err)
	}
}
