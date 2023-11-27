// Quiz reads from a csv file where its lines are in the format question,answer
package main

import (
	"encoding/csv"
	"flag"
	"fmt"
	"os"
	"strings"
)

type problem struct {
	question, answer string
}

type quiz []problem

func main() {
	csv := flag.String("csv", "problems.csv", "Filename of the quiz csv file in the format\"question,answer\". Defaults to \"problems.csv\".")
	flag.Parse()
	q := openCSV(*csv)
	play(q)
}

// errorExit displays error and then exit
func errorExit(e error) {
	fmt.Println(e)
	os.Exit(1)
}

// openCSV file and returns a quiz
func openCSV(filename string) quiz {
	f, err := os.Open(filename)
	if err != nil {
		errorExit(err)
	}
	reader := csv.NewReader(f)
	records, err :=  reader.ReadAll()
	if err != nil {
		errorExit(err)
	}
	return parseLines(records)
}

// parseLines parse the csv records to a quiz
func parseLines(records [][]string) quiz {
	q := make(quiz, len(records))
	for i, v := range records {
		q[i].question = strings.Join(v[:len(v)-1], "") //it joins all the fields beside the last, allowing questions with a ","
		q[i].answer = v[len(v)-1]
	}
	return q
}

// play display question and receive answer of a quiz
func play(q quiz) {
	points := 0
	var input string

	for i, problem := range q {
		fmt.Printf("Question %d: %s\n", i+1, problem.question)
		_, err := fmt.Scanf("%s\n", &input)
		if err != nil {
			errorExit(err)
		}
		if input == problem.answer {
			points++
		} 
	}

	fmt.Printf("Game over! You scored %v of %v points.", points, len(q))
}
