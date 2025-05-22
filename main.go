package main

import (
	"fmt"
	"os/exec"
	"time"

	"github.com/getlantern/systray"
)

// Global variables
var (
	isSleepPrevented bool
	timerDuration    time.Duration
	timerEndTime     time.Time
	timerActive      bool
)

func main() {
	// Start the systray
	systray.Run(onReady, onExit)
}

func onReady() {
	// Set up the menu bar icon
	systray.SetIcon(getIcon("default"))
	systray.SetTitle("Jolt")
	systray.SetTooltip("Jolt - Prevent your Mac from sleeping")

	// Create menu items
	mToggle := systray.AddMenuItem("Enable Jolt", "Toggle sleep prevention")
	systray.AddSeparator()

	// Timer submenu
	mTimerSubmenu := systray.AddMenuItem("Set Timer", "Set a timer for sleep prevention")
	m15min := mTimerSubmenu.AddSubMenuItem("15 minutes", "Prevent sleep for 15 minutes")
	m30min := mTimerSubmenu.AddSubMenuItem("30 minutes", "Prevent sleep for 30 minutes")
	m1hour := mTimerSubmenu.AddSubMenuItem("1 hour", "Prevent sleep for 1 hour")
	m2hour := mTimerSubmenu.AddSubMenuItem("2 hours", "Prevent sleep for 2 hours")
	mCustom := mTimerSubmenu.AddSubMenuItem("Custom...", "Set a custom timer duration")
	mCancelTimer := mTimerSubmenu.AddSubMenuItem("Cancel Timer", "Cancel the active timer")
	mCancelTimer.Disable()

	systray.AddSeparator()
	mQuit := systray.AddMenuItem("Quit", "Quit Jolt")

	// Handle menu item clicks
	go func() {
		for {
			select {
			case <-mToggle.ClickedCh:
				toggleSleepPrevention()
				updateMenuItems(mToggle, mCancelTimer)
			case <-m15min.ClickedCh:
				setTimer(15 * time.Minute)
				updateMenuItems(mToggle, mCancelTimer)
			case <-m30min.ClickedCh:
				setTimer(30 * time.Minute)
				updateMenuItems(mToggle, mCancelTimer)
			case <-m1hour.ClickedCh:
				setTimer(time.Hour)
				updateMenuItems(mToggle, mCancelTimer)
			case <-m2hour.ClickedCh:
				setTimer(2 * time.Hour)
				updateMenuItems(mToggle, mCancelTimer)
			case <-mCustom.ClickedCh:
				// In a real implementation, this would show a dialog
				// For now, we'll just set a 45-minute timer as an example
				setTimer(45 * time.Minute)
				updateMenuItems(mToggle, mCancelTimer)
			case <-mCancelTimer.ClickedCh:
				cancelTimer()
				updateMenuItems(mToggle, mCancelTimer)
			case <-mQuit.ClickedCh:
				systray.Quit()
				return
			}
		}
	}()

	// Timer monitoring goroutine
	go func() {
		ticker := time.NewTicker(1 * time.Second)
		defer ticker.Stop()

		for range ticker.C {
			if timerActive && time.Now().After(timerEndTime) {
				// Timer has expired
				disableSleepPrevention()
				timerActive = false
				updateMenuItems(mToggle, mCancelTimer)

				// Show notification
				showNotification("Jolt Timer Expired", "Sleep prevention has been disabled.")
			}

			if timerActive {
				// Update the remaining time in the menu
				remaining := timerEndTime.Sub(time.Now())
				hours := int(remaining.Hours())
				minutes := int(remaining.Minutes()) % 60
				seconds := int(remaining.Seconds()) % 60

				prefix := "Jolt"
				if isSleepPrevented {
					prefix = "⚡️ Jolt"
				}

				if hours > 0 {
					systray.SetTitle(fmt.Sprintf("%s (%d:%02d:%02d)", prefix, hours, minutes, seconds))
				} else {
					systray.SetTitle(fmt.Sprintf("%s (%d:%02d)", prefix, minutes, seconds))
				}
			}
		}
	}()
}

func onExit() {
	// Make sure to disable sleep prevention when exiting
	if isSleepPrevented {
		disableSleepPrevention()
	}
}

func toggleSleepPrevention() {
	if isSleepPrevented {
		disableSleepPrevention()
	} else {
		enableSleepPrevention()
	}
}

func enableSleepPrevention() {
	// Use caffeinate command to prevent sleep
	cmd := exec.Command("caffeinate", "-d", "-i", "-u", "-s")
	err := cmd.Start()
	if err != nil {
		fmt.Println("Error enabling sleep prevention:", err)
		return
	}

	isSleepPrevented = true
	systray.SetTitle("⚡️ Jolt")
	systray.SetTooltip("Jolt - Sleep Prevention Active")
}

func disableSleepPrevention() {
	// Kill all caffeinate processes
	if err := exec.Command("killall", "caffeinate").Run(); err != nil {
		fmt.Println("Error disabling sleep prevention:", err)
	}

	isSleepPrevented = false
	timerActive = false
	systray.SetTitle("Jolt")
	systray.SetTooltip("Jolt - Prevent your Mac from sleeping")
}

func setTimer(duration time.Duration) {
	timerDuration = duration
	timerEndTime = time.Now().Add(duration)
	timerActive = true

	if !isSleepPrevented {
		enableSleepPrevention()
	}
}

func cancelTimer() {
	timerActive = false
	if isSleepPrevented {
		systray.SetTitle("⚡️ Jolt")
	} else {
		systray.SetTitle("Jolt")
	}
}

func updateMenuItems(mToggle, mCancelTimer *systray.MenuItem) {
	if isSleepPrevented {
		mToggle.SetTitle("Disable Jolt")
	} else {
		mToggle.SetTitle("Enable Jolt")
	}

	if timerActive {
		mCancelTimer.Enable()
	} else {
		mCancelTimer.Disable()
	}
}

func showNotification(title, message string) {
	// Use AppleScript to show a notification
	script := fmt.Sprintf(`display notification "%s" with title "%s"`, message, title)
	if err := exec.Command("osascript", "-e", script).Run(); err != nil {
		fmt.Println("Failed to show notification:", err)
	}
}

// getIcon returns a simple placeholder icon
// We're using emoji characters in the title instead of custom icons
func getIcon(state string) []byte {
	// Return a simple placeholder icon (1x1 transparent pixel)
	return []byte{0}
}
