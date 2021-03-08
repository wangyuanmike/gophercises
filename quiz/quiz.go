package main

import (
	"encoding/csv"
	"fmt"
	"os"
	"strings"
)

type problem struct {
	question string
	answer   string
}

func exit(errMsg string) {
	fmt.Println(errMsg)
	os.Exit(1)
}

func parseProblem(problemFile string) []problem {
	f, err := os.Open(problemFile)
	defer f.Close()
	if err != nil {
		exit("Cannot open problem file...")
	}

	csvReader := csv.NewReader(f)
	records, err := csvReader.ReadAll()
	if err != nil {
		exit("Cannot read problem file...")
	}

	var problems = make([]problem, len(records))
	for i := range records {
		problems[i].question = records[i][0]
		problems[i].answer = strings.TrimSpace(records[i][1])
	}
	return problems
}

func executeQuiz(problemFile string) {
	problems := parseProblem(problemFile)

	correctCount := 0
	var answer string
	fmt.Println("Quiz starts...")
	for i := range problems {
		fmt.Println(problems[i].question)
		fmt.Scanln(&answer)
		if answer == problems[i].answer {
			correctCount++
		}
	}

	fmt.Printf("You scored %d out of %d questions", correctCount, len(problems))
}

func main() {
	executeQuiz("problems.csv")
}
