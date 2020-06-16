package main

import (
	"context"
	"encoding/csv"
	"fmt"
	"log"
	"math/rand"
	"os"
	"strings"
	"time"
)

func main() {
	var filePath = "problems.csv"
	fmt.Printf("Enter the file path to the quiz: (%s)", filePath)
	fmt.Scanln(&filePath)
	// os package read file
	file, err := os.Open(filePath)
	defer file.Close()

	if err != nil {
		log.Fatal(err)
	}

	// csv package create csv reader
	reader := csv.NewReader(file)
	records, err := reader.ReadAll()
	if err != nil {
		log.Fatalln(err)
	}

	shuffle(records)

	total, correct := doQuizWithTimeout(context.Background(), records)

	fmt.Printf("Correct: %d/%d", correct, total)
}

func shuffle(records [][]string) {
	rand.Seed(time.Now().UnixNano())
	rand.Shuffle(len(records), func(i, j int) {
		records[i], records[j] = records[j], records[i]
	})
}

func doQuizWithTimeout(ctx context.Context, records [][]string) (total, correct int) {
	// wait for user interaction
	fmt.Print("Press the Enter key to start the quiz")
	fmt.Scanln()

	// set up timer - create context
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	correctChan := make(chan bool, len(records))
	doneChan := make(chan bool)

	go func(records [][]string) {
		for _, record := range records {
			fmt.Printf("%s = ", record[0])
			var answer string
			fmt.Scanln(&answer)
			if strings.EqualFold(record[1], strings.Trim(answer, " ")) {
				correctChan <- true
			}
		}
		doneChan <- true
	}(records)

	select {
	case <-ctx.Done():
		fmt.Println("\nTime is up.")
	case <-doneChan:
		cancel()
	}
	return len(records), len(correctChan)
}
