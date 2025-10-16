//go:build linux
// +build linux

package main

import (
	"log"
	"os/exec"
)

// notify envia uma notificação de desktop no Linux usando notify-send.
func notify(message string) {
	// Chamamos o binário do sistema 'notify-send'
	cmd := exec.Command("notify-send", "-i", "pomodoro-icon", "Pomodoro CLI", message)

	// Tentativa de execução. Falha silenciosamente se o comando não puder ser executado.
	if err := cmd.Run(); err != nil {
		log.Printf("Failed to send Linux notification (is notify-send installed?): %v", err)
	}
}
