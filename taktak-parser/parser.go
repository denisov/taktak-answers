package main

import (
	"crypto/tls"
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"github.com/denisov/taktak-answers/models"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"
)

// parseSolutions парсит новые решения
func parseSolutions() {
	// перебираем первые 10 страниц. На странице решения упорядочены по последнему комментарию, поэтому не понятно когда нужн остановиться
	// 10 страниц хватит в любом случае
	for page := 1; page < 10; page++ {
		parsePage(page)
	}
}

// parsePage парсит страницу решений
// на странице идут решения упорядоченные по дате последнего комментария по убыванию (новые сверху)
// однако в качестве решения учитывается первый комментарий
func parsePage(page int) {
	urlToParse := fmt.Sprintf("https://taktaktak.ru/problems/solved/?page=%d&ajax=2", page)
	log.Println(urlToParse)

	resp := getResponse(urlToParse)
	doc, err := goquery.NewDocumentFromResponse(resp)
	if err != nil {
		log.Fatalf("Не могу создать документ для парсинга URL: %s. %s", urlToParse, err)
	}
	time.Sleep(2 * time.Second)
	doc.Find(".item").Each(func(i int, selection *goquery.Selection) {
		title := selection.Find(".content .title a").First()
		problemPageUrl, exists := title.Attr("href")
		if !exists {
			html, _ := title.Html()
			log.Printf("Не могу найти аттрибут 'href' в '%s'", html)
			return
		}

		urlSplitted := strings.Split(problemPageUrl, "/")
		problemId, err := strconv.Atoi(urlSplitted[len(urlSplitted)-1])
		if err != nil {
			log.Fatal(err)
		}

		exists, err = models.SolutionCheckExists(problemId)
		if err != nil {
			log.Fatal(err)
		}

		problemMsg := fmt.Sprintf("ProblemId: %d\t", problemId)
		if exists {
			log.Println(problemMsg + "дата решения в базе есть, поропускаем")
		} else {
			log.Println(problemMsg + "даты решения в базе нет, парсим страницу проблемы")
			parseProblemPage(problemId)
		}
	})
}

// parseProblemPage получает дату первого решения со страницы проблемы
func parseProblemPage(problemId int) {
	urlToParse := fmt.Sprintf("https://taktaktak.ru/problem/%d", problemId)

	resp := getResponse(urlToParse)
	doc, err := goquery.NewDocumentFromResponse(resp)
	time.Sleep(2 * time.Second)
	if err != nil {
		log.Fatalf("Не могу создать документ для парсинга URL: %s. %s", urlToParse, err)
	}
	answerDateText := doc.Find(".answer").First().Find(".date span").First().Text()
	if answerDateText == "" {
		log.Fatalf("Не могу найти дату решения на странице проблемы %d", problemId)
	}
	answerDate, err := parseDate(answerDateText)
	if err != nil {
		log.Fatalf("Не могу открыть распарсить дату %s на странице проблемы %d", answerDateText, problemId)
	}

	log.Printf("%s \t=> %s сохраняем в базу\n", answerDateText, answerDate)
	models.SolutionsAdd(problemId, answerDate)
}

func getResponse(urlToParse string) *http.Response {

	// FIXME не надо каждый раз создавать
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	client := &http.Client{Transport: tr}
	resp, err := client.Get(urlToParse)
	if err != nil {
		log.Fatalf("Не могу открыть URL: %s. %s", urlToParse, err)
	}

	return resp
}
