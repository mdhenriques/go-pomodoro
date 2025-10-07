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
	sessionHours := promptForDuration("Enter total session duration (hours): ")

	fmt.Println("\nConfiguration:")
	fmt.Printf("Study Cycle: %d minutes\n", studyDuration)
	fmt.Printf("Rest Cycle:  %d minutes\n", restDuration)
	fmt.Printf("Session:     %d hours\n", sessionHours)

	totalSessionTime := time.Duration(sessionHours) * time.Hour

	startPomodoro(studyDuration, restDuration, totalSessionTime)
}

func startPomodoro(studyMins, restMins int, totalTime time.Duration) {
	startTime := time.Now()

	//Create Ticker durations
	studyDuration := time.Duration(studyMins) * time.Minute
	restDuration := time.Duration(restMins) * time.Minute

	fmt.Printf("\nStarting Pomodoro session for %v...\n", totalTime)

	//Main loop for the session
	for time.Since(startTime) < totalTime {
		runTimer("STUDY", studyDuration)

		if time.Since(startTime) >= totalTime{
			break
		}

		runTimer("REST", restDuration)
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

	fmt.Print(promptText)

	inputText, err := reader.ReadString('\n')
	if err != nil {
		log.Fatalf("Failed to read user input: %v", err)
	}

	inputText = strings.TrimSpace(inputText)

	duration, err := strconv.Atoi(inputText)
	if err != nil {
		log.Fatalf("Invalid input: '%s' is not a valid number. Error: %v", inputText, err)
	}

	if duration <= 0 {
		fmt.Println("Duration must be a positive number. Please try again")
	}

	return duration
}