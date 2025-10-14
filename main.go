package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
	"time"
	"os/exec"
)

const (
	toastTemplateTypeText02 = 2
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

	notify("Session Complete!")
	fmt.Println("\n Pomodoro Session Complete!")
}

func runTimer(cycleType string, duration time.Duration) {
	fmt.Printf("\n⏰ Starting %s cycle for %v...\n", cycleType, duration)
    
    // 1. Setup the Ticker and Timer
	ticker := time.NewTicker(1 * time.Second)
	defer ticker.Stop()
	 // Ticks every 1 second
	timer := time.NewTimer(duration)
	defer timer.Stop()         // Timer to signal when the total duration is up
	
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