package main

import (
	"fmt"
	"github.com/goodsign/monday"
	"log"
	"regexp"
	"strings"
	"time"
)

var loc *time.Location

func init() {
	var err error
	loc, err = time.LoadLocation("Asia/Novosibirsk")
	if err != nil {
		log.Fatal("Не могу загрузить location ", err)
	}
}

// parseDate парсит дату в формате тактака
func parseDate(taktakDate string) (time.Time, error) {
	taktakDate = strings.Trim(taktakDate, " ")
	// дата без года - добавляем год
	re := regexp.MustCompile(`(\d+ [а-я]+), (\d+)`)
	taktakDate = re.ReplaceAllString(taktakDate, fmt.Sprintf("$1 %d, $2", time.Now().Year()))

	// добавляем текущую дату
	re = regexp.MustCompile(`(сегодня), (\d+:\d+)`)
	taktakDate = re.ReplaceAllString(taktakDate, time.Now().Format("2 January 2006")+", $2")

	// добавляем вчерашнюю дату
	re = regexp.MustCompile(`(вчера), (\d+:\d+)`)
	yesterday := time.Now().AddDate(0, 0, -1)
	taktakDate = re.ReplaceAllString(taktakDate, yesterday.Format("2 January 2006")+", $2")

	if taktakDate == "1 час назад" {
		return time.Now().Add(-1 * time.Hour), nil
	}

	if taktakDate == "2 часа назад" {
		return time.Now().Add(-1 * time.Hour), nil
	}

	parsedDate, err := monday.ParseInLocation(
		"2 January 2006, 15:04",
		taktakDate,
		loc,
		monday.LocaleRuRU,
	)
	if err != nil {
		return time.Time{}, fmt.Errorf("Не могу распарсить дату %s. %s", taktakDate, err)
	}

	return parsedDate, nil
}
