package main

import (
	"github.com/broemp/sshChat/config"
	"github.com/broemp/sshChat/models"
)

func main() {
	config.LoadConfig()
	app := models.NewApp()
	app.Start()
}
