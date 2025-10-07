package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
	"time"
)

func main() {
	fmt.Println("--- GO Pomodoro CLI ---")

	// Get configuration from user
	studyDuration := promptForDuration("Enter study duration (minutes): ")
	restDuration := promptForDuration("Enter rest duration (minutes): ")
	numCycles := promptForDuration("Enter total number of cycles: ")

	fmt.Println("\nConfiguration:")
	fmt.Printf("Study Cycle:  %d minutes\n", studyDuration)
	fmt.Printf("Rest Cycle:   %d minutes\n", restDuration)
	fmt.Printf("Total Cycles: %d\n", numCycles)

	startPomodoro(studyDuration, restDuration, numCycles)
}

func startPomodoro(studyMins, restMins int, numCycles int) {
	studyDuration := time.Duration(studyMins) * time.Minute
	restDuration := time.Duration(restMins) * time.Minute

	totalDuration := (studyDuration + restDuration) * time.Duration(numCycles)

	fmt.Printf("\nStarting Pomodoro session for %d cycles (Est. Total Time: %v)\n", numCycles, totalDuration)

	for i := 1; i <= numCycles; i++ {
		fmt.Printf("\n--- CYCLE %d of %d ---\n", i, numCycles)
		
		runTimer("STUDY", studyDuration)

		if i < numCycles {
			runTimer("REST", restDuration)
		}
	}

	fmt.Println("\n Pomodoro Session Complete!")
}

func runTimer(cycleType string, duration time.Duration) {
	fmt.Printf("\n Starting %s cycle for %v...\n", cycleType, duration)

	timer := time.NewTimer(duration)

	<-timer.C
	fmt.Printf("%s cycle complete!\n", cycleType)
}

func promptForDuration(promptText string) int {
	reader := bufio.NewReader(os.Stdin)

	for {
		fmt.Print(promptText)

		inputText, err := reader.ReadString('\n')
		if err != nil {
			log.Fatalf("Failed to read user input: %v", err)
		}

		inputText = strings.TrimSpace(inputText)

		duration, err := strconv.Atoi(inputText)
		if err != nil {
			fmt.Printf("Invalid input: '%s' is not a whole number. Please enter a valid integer.\n", inputText)
			continue
		}

		if duration <= 0 {
			fmt.Println("Duration or cycle count must be a positive number. Please try again.")
			continue
		}

		return duration
	}

}