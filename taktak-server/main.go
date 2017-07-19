package main

import (
	"encoding/json"
	"github.com/denisov/taktak-answers/models"
	"github.com/jinzhu/now"
	"log"
	"net/http"
	"os"
	"strconv"
)

const PORT = 8001

type response struct {
	ThisWeekAnswers int `json:"this_week_answers"`
}

func main() {
	log.SetOutput(os.Stdout)
	models.InitDb("/home/andrey/taktak/answers.db")
	now.FirstDayMonday = true

	http.HandleFunc("/", func(writer http.ResponseWriter, request *http.Request) {
		writer.Header().Set("Content-type", "application/json")
		writer.Header().Set("Access-Control-Allow-Origin", "*")

		thisWeekAnswers, err := models.SolutionsGetCountFrom(now.BeginningOfWeek())
		if err != nil {
			log.Fatalln(err)
		}

		json.NewEncoder(writer).Encode(response{thisWeekAnswers})
		log.Println("req")
	})

	log.Println("TakTak answer server ready for connections on " + strconv.Itoa(PORT))
	log.Fatal(http.ListenAndServe(":"+strconv.Itoa(PORT), nil))
	// не daemon режим
}
