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
		fmt.Println(err)
		os.Exit(1)
	}

	// Load external config file
	config, err := config.Load("/Library/UserSetup/config.json")
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	// Remove the config file after running
	defer os.Remove("/tmp/config.json") // TODO: Needed?

	// Get the Kandji variables from the plist
	var kVar kandji.KandjiProfileVars = kandji.KandjiProfileVars{}
	kVar, err = kandji.GetPlistInfo()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	// Get the Mac Model for the device name
	model, err := mac.GetModel()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	// Initialize setup variables
	var sVar shared.SetupVars = shared.SetupVars{
		DeviceID:     kVar.DeviceID,
		DeviceName:   "",
		FullName:     kVar.FullName,
		Username:     mac.CreateShortname(kVar.FullName),
		TempPassword: config.TempPassword,
		UserRole:     "standard",
		Blueprint:    "",
		Confirm:      false,
		DeleteSpare:  false,
	}

	// Set the device name based on the user's full name and the device model, if they are not empty
	if kVar.FullName != "" && model != "" {
		sVar.DeviceName = fmt.Sprintf("%s %s", kVar.FullName, model)
	}

	// Verify there is a user assigned to the device
	if sVar.FullName == "" {
		fmt.Println("No user assigned to device - please assign a user in the Kandji portal.")
		fmt.Println("If a user is already assigned, sync Kandji via Self Service or sudo kandji run. Then, try again.")
		os.Exit(1)
	}

	// Update setup vars with user input
	err = ui.RunMenu(&sVar)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	// Determine blueprint based on user role
	devBp := config.DevBlueprint
	standardBp := config.StandardBlueprint
	if sVar.UserRole == "dev" {
		sVar.Blueprint = devBp
	} else {
		sVar.Blueprint = standardBp
	}

	// Make the user with the shell script with output
	createResult, err := mac.CreateUser(&sVar)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	fmt.Println(createResult)

	// Send the API request to change the blueprint
	err = kandji.UpdateBlueprint(&sVar, config)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	// Send the API request to change the computer name
	err = kandji.UpdateComputerName(&sVar, config)
	if err != nil {
		fmt.Println(err)
	}

	// Delete the spare user if the user selected to do so
	if sVar.DeleteSpare {
		err = kandji.DeleteUser(&sVar, config, "spare")
		if err != nil {
			fmt.Println(err)
		}
	}

	// Verify the new computer details with the Kandji API
	details, err := kandji.GetComputerDetails(&sVar, config)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	summaryStr := fmt.Sprintf("Assigned User: %s\nBlueprint: %s\nDevice Name: %s\n", details.User.Name, details.BlueprintName, details.DeviceName)
	fmt.Println(summaryStr)

}
