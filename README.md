# Jolt

Jolt is a lightweight macOS menu bar application that prevents your Mac from going to sleep, dimming your screen, or starting the screensaver.

## Features

- Lives in the menu bar (no dock icon)
- Controllable via the menu bar icon
- Integrated timer with optional notifications
- Support for both retina displays and dark mode

## Requirements

- macOS 10.13 or later
- Go 1.24 or later

## Building from Source

1. Clone the repository:
   ```
   git clone https://github.com/nilicule/joltApp.git
   cd jolt
   ```

2. Install dependencies:
   ```
   go get github.com/getlantern/systray
   ```

3. Build the application:
   ```
   chmod +x build.sh
   ./build.sh
   ```

4. The build script will:
   - Create a macOS application bundle (`build/Jolt.app`) that you can run by double-clicking on it
   - Package the application into a zip file with the version number from `version.txt`
   - Store the zip file in the `releases` directory

## Usage

- **Enable/Disable**: Click on the "Enable Jolt" menu item to toggle sleep prevention.
- **Timer**: Set a timer for sleep prevention using the "Set Timer" submenu.
- **Start at Login**: Click on the "Start at Login" menu item to make Jolt start automatically when you log in to your Mac. A checkmark indicates when this feature is enabled.
- **Quit**: Click on the "Quit" menu item to exit the application.

## Icon Implementation

Jolt supports custom icon files for the menu bar:

- When sleep prevention is active: Uses `icon_active.png` from the assets/icons directory
- When sleep prevention is inactive: Uses `icon_default.png` from the assets/icons directory

Jolt now displays only the icon in the menu bar, without any text, providing a cleaner and more minimal interface.

### Customizing Icons

To use custom icons with Jolt:

1. Create your icon files according to the specifications in the `assets/icons/README.md` file
2. Place them in the `assets/icons` directory with the following names:
   - `icon_default.png` - The default icon shown when sleep prevention is inactive
   - `icon_active.png` - The icon shown when sleep prevention is active
3. Rebuild the application using the build.sh script

See the README in the `assets/icons` directory for more details on the implementation.

## How It Works

Jolt uses the macOS `caffeinate` command to prevent sleep, screen dimming, and screensaver activation. When enabled, it runs the `caffeinate` command with the following options:

- `-d`: Prevent the display from sleeping
- `-i`: Prevent the system from idle sleeping
- `-s`: Prevent the system from sleeping
- `-u`: Declare that the user is active (to prevent the screensaver from starting)

## License

[MIT License](LICENSE)
