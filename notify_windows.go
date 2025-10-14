//go:build windows
// +build windows

package main

import (
	"fmt"
	"log"
	"os/exec"
	"strings"
)

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