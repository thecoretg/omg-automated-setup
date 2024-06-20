package setup

import (
	"fmt"
	"os"

	"github.com/thecoretg/omg-user-automation/internal/config"
	"github.com/thecoretg/omg-user-automation/internal/kandji"
	"github.com/thecoretg/omg-user-automation/internal/mac"
)

func RunProgram() {
	// The main function that runs the program - run this function in main.go

	// Check if root - if not, exit. We'll need root to create a user
	err := mac.CheckRoot()
	if err != nil {
		fmt.Println("This command must be run as root. Please use sudo or switch to the root user.")
		os.Exit(1)
	}

	// Load external config file
	config, err := config.Load("/Library/UserSetup/config.json")
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	// Get the Kandji device ID from the plist
	devId, err := kandji.GetDeviceID()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	// Use device ID to get all of the initial details from the Kandji API
	initialDetails, err := kandji.GetComputerDetails(devId, config)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	switch initialDetails.User.Name {
	case "":
		// If assigned user is empty, run the spare user setup
		RunSpareLogic(config, initialDetails)

	default:
		// If assigned user is not empty, run the user setup
		RunAssignedLogic(config, initialDetails)
	}
}
