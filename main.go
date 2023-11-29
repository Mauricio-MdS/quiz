// Quiz reads from a csv file where its lines are in the format question,answer
package main

import (
	"encoding/csv"
	"flag"
	"fmt"
	"math/rand"
	"os"
	"strings"
	"time"
)

type problem struct {
	question, answer string
}

type quiz []problem



func main() {
	csv := flag.String("csv", "problems.csv", "Filename of the quiz csv file in the format\"question,answer\". Defaults to \"problems.csv\".")
	random := flag.Bool("random", false, "If true, will randomize the question order. Defaults to false.")
	timeLimit := flag.Int("limit", 30, "Time in seconds to complete the quiz. Defaults to 30s.")
	flag.Parse()
	q := openCSV(*csv)
	if *random {
		q = shuffleQuiz(q)
	}
	play(q, *timeLimit)
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
func play(q quiz, timeLimit int) {
	points := 0
	answerChan := make(chan string)

	fmt.Printf("The quiz will start with a timer of %v seconds. Press enter to start.", timeLimit)
	fmt.Scanln()
	timer := time.NewTimer(time.Duration(time.Second * time.Duration(timeLimit)))

	loop:
	for i, problem := range q {
		fmt.Printf("Question %d: %s\n", i+1, problem.question)
		go playerInput(answerChan)
		select {
		case <- timer.C:
			break loop
		case ans := <-answerChan:
			if strings.TrimSpace(strings.ToLower(ans)) == strings.TrimSpace(strings.ToLower(problem.answer)) {
				points++
			} 
		}		
	}

	fmt.Printf("Game over! You scored %v of %v points.", points, len(q))
}

// playerInput gets input for the player and puts in the answerChan
func playerInput(answerChan chan string){
	var input string
	fmt.Scanf("%s\n", &input)
	answerChan <- input
}


//shuffleQuiz randomizes the order of the quiz
func shuffleQuiz(q quiz) quiz {
	for i := range q {
		r := rand.Intn(len(q))
		q[i], q[r] = q[r], q[i]
	}
	return q
}
