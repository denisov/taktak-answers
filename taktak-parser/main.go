package main

import (
	"github.com/denisov/taktak-answers/models"
	"log"
	"os"
)

func main() {
	log.SetOutput(os.Stdout)

	// TODO вынести в конфиг
	models.InitDb("/home/andrey/taktak/answers.db")

	parseSolutions()
}
