package main

import "github.com/broemp/sshChat/models"

func main() {
	app := models.NewApp()
	app.Start()
}
