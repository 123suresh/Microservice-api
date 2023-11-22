package main

import (
	"os"

	"example.com/dynamicWordpressBuilding/internal/controller"
	"example.com/dynamicWordpressBuilding/internal/repository"
	"example.com/dynamicWordpressBuilding/internal/service"
	"github.com/joho/godotenv"
	"github.com/sirupsen/logrus"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		logrus.Fatalf("Error getting env, not coming through %v", err)
	}
	logrus.Info("Successfully loaded env file")
	repo := repository.NewRepo()
	svc := service.NewService(repo)
	ctl := controller.NewController(svc)
	port := os.Getenv("PORT")
	if port == "" {
		port = "8081"
	}
	err = ctl.Router.Run(":" + port)
}
