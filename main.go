package main

import (
	"bufio"
	"flag"
	"fmt"
	"log"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"time"
)

const (
	toastTemplateTypeText02 = 2
)

//SessionStats holds all the data for the summary
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

	flag.Parse()

	//Validate and get study duration
	studyDuration := getDurationFromFlagOrPrompt(*studyTimeFlag, "Enter study duration (minutes): ")

	restDuration := getDurationFromFlagOrPrompt(*restTimeFlag, "Enter rest duration (minutes): ")

	numCycles := getDurationFromFlagOrPrompt(*cycleNumberFlag, "Enter number of pomodoro cycles: ")

	fmt.Println("\nConfiguration:")
	fmt.Printf("Study Cycle:  %d minutes\n", studyDuration)
	fmt.Printf("Rest Cycle:   %d minutes\n", restDuration)
	fmt.Printf("Total Cycles: %d\n", numCycles)

	// Session stats and start tracking
	sessionStats := &SessionStats{
		StartTime: time.Now(),
		TotalCycles: numCycles,
		TotalStudyTime: time.Duration(studyDuration) * time.Minute * time.Duration(numCycles),
		TotalRestTime: time.Duration(restDuration) * time.Minute * time.Duration(numCycles-1), //Last cycle doesnt have rest
	}

	startPomodoro(studyDuration, restDuration, numCycles, sessionStats)

	//Update end time and show summary
	sessionStats.EndTime = time.Now()
	displaySessionSummary(sessionStats)
}

// getDurationFromFlagOrPrompt handles the logic for getting a value either from flag or prompt
func getDurationFromFlagOrPrompt(flagValue int, promptText string) int {
	//If flag was provided and has a valid value (> 0), use it
	if flagValue > 0 {
		return flagValue
	}

	// If flag was provided but value is invalid (<= 0), show error and prompt
	if flagValue < 0 {
		fmt.Printf("Invalid flag value: %d. Value must be positive.\n", flagValue)
		fmt.Println("Please enter a valid value: ")
	}

	return promptForDuration(promptText)
}

func displaySessionSummary(stats *SessionStats) {
	fmt.Println() // Add some space before the summary
	
	boxWidth := 38 // Slightly narrower for better fit
	
	// Top border
	fmt.Printf("╔%s╗\n", strings.Repeat("═", boxWidth-2))
	
	// Title line - centered
	title := "SESSION SUMMARY"
	titlePadding := (boxWidth - 2 - len(title)) / 2
	fmt.Printf("║%s%s%s║\n", 
		strings.Repeat(" ", titlePadding),
		title,
		strings.Repeat(" ", boxWidth-2-len(title)-titlePadding))
	
	// Separator
	fmt.Printf("╠%s╣\n", strings.Repeat("═", boxWidth-2))
	
	// Content lines - properly aligned
	fmt.Printf("║ Started:   %s║\n", formatLine(stats.StartTime.Format("15:04:05"), boxWidth-13))
	fmt.Printf("║ Ended:     %s║\n", formatLine(stats.EndTime.Format("15:04:05"), boxWidth-13))
	fmt.Printf("║ Cycles:    %s║\n", formatLine(fmt.Sprintf("%d/%d completed", stats.CompletedCycles, stats.TotalCycles), boxWidth-13))
	fmt.Printf("║ Study:     %s║\n", formatLine(fmt.Sprintf("%d minutes", int(stats.TotalStudyTime.Minutes())), boxWidth-13))
	fmt.Printf("║ Rest:      %s║\n", formatLine(fmt.Sprintf("%d minutes", int(stats.TotalRestTime.Minutes())), boxWidth-13))
	
	// Bottom border
	fmt.Printf("╚%s╝\n", strings.Repeat("═", boxWidth-2))
}

// Simplified formatting function
func formatLine(text string, width int) string {
	if len(text) > width {
		return text[:width]
	}
	return text + strings.Repeat(" ", width-len(text))
}


func notify(title string) {
	// 1. Simple escaping: Replace single quotes with double single quotes ('' is a literal ' in PS)
	safeTitle := strings.ReplaceAll(title, "'", "''")

	// 2. Build the PowerShell command string.
	// We use a simplified template ('ToastText02') and avoid multi-line XML definition
	// to keep the command string short and reliable for os/exec.
	script := fmt.Sprintf(`
        [Windows.UI.Notifications.ToastNotificationManager, Windows.UI.Notifications, ContentType = WindowsRuntime] | Out-Null;
        $template = [Windows.UI.Notifications.ToastNotificationManager]::GetTemplateContent(%d);
        $template.GetElementsByTagName("text")[0].AppendChild($template.CreateTextNode('%s')) | Out-Null;
        $template.GetElementsByTagName("text")[1].AppendChild($template.CreateTextNode('')) | Out-Null;
        $toast = [Windows.UI.Notifications.ToastNotification]::new($template);
        [Windows.UI.Notifications.ToastNotificationManager]::CreateToastNotifier("GO Pomodoro CLI").Show($toast);
    `, toastTemplateTypeText02, safeTitle)

	cmd := exec.Command("powershell", "-ExecutionPolicy", "Bypass", "-Command", script)

	if err := cmd.Run(); err != nil {
		log.Printf("Failed to show Windows notification: %v", err)
		// Fallback to terminal bell
		fmt.Print("\a")
	}
}

func startPomodoro(studyMins, restMins int, numCycles int, stats *SessionStats) {
	studyDuration := time.Duration(studyMins) * time.Minute
	restDuration := time.Duration(restMins) * time.Minute

	totalDuration := (studyDuration + restDuration) * time.Duration(numCycles)

	fmt.Printf("\nStarting Pomodoro session for %d cycles (Est. Total Time: %v)\n", numCycles, totalDuration)

	for i := 1; i <= numCycles; i++ {
		fmt.Printf("\n--- CYCLE %d of %d ---\n", i, numCycles)
		
		runTimer("STUDY", studyDuration)
		stats.CompletedCycles = i // Update completed cycles after each study session

		if i < numCycles {
			runTimer("REST", restDuration)
		}
	}

	notify("Session Complete!")
	fmt.Println("\nPomodoro Session Complete!")
}

func runTimer(cycleType string, duration time.Duration) {
	fmt.Printf("\n⏰ Starting %s cycle for %v...\n", cycleType, duration)

	// 1. Setup the Ticker and Timer
	ticker := time.NewTicker(1 * time.Second)
	defer ticker.Stop()
	// Ticks every 1 second
	timer := time.NewTimer(duration)
	defer timer.Stop() // Timer to signal when the total duration is up

	remaining := duration // Start remaining time at the full duration

	// 2. Main Countdown Loop
	for {
		select {
		case <-timer.C:
			// Timer expired - cycle complete
			fmt.Printf("\rTime remaining: 00:00%s\n", strings.Repeat(" ", 5))

			// Send notification
			if cycleType == "STUDY" {
				notify("Study Complete! It's REST time.")
			} else {
				notify("Rest Complete! Time to go back to WORK!")
			}

			fmt.Printf("%s cycle complete!\n", cycleType)
			return

		case <-ticker.C:
			// Update countdown display
			remaining -= 1 * time.Second

			minutes := int(remaining.Minutes())
			seconds := int(remaining.Seconds()) % 60

			timeStr := fmt.Sprintf("%02d:%02d", minutes, seconds)
			fmt.Printf("\rTime remaining: %s", timeStr)
		}
	}
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
