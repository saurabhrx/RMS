package main

import (
	"RMS/database"
	"RMS/routes"
	"fmt"
	"github.com/joho/godotenv"
	"github.com/sirupsen/logrus"
	"net/http"
	"os"
)

func main() {
	if err := godotenv.Load(); err != nil {
		fmt.Println("error in loading env file")
		return
	}

	host := os.Getenv("DB_HOST")
	post := os.Getenv("DB_PORT")
	databaseName := os.Getenv("DB_NAME")
	user := os.Getenv("DB_USER")
	password := os.Getenv("DB_PASSWORD")

	if err := database.ConnectToDB(host, post, user, password, databaseName); err != nil {
		logrus.Panicf("failed to connect to database : %+v", err)
	}
	fmt.Println("database connected")

	srv := routes.SetupTodoRoutes()

	if srvErr := http.ListenAndServe(":8080", srv); srvErr != nil {
		logrus.Panicf("failed to connect to server %+v", srvErr)
		return
	}

	if DBCloseErr := database.CloseDBConnection(); DBCloseErr != nil {
		logrus.Panicf("failed to close database %+v", DBCloseErr)
		return
	}

}
