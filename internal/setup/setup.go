package setup

import (
	"fmt"
	"os"

	"github.com/thecoretg/omg-user-automation/internal/config"
	"github.com/thecoretg/omg-user-automation/internal/kandji"
	"github.com/thecoretg/omg-user-automation/internal/mac"
	"github.com/thecoretg/omg-user-automation/internal/shared"
	"github.com/thecoretg/omg-user-automation/internal/ui"
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

	// Initialize setup variables
	var sVar shared.SetupVars = shared.SetupVars{
		SetupType:    "",
		DeviceID:     initialDetails.DeviceID,
		FullName:     initialDetails.User.Name,
		Username:     mac.CreateShortname(initialDetails.User.Name),
		TempPassword: config.TempPassword,
		UserRole:     "standard",
		Blueprint:    "",
		Confirm:      false,
		DeleteSpare:  false,
	}

	// Run setup menu to determine Spare or User
	err = ui.RunSetupTypeMenu(&sVar)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	switch sVar.SetupType {
	case "Spare":
		// If the user selects Spare, run the spare setup which doesnt exist yet
		fmt.Println("Spare setup not yet implemented.")
		os.Exit(0)
	case "User":
		// If the user selects User, run user setup menu
		summary, err := runUserSetup(config, &sVar)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		fmt.Println(summary)
	}
}

func runUserSetup(conf *config.Config, sVar *shared.SetupVars) (string, error) {

	err := ui.RunUserMenu(sVar)
	if err != nil {
		return "", fmt.Errorf("error running user setup menu: %v", err)
	}

	// Make the user with the shell script with output
	createResult, err := mac.CreateUser(sVar)
	if err != nil {
		return "", fmt.Errorf("error creating user: %v", err)
	} else { // If the user was created, print the output
		fmt.Println(createResult)
	}

	// Determine blueprint based on user role
	devBp := conf.DevBlueprint
	standardBp := conf.StandardBlueprint
	if sVar.UserRole == "dev" {
		sVar.Blueprint = devBp
	} else {
		sVar.Blueprint = standardBp
	}

	// Send the API request to change the blueprint
	err = kandji.UpdateBlueprint(sVar, conf)
	if err != nil {
		return "", fmt.Errorf("error updating blueprint: %v", err)
	}

	// Delete the spare user if the user selected to do so
	if sVar.DeleteSpare {
		err = kandji.DeleteUser(sVar, conf, "spare")
		if err != nil {
			return "", fmt.Errorf("error deleting spare user: %v", err)
		}
	}

	// Verify the new computer details with the Kandji API
	details, err := kandji.GetComputerDetails(sVar.DeviceID, conf)
	if err != nil {
		return "", fmt.Errorf("error getting computer details: %v", err)
	}

	summaryStr := fmt.Sprintf("Assigned User: %s\nBlueprint: %s\n", details.User.Name, details.BlueprintName)
	return summaryStr, nil
}
