package models

import (
	"fmt"
	"log"
	"time"
)

type Solution struct {
	ProblemId    int
	ProblemTitle string
	ProblemDate  time.Time
	SolutionDate time.Time
	ExpertId     int
	ExpertName   string
}

func solutionsInit() error {
	log.Println("Initing solutions table")

	_, err := db.Exec(`
		CREATE TABLE IF NOT EXISTS solutions (
			problem_id int NOT NULL PRIMARY KEY,
			problem_title text NOT NULL DEFAULT '',
			problem_date date NOT NULL DEFAULT '0000-00-00 00:00:00',
			solution_date date NOT NULL DEFAULT '0000-00-00 00:00:00',
			expert_id int NOT NULL DEFAULT 0,
			expert_name text NOT NULL DEFAULT ''
		)
	`)
	if err != nil {
		return err
	}

	_, err = db.Exec("CREATE INDEX IF NOT EXISTS solution_date ON solutions (solution_date)")
	if err != nil {
		return err
	}

	return nil
}

// SolutionsAdd добавляет информацию о решении проблемы
func SolutionsAdd(problemId int, solutionDate time.Time) error {
	stmt, err := db.Prepare("INSERT INTO solutions(problem_id, solution_date) VALUES (?, ?)")
	if err != nil {
		return err
	}
	defer stmt.Close()
	_, err = stmt.Exec(problemId, solutionDate)
	if err != nil {
		return fmt.Errorf("Can't add solution: %s", err)
	}

	return nil
}

// SolutionCheckExists проверяет существует ли в базе информация об указанной проблеме
func SolutionCheckExists(problemId int) (bool, error) {
	var count int
	err := db.QueryRow("SELECT count(*) FROM solutions WHERE problem_id=?", problemId).Scan(&count)
	if err != nil {
		return false, err
	}
	if count > 0 {
		return true, nil
	} else {
		return false, nil
	}
}

// SolutionsGetCountFrom возвращает число решений с указаного времени
func SolutionsGetCountFrom(from time.Time) (int, error) {
	var count int
	err := db.QueryRow("SELECT COUNT (*) FROM solutions WHERE solution_date > ?", from).Scan(&count)
	if err != nil {
		return 0, err
	}
	return count, nil
}
