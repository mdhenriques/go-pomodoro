package main

import (
	"bufio"
	"flag"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
	"time"

	"golang.org/x/term"
)

const (
	StateRunning TimerState = iota
	StatePaused
	StateStopped
)

const (
	TimerCompleted = "COMPLETED"
	TimerQuit      = "QUIT"
	CommandPause   = "pause"
	CommandResume  = "resume"
	CommandQuit    = "quit"
)

type TimerState int

// SessionStats tracks timing data for the pomodoro session
type SessionStats struct {
	StartTime       time.Time
	EndTime         time.Time
	TotalCycles     int
	CompletedCycles int
	TotalStudyTime  time.Duration
	TotalRestTime   time.Duration
}

func main() {
	fmt.Println("--- GO Pomodoro CLI ---")

	studyTimeFlag := flag.Int("study", 0, "study duration in minutes")
	restTimeFlag := flag.Int("rest", 0, "rest duration in minutes")
	cycleNumberFlag := flag.Int("cycles", 0, "number of pomodoro cycles")
	testUnitSecondsFlag := flag.Bool("test-unit-seconds", false, "If true, treats cycles durations as SECONDS instead of minutes for quick testing.")

	flag.Parse()

	studyDuration := getDurationFromFlagOrPrompt(*studyTimeFlag, "Enter study duration (minutes): ")
	restDuration := getDurationFromFlagOrPrompt(*restTimeFlag, "Enter rest duration (minutes): ")
	numCycles := getDurationFromFlagOrPrompt(*cycleNumberFlag, "Enter number of pomodoro cycles: ")

	timeUnit := time.Minute
	unitForPrompt := "minutes"
	if *testUnitSecondsFlag {
		timeUnit = time.Second
		unitForPrompt = "seconds"
	}

	studyDurationFinal := time.Duration(studyDuration) * timeUnit
	restDurationFinal := time.Duration(restDuration) * timeUnit

	fmt.Println("\nConfiguration:")
	fmt.Printf("Study Cycle:  %d %s\n", studyDuration, unitForPrompt)
	fmt.Printf("Rest Cycle:   %d %s\n", restDuration, unitForPrompt)
	fmt.Printf("Total Cycles: %d\n", numCycles)

	sessionStats := &SessionStats{
		StartTime:      time.Now(),
		TotalCycles:    numCycles,
		TotalStudyTime: studyDurationFinal * time.Duration(numCycles),
		TotalRestTime:  restDurationFinal * time.Duration(numCycles-1),
	}

	startPomodoro(studyDurationFinal, restDurationFinal, numCycles, sessionStats)

	sessionStats.EndTime = time.Now()
	displaySessionSummary(sessionStats)
}

func getDurationFromFlagOrPrompt(flagValue int, promptText string) int {
	if flagValue > 0 {
		return flagValue
	}

	if flagValue < 0 {
		fmt.Printf("Invalid flag value: %d. Value must be positive.\n", flagValue)
		fmt.Println("Please enter a valid value: ")
	}

	return promptForDuration(promptText)
}

func startPomodoro(studyDuration time.Duration, restDuration time.Duration, numCycles int, stats *SessionStats) {
	totalDuration := (studyDuration + restDuration) * time.Duration(numCycles)

	fmt.Printf("\nStarting Pomodoro session for %d cycles (Est. Total Time: %v)\n", numCycles, totalDuration)

	for i := 1; i <= numCycles; i++ {
		fmt.Printf("\n--- CYCLE %d of %d ---\n", i, numCycles)

		studyResult := runTimer("STUDY", studyDuration)
		if studyResult == TimerQuit {
			return
		}
		notify("Study Complete! Coffee Break")
		stats.CompletedCycles = i

		if i < numCycles {
			restResult := runTimer("REST", restDuration)
			if restResult == TimerQuit {
				return
			}
			notify("Rest ended. Keep focusing!!")
		}
	}

	notify("Session Complete!")
	fmt.Println("\nPomodoro Session Complete!")
}

func runTimer(cycleType string, duration time.Duration) string {
	fmt.Printf("\n⏰ Starting %s cycle for %v...\n", cycleType, duration)
	fmt.Println("Controls: [p] pause, [r] resume, [q] quit")

	ticker := time.NewTicker(1 * time.Second)
	defer ticker.Stop()

	tickCh := make(chan time.Time, 1)
	go func() {
		defer close(tickCh)
		for t := range ticker.C {
			tickCh <- t
		}
	}()

	remaining := duration
	state := StateRunning

	timer := time.NewTimer(remaining)
	defer func() {
		if timer != nil {
			timer.Stop()
		}
	}()

	inputCh := make(chan string)
	doneCh := make(chan struct{})
	go listenForInput(inputCh, doneCh)
	defer close(doneCh)

	cleanupSpaces := strings.Repeat(" ", 80)

	// Force initial display
	tickCh <- time.Now()

	for {
		select {
		case <-timer.C:
			fmt.Printf("\rTime remaining: 00:00%s\n", strings.Repeat(" ", 20))
			fmt.Printf("%s cycle complete!\n", cycleType)
			return cycleType

		case <-tickCh:
			if state == StateRunning {
				remaining -= 1 * time.Second

				minutes := int(remaining.Minutes())
				seconds := int(remaining.Seconds()) % 60

				timeStr := fmt.Sprintf("%02d:%02d", minutes, seconds)
				fmt.Printf("\rTime remaining: %s%s", timeStr, strings.Repeat(" ", 25))

				if remaining <= 0 {
					return cycleType
				}
			}

		case command := <-inputCh:
			switch command {
			case CommandPause:
				if state == StateRunning {
					state = StatePaused

					// Drain timer channel if it already fired
					if !timer.Stop() {
						select {
						case <-timer.C:
						default:
						}
					}

					fmt.Printf("\r⏸️  PAUSED - Remaining: %s [r]esume%s",
						formatDuration(remaining), cleanupSpaces)
				}

			case CommandResume:
				if state == StatePaused {
					state = StateRunning
					timer.Reset(remaining)

					fmt.Printf("\r▶️  RESUMED - Remaining: %s [p]ause%s",
						formatDuration(remaining), cleanupSpaces)

					tickCh <- time.Now()
				}

			case CommandQuit:
				fmt.Printf("\r%s\r", strings.Repeat(" ", 30))
				fmt.Println("\nTimer cancelled by user!")
				return TimerQuit
			}
		}
	}
}

func displaySessionSummary(stats *SessionStats) {
	fmt.Println()

	boxWidth := 38

	fmt.Printf("╔%s╗\n", strings.Repeat("═", boxWidth-2))

	title := "SESSION SUMMARY"
	titlePadding := (boxWidth - 2 - len(title)) / 2
	fmt.Printf("║%s%s%s║\n",
		strings.Repeat(" ", titlePadding),
		title,
		strings.Repeat(" ", boxWidth-2-len(title)-titlePadding))

	fmt.Printf("╠%s╣\n", strings.Repeat("═", boxWidth-2))

	fmt.Printf("║ Started:   %s║\n", formatLine(stats.StartTime.Format("15:04:05"), boxWidth-13))
	fmt.Printf("║ Ended:     %s║\n", formatLine(stats.EndTime.Format("15:04:05"), boxWidth-13))
	fmt.Printf("║ Cycles:    %s║\n", formatLine(fmt.Sprintf("%d/%d completed", stats.CompletedCycles, stats.TotalCycles), boxWidth-13))
	fmt.Printf("║ Study:     %s║\n", formatLine(formatTotalTime(stats.TotalStudyTime), boxWidth-13))
	fmt.Printf("║ Rest:      %s║\n", formatLine(formatTotalTime(stats.TotalRestTime), boxWidth-13))

	fmt.Printf("╚%s╝\n", strings.Repeat("═", boxWidth-2))
}

func listenForInput(inputCh chan<- string, doneCh <-chan struct{}) {
	fd := int(os.Stdin.Fd())
	oldState, err := term.MakeRaw(fd)
	if err != nil {
		log.Printf("Failed to set terminal to raw mode, falling back to buffered input: %v", err)
		listenForPause(inputCh, doneCh)
		return
	}
	defer term.Restore(fd, oldState)

	for {
		select {
		case <-doneCh:
			return
		default:
		}

		buf := make([]byte, 1)

		if _, err := os.Stdin.Read(buf); err != nil {
			continue
		}

		char := rune(buf[0])

		switch char {
		case 'p', 'P':
			inputCh <- CommandPause
		case 'r', 'R':
			inputCh <- CommandResume
		case 'q', 'Q':
			inputCh <- CommandQuit
		case 3: // Ctrl+C in raw mode
			inputCh <- CommandQuit
		}
	}
}

// listenForPause is fallback when raw terminal mode isn't available
func listenForPause(inputCh chan<- string, doneCh <-chan struct{}) {
	reader := bufio.NewReader(os.Stdin)
	for {
		select {
		case <-doneCh:
			return
		default:
		}
		inputText, err := reader.ReadString('\n')
		if err != nil {
			continue
		}

		inputChar := strings.TrimSpace(inputText)
		if len(inputChar) == 0 {
			continue
		}
		char := rune(inputChar[0])

		switch char {
		case 'p', 'P':
			inputCh <- CommandPause
		case 'r', 'R':
			inputCh <- CommandResume
		case 'q', 'Q':
			inputCh <- CommandQuit
		}
	}
}

func formatLine(text string, width int) string {
	if len(text) > width {
		return text[:width]
	}
	return text + strings.Repeat(" ", width-len(text))
}

func formatTotalTime(d time.Duration) string {
	minutes := int(d.Minutes())

	if minutes == 0 && d > 0 {
		return fmt.Sprintf("%d seconds", int(d.Seconds()))
	}
	return fmt.Sprintf("%d minutes", minutes)
}

func formatDuration(d time.Duration) string {
	minutes := int(d.Minutes())
	seconds := int(d.Seconds()) % 60
	return fmt.Sprintf("%02d:%02d", minutes, seconds)
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