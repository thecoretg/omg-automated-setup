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

func RunAssignedLogic(conf *config.Config, initialDetails kandji.DeviceDetails) {
	// Logic for when a user is assigned in Kandji - check if the user exists and create if not
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

	summary, err := CreateAssignedUser(conf, &initialDetails)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	fmt.Println(summary)
}

func CreateAssignedUser(conf *config.Config, initialDetails *kandji.DeviceDetails) (string, error) {
	// User creation logic for when a user is assigned in Kandji
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
	exists, err := mac.CheckUserExists(sVar.Username)
	if err != nil {
		return "", fmt.Errorf("error checking if user exists: %v", err)
	}

	// Store post-user creation instructions for end of program
	loginInstr := "IMPORTANT: Log out of the current user and log in to the new user to ensure it gets secure token before shutting down."
	var userInstr string
	switch exists {
	case true:
		// If the user already exists, skip user creation
		userInstr = fmt.Sprintf("User %s already exists on this Mac. User creation was skipped.\n%s", sVar.Username, loginInstr)

	case false:
		// If the user does not exist, create the user
		err := mac.CreateUser(&sVar)
		if err != nil {
			return "", fmt.Errorf("error creating user: %v", err)
		}
		userInstr = fmt.Sprintf("User %s (%s user) created with default password.\n%s", sVar.Username, sVar.UserRole, loginInstr)
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
		err = kandji.DeleteUser(&sVar, conf, "spareuser")
		if err != nil {
			return "", fmt.Errorf("error deleting spare user: %v", err)
		}
	}

	// Verify the new computer details with the Kandji API
	details, err := kandji.GetComputerDetails(sVar.DeviceID, conf)
	if err != nil {
		return "", fmt.Errorf("error getting computer details: %v", err)
	}

	summaryStr := fmt.Sprintf("Assigned User: %s\nBlueprint: %s\n\n%s\n\n", details.User.Name, details.BlueprintName, userInstr)
	return summaryStr, nil
}
