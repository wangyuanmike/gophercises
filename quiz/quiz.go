package main

import (
	"encoding/csv"
	"flag"
	"fmt"
	"os"
	"strings"
	"time"
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
		exit("Failed to open problem file...")
	}

	csvReader := csv.NewReader(f)
	records, err := csvReader.ReadAll()
	if err != nil {
		exit("Failed to read problem file...")
	}

	var problems = make([]problem, len(records))
	for i, record := range records {
		problems[i] = problem{
			question: record[0],
			answer:   strings.TrimSpace(record[1]),
		}
	}
	return problems
}

func executeQuiz(problemFile string, timeout int) {
	problems := parseProblem(problemFile)

	correctCount := 0
	fmt.Println("Quiz starts...")
	timer := time.NewTimer(time.Duration(timeout) * time.Second)
	for i, problem := range problems {
		fmt.Printf("Question #%d: %s\n", i, problem.question)
		answerCh := make(chan string)
		go func() {
			var answer string
			fmt.Scanln(&answer)
			answerCh <- answer
		}()
		select {
		case <-timer.C:
			fmt.Printf("You scored %d out of %d questions\n", correctCount, len(problems))
			exit("Timeout, quiz is finished...")
		case answer := <-answerCh:
			if answer == problem.answer {
				correctCount++
			}
		}
	}

	fmt.Printf("You scored %d out of %d questions\n", correctCount, len(problems))
}

func main() {
	file := flag.String("file", "problems.csv", "path of problem file")
	timeout := flag.Int("timeout", 30, "timer duration")
	flag.Parse()

	executeQuiz(*file, *timeout)
}
