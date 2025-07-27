//go:build windows
// +build windows

// Package service provides utilities for running the application as a Windows service,
// including handling Service Control Manager (SCM) commands such as start and stop.
package service

import (
	"github.com/rs/zerolog/log"

	"golang.org/x/sys/windows/svc"
)

// SCM command handler (start, stop)
type windowsServiceHandler struct {
	runApp func()
}

// Execute processes SCM requests such as start and stop, and runs the main application logic.
func (h *windowsServiceHandler) Execute(_ []string, r <-chan svc.ChangeRequest, s chan<- svc.Status) (bool, uint32) {
	s <- svc.Status{State: svc.StartPending}
	go h.runApp()
	s <- svc.Status{State: svc.Running, Accepts: svc.AcceptStop}

	// Handle SCM requests (e.g., stop, shutdown)
	for c := range r {
		switch c.Cmd {
		case svc.Stop, svc.Shutdown:
			s <- svc.Status{State: svc.StopPending}
			return false, 0
		}
	}

	return false, 0
}

// IsServiceMode checks if the current process is running as a Windows service.
func IsServiceMode() bool {
	isService, _ := svc.IsWindowsService()
	return isService
}

// RunAsService runs the given application function as a Windows service
// and integrates with the Windows Service Control Manager (SCM).
func RunAsService(runApp func()) {
	err := svc.Run("discord-invite", &windowsServiceHandler{runApp: runApp})
	if err != nil {
		log.Fatal().Msgf("Service fail with error: %v", err)
	}
}
