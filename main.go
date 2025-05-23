package main

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"time"

	"github.com/getlantern/systray"
)

// Global variables
var (
	isSleepPrevented bool
	timerDuration    time.Duration
	timerEndTime     time.Time
	timerActive      bool
	autoStartEnabled bool
)

func main() {
	// Start the systray
	systray.Run(onReady, onExit)
}

func onReady() {
	// Set up the menu bar icon
	systray.SetIcon(getIcon("default"))
	systray.SetTitle("")
	systray.SetTooltip("Jolt - Prevent your Mac from sleeping")

	// Check if auto-start is enabled
	autoStartEnabled = isAutoStartEnabled()

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

	// Auto-start menu item
	mAutoStart := systray.AddMenuItem("Start at Login", "Toggle start at login")

	systray.AddSeparator()
	mQuit := systray.AddMenuItem("Quit", "Quit Jolt")

	// Handle menu item clicks
	go func() {
		for {
			select {
			case <-mToggle.ClickedCh:
				toggleSleepPrevention()
				updateMenuItems(mToggle, mCancelTimer, mAutoStart)
			case <-m15min.ClickedCh:
				setTimer(15 * time.Minute)
				updateMenuItems(mToggle, mCancelTimer, mAutoStart)
			case <-m30min.ClickedCh:
				setTimer(30 * time.Minute)
				updateMenuItems(mToggle, mCancelTimer, mAutoStart)
			case <-m1hour.ClickedCh:
				setTimer(time.Hour)
				updateMenuItems(mToggle, mCancelTimer, mAutoStart)
			case <-m2hour.ClickedCh:
				setTimer(2 * time.Hour)
				updateMenuItems(mToggle, mCancelTimer, mAutoStart)
			case <-mCustom.ClickedCh:
				// In a real implementation, this would show a dialog
				// For now, we'll just set a 45-minute timer as an example
				setTimer(45 * time.Minute)
				updateMenuItems(mToggle, mCancelTimer, mAutoStart)
			case <-mCancelTimer.ClickedCh:
				cancelTimer()
				updateMenuItems(mToggle, mCancelTimer, mAutoStart)
			case <-mAutoStart.ClickedCh:
				toggleAutoStart()
				updateMenuItems(mToggle, mCancelTimer, mAutoStart)
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
				updateMenuItems(mToggle, mCancelTimer, mAutoStart)

				// Show notification
				showNotification("Jolt Timer Expired", "Sleep prevention has been disabled.")
			}

			if timerActive {
				// Update the remaining time in the menu
				remaining := timerEndTime.Sub(time.Now())
				hours := int(remaining.Hours())
				minutes := int(remaining.Minutes()) % 60
				seconds := int(remaining.Seconds()) % 60

				if hours > 0 {
					systray.SetTitle(fmt.Sprintf("%d:%02d:%02d", hours, minutes, seconds))
				} else {
					systray.SetTitle(fmt.Sprintf("%d:%02d", minutes, seconds))
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
	systray.SetIcon(getIcon("active"))
	systray.SetTitle("")
	systray.SetTooltip("Jolt - Sleep Prevention Active")
}

func disableSleepPrevention() {
	// Kill all caffeinate processes
	if err := exec.Command("killall", "caffeinate").Run(); err != nil {
		fmt.Println("Error disabling sleep prevention:", err)
	}

	isSleepPrevented = false
	timerActive = false
	systray.SetIcon(getIcon("default"))
	systray.SetTitle("")
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
	systray.SetTitle("")
}

func updateMenuItems(mToggle, mCancelTimer, mAutoStart *systray.MenuItem) {
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

	if autoStartEnabled {
		mAutoStart.SetTitle("Start at Login")
		mAutoStart.Check()
	} else {
		mAutoStart.SetTitle("Start at Login")
		mAutoStart.Uncheck()
	}
}

// toggleAutoStart toggles the auto-start feature
func toggleAutoStart() {
	if autoStartEnabled {
		disableAutoStart()
	} else {
		enableAutoStart()
	}
}

func showNotification(title, message string) {
	// Use AppleScript to show a notification
	script := fmt.Sprintf(`display notification "%s" with title "%s"`, message, title)
	if err := exec.Command("osascript", "-e", script).Run(); err != nil {
		fmt.Println("Failed to show notification:", err)
	}
}

// getIcon loads and returns the appropriate icon based on the state
func getIcon(state string) []byte {
	var iconName string

	// Determine which icon to load based on the state
	if state == "active" {
		iconName = "icon_active.png"
	} else {
		iconName = "icon_default.png"
	}

	// List of possible icon locations (in order of preference)
	iconPaths := []string{
		fmt.Sprintf("assets/icons/%s", iconName),                                     // Development location
		fmt.Sprintf("Contents/Resources/assets/icons/%s", iconName),                  // macOS app bundle location
		fmt.Sprintf("%s/Contents/Resources/assets/icons/%s", getAppPath(), iconName), // Absolute path in app bundle
	}

	// Try each location until we find the icon
	for _, iconPath := range iconPaths {
		iconData, err := os.ReadFile(iconPath)
		if err == nil {
			return iconData
		}
	}

	// If we couldn't find the icon, log an error and return a placeholder
	fmt.Printf("Error loading icon %s: not found in any location\n", iconName)
	// Return a simple placeholder icon (1x1 transparent pixel)
	return []byte{0}
}

// getAppPath returns the path to the application bundle
func getAppPath() string {
	// Get the executable path
	execPath, err := os.Executable()
	if err != nil {
		fmt.Println("Error getting executable path:", err)
		return ""
	}

	// For a macOS app bundle, the executable is in Contents/MacOS/
	// So we need to go up two directories to get the app bundle path
	appPath := execPath
	for i := 0; i < 2; i++ {
		appPath = filepath.Dir(appPath)
	}

	return appPath
}

// isAutoStartEnabled checks if the app is set to start at login
func isAutoStartEnabled() bool {
	// Get the path to the app bundle
	appPath, err := exec.Command("pwd").Output()
	if err != nil {
		fmt.Println("Error getting current directory:", err)
		return false
	}

	// Trim newline from the output
	appPathStr := fmt.Sprintf("%s/build/Jolt.app", string(appPath[:len(appPath)-1]))

	// AppleScript to check if the app is in Login Items
	script := fmt.Sprintf(`
		tell application "System Events"
			get the name of every login item whose path contains "%s"
		end tell
	`, appPathStr)

	output, err := exec.Command("osascript", "-e", script).Output()
	if err != nil {
		fmt.Println("Error checking auto-start status:", err)
		return false
	}

	// If the output is not empty, the app is in Login Items
	return len(output) > 0
}

// enableAutoStart adds the app to Login Items
func enableAutoStart() {
	// Get the path to the app bundle
	appPath, err := exec.Command("pwd").Output()
	if err != nil {
		fmt.Println("Error getting current directory:", err)
		return
	}

	// Trim newline from the output
	appPathStr := fmt.Sprintf("%s/build/Jolt.app", string(appPath[:len(appPath)-1]))

	// AppleScript to add the app to Login Items
	script := fmt.Sprintf(`
		tell application "System Events"
			make login item at end with properties {path:"%s", hidden:false}
		end tell
	`, appPathStr)

	if err := exec.Command("osascript", "-e", script).Run(); err != nil {
		fmt.Println("Error enabling auto-start:", err)
	} else {
		autoStartEnabled = true
	}
}

// disableAutoStart removes the app from Login Items
func disableAutoStart() {
	// Get the path to the app bundle
	appPath, err := exec.Command("pwd").Output()
	if err != nil {
		fmt.Println("Error getting current directory:", err)
		return
	}

	// Trim newline from the output
	appPathStr := fmt.Sprintf("%s/build/Jolt.app", string(appPath[:len(appPath)-1]))

	// AppleScript to remove the app from Login Items
	script := fmt.Sprintf(`
		tell application "System Events"
			delete (every login item whose path contains "%s")
		end tell
	`, appPathStr)

	if err := exec.Command("osascript", "-e", script).Run(); err != nil {
		fmt.Println("Error disabling auto-start:", err)
	} else {
		autoStartEnabled = false
	}
}
