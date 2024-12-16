package main

import (
	"encoding/csv"
	"flag"
	"fmt"
	"os"
	"time"
)

type Problem struct {
	question string
	answer   string
}

func parseProblem(lines [][]string) []Problem {
	r := make([]Problem, len(lines))
	for i, line := range lines {
		r[i].question = line[0]
		r[i].answer = line[1]
	}
	return r
}

func getProblems(filename string) ([]Problem, error) {
	if fObj, err := os.Open(filename); err == nil {
		csvR := csv.NewReader(fObj)
		if cLines, err := csvR.ReadAll(); err == nil {
			return parseProblem(cLines), nil
		} else {
			return nil, err
		}
	} else {
		return nil, err
	}
}

func main() {
	fName := flag.String("file", "./static/problems.csv", "File to read")

	timer := flag.Int("t", 15, "Time in seconds")

	flag.Parse()

	problems, err := getProblems(*fName)

	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	correctAnswers := 0

	t := time.NewTimer(time.Duration(*timer) * time.Second)
	ansC := make(chan string)

problemLoop:
	for i, problem := range problems {
		var answer string
		fmt.Println("Problem: ", problem.question)

		go func() {
			fmt.Scanf("%s", &answer)
			ansC <- answer
		}()

		select {
		case <-t.C:
			fmt.Println("OK")
			break problemLoop
		case answer := <-ansC:
			if answer == problem.answer {
				correctAnswers++
			}
			if i == len(problems)-1 {
				close(ansC)
			}
		}
	}

	fmt.Printf("Your result is %d out of %d\n ", correctAnswers, len(problems))
	fmt.Println("Press enter to exit")
}

//TIP See GoLand help at <a href="https://www.jetbrains.com/help/go/">jetbrains.com/help/go/</a>.
// Also, you can try interactive lessons for GoLand by selecting 'Help | Learn IDE Features' from the main menu.
