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

	switch initialDetails.User.Name {
	case "":
		title := "No assigned user in Kandji. Continue with Spare User setup?"
		// Run yes/no menu to confirm
		confirm, err := ui.YesNoMenu(title)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		if !confirm {
			fmt.Println("Exiting program.")
			os.Exit(0)
		}

		summary, err := RunSpareSetup(config, &initialDetails)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		fmt.Println(summary)

	default:
		// If assigned user is not empty, run the user setup
		title := fmt.Sprintf("Assigned user detected in Kandji: %s. Continue with user setup?", initialDetails.User.Name)
		confirm, err := ui.YesNoMenu(title)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		if !confirm {
			fmt.Println("Exiting program.")
			os.Exit(0)
		}

		summary, err := RunUserSetup(config, &initialDetails)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		fmt.Println(summary)
	}
}

func RunSpareSetup(conf *config.Config, initialDetails *kandji.DeviceDetails) (string, error) {
	// Initialize setup variables
	var sVar shared.SetupVars = shared.SetupVars{
		DeviceID:    initialDetails.DeviceID,
		FullName:    "Spare User",
		Username:    "spareuser",
		Password:    conf.SparePassword,
		UserRole:    "standard",
		Blueprint:   "",
		Confirm:     false,
		DeleteSpare: false,
	}

	// Make the spare user with the shell script with output
	createResult, err := mac.CreateUser(&sVar)
	if err != nil {
		return "", fmt.Errorf("error creating user: %v", err)
	} else { // If the user was created, print the output
		fmt.Println(createResult)
	}

	userInstr := "Spare User (spareuser) created with default password.\nIMPORTANT: Log out of the current user and log in to the spare user to ensure it gets secure token before shutting down."
	return userInstr, nil

}

func RunUserSetup(conf *config.Config, initialDetails *kandji.DeviceDetails) (string, error) {

	// Initialize setup variables
	var sVar shared.SetupVars = shared.SetupVars{
		DeviceID:    initialDetails.DeviceID,
		FullName:    initialDetails.User.Name,
		Username:    mac.CreateShortname(initialDetails.User.Name),
		Password:    conf.TempPassword,
		UserRole:    "standard",
		Blueprint:   "",
		Confirm:     false,
		DeleteSpare: false,
	}

	err := ui.RunUserMenu(&sVar)
	if err != nil {
		return "", fmt.Errorf("error running user setup menu: %v", err)
	}

	// Make the user with the shell script with output
	_, err = mac.CreateUser(&sVar)
	if err != nil {
		return "", fmt.Errorf("error creating user: %v", err)
	}

	// Determine blueprint based on user role
	devBp := conf.DevBlueprint
	standardBp := conf.StandardBlueprint
	if sVar.UserRole == "admin" {
		sVar.Blueprint = devBp
	} else {
		sVar.Blueprint = standardBp
	}

	// Send the API request to change the blueprint
	err = kandji.UpdateBlueprint(&sVar, conf)
	if err != nil {
		return "", fmt.Errorf("error updating blueprint: %v", err)
	}

	// Delete the spare user if the user selected to do so
	if sVar.DeleteSpare {
		err = kandji.DeleteUser(&sVar, conf, "spare")
		if err != nil {
			return "", fmt.Errorf("error deleting spare user: %v", err)
		}
	}

	// Verify the new computer details with the Kandji API
	details, err := kandji.GetComputerDetails(sVar.DeviceID, conf)
	if err != nil {
		return "", fmt.Errorf("error getting computer details: %v", err)
	}

	userInstr := fmt.Sprintf("User %s (%s user) created with default password.\nIMPORTANT: Log out of the current user and log in to the spare user to ensure it gets secure token before shutting down.", sVar.Username, sVar.UserRole)
	summaryStr := fmt.Sprintf("Assigned User: %s\nBlueprint: %s\n\n%s\n", details.User.Name, details.BlueprintName, userInstr)
	return summaryStr, nil
}
