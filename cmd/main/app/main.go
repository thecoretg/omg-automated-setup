package main

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"

	"github.com/thecoretg/omg-user-automation/internal/kandji"
	"github.com/thecoretg/omg-user-automation/internal/mac"
	"github.com/thecoretg/omg-user-automation/internal/types"
	"github.com/thecoretg/omg-user-automation/internal/ui"
)

func main() {

	// Check if the script is running as root
	err := mac.CheckRoot()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	config, err := loadConfig("/tmp/config.json")
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	defer os.Remove("/tmp/config.json")

	var kv kandji.KandjiProfileVars = kandji.KandjiProfileVars{}
	// Get the Kandji variables from the plist
	kv, err = kandji.GetPlistInfo()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	// Get the Mac Model
	macModel, err := mac.GetModel()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	var sv types.SetupVars = types.SetupVars{
		DeviceID:     kv.DeviceID,
		DeviceName:   "", // default to empty
		FullName:     kv.FullName,
		Username:     createUsername(kv.FullName),
		TempPassword: config.TempPassword,
		UserRole:     "standard", // default to standard
		Blueprint:    "",         // default to empty
		Confirm:      false,      // default to false
		DeleteSpare:  false,      // default to false
	}

	// Set the device name based on the user's full name and the device model, if they are not empty
	if kv.FullName != "" && macModel != "" {
		sv.DeviceName = fmt.Sprintf("%s %s", kv.FullName, macModel)
	}

	// Verify there is a user assigned to the device
	if sv.FullName == "" {
		fmt.Println("No user assigned to device - please assign a user in the Kandji portal.")
		fmt.Println("If a user is already assigned, sync Kandji via Self Service or sudo kandji run. Then, try again.")
		os.Exit(1)
	}

	// Update setup vars with user input
	err = ui.RunMenu(&sv)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	// Determine blueprint based on user role
	devBp := config.DevBlueprint
	standardBp := config.StandardBlueprint
	if sv.UserRole == "dev" {
		sv.Blueprint = devBp
	} else {
		sv.Blueprint = standardBp
	}

	// Make the user with the shell script with output
	createResult, err := mac.CreateUser(&sv)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	fmt.Println(createResult)

	// Send the API request to change the blueprint
	err = kandji.UpdateBlueprint(&sv, config)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	// Send the API request to change the computer name
	err = kandji.UpdateComputerName(&sv, config)
	if err != nil {
		fmt.Println(err)
	}

	// Delete the spare user if the user selected to do so
	if sv.DeleteSpare {
		err = kandji.DeleteUser(&sv, config, "spare")
		if err != nil {
			fmt.Println(err)
		}
	}

	// Verify the new computer details with the Kandji API
	details, err := kandji.GetComputerDetails(&sv, config)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	summaryStr := fmt.Sprintf("Assigned User: %s\nBlueprint: %s\nDevice Name: %s\n", details.User.Name, details.BlueprintName, details.DeviceName)
	fmt.Println(summaryStr)

}

func createUsername(fullName string) string {
	loweCaseName := strings.ToLower(fullName)
	username := strings.ReplaceAll(loweCaseName, " ", "")
	return username
}

func loadConfig(path string) (*types.Config, error) {
	file, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("error reading config file: %v", err)
	}

	var config types.Config
	if err := json.Unmarshal(file, &config); err != nil {
		return nil, fmt.Errorf("error unmarshalling config: %v", err)
	}

	return &config, nil
}
