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

func RunSpareLogic(conf *config.Config, initialDetails kandji.DeviceDetails) {
	// Logic for when no user is assigned in Kandji - check if the spare user exists and create if not
	exists, err := mac.CheckUserExists("spareuser")
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	if exists {
		fmt.Println("No assigned user in Kandji, but Spare User already exists. Exiting program.")
		os.Exit(0)
	}

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

	summary, err := CreateSpareUser(conf, &initialDetails)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	fmt.Println(summary)
}

func CreateSpareUser(conf *config.Config, initialDetails *kandji.DeviceDetails) (string, error) {
	// User creation logic for when no user is assigned in Kandji
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
	err := mac.CreateUser(&sVar)
	if err != nil {
		return "", fmt.Errorf("error creating user: %v", err)
	}

	userInstr := "Spare User (spareuser) created with default password.\nIMPORTANT: Log out of the current user and log in to the spare user to ensure it gets secure token before shutting down."
	return userInstr, nil

}
