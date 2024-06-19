package main

import (
	_ "embed"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/charmbracelet/huh"
	"howett.net/plist"
)

//go:embed create_user.sh
var createUserScript string

type Config struct {
	ApiToken          string `json:"KANDJI_API_TOKEN"`
	ApiUrl            string `json:"KANDJI_API_URL"`
	StandardBlueprint string `json:"KANDJI_STANDARD_BLUEPRINT"`
	DevBlueprint      string `json:"KANDJI_DEV_BLUEPRINT"`
	TempPassword      string `json:"TEMP_PASSWORD"`
}

type KandjiProfileVars struct {
	// Needed Kandji variables from variables plist:
	// https://support.kandji.io/support/solutions/articles/72000560519-global-variables
	DeviceID    string `plist:"DEVICE_ID"`
	FullName    string `plist:"FULL_NAME"`
	EmailPrefix string `plist:"EMAIL_PREFIX"`
}

type SetupVars struct {
	// User options for the setup menu
	DeviceID  string
	FullName  string
	Username  string
	UserRole  string
	Confirm   bool
	Blueprint string
}

func main() {
	config, err := loadConfig("/tmp/config.json")
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	defer os.Remove("/tmp/config.json")

	var kv KandjiProfileVars = KandjiProfileVars{}
	// Get the Kandji variables from the plist
	kv, err = getPlistInfo()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	var sv SetupVars = SetupVars{
		DeviceID:  kv.DeviceID,
		FullName:  kv.FullName,
		Username:  kv.EmailPrefix,
		UserRole:  "standard", // default to standard
		Blueprint: "",         // default to empty
		Confirm:   false,      // default to false
	}

	// Verify there is a user assigned to the device
	if sv.FullName == "" {
		fmt.Println("No user assigned to device - please assign a user in the Kandji portal.")
		fmt.Println("If a user is already assigned, sync Kandji via Self Service or sudo kandji run. Then, try again.")
		os.Exit(1)
	}

	// Update setup vars with user input
	err = setupMenu(&sv)
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
	createResult, err := createUser(&sv, config)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	fmt.Println(createResult)

	// Send the API request to change the blueprint
	bpName, err := updateBlueprint(&sv, config)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	fmt.Printf("Blueprint updated to: %s\n", bpName)

	// Verify the new computer details with the Kandji API
	verifyResult, err := verifyNewComputerDetails(&sv, config)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	fmt.Println(verifyResult)

}

func setupMenu(sv *SetupVars) error {

	roleForm := huh.NewForm(
		huh.NewGroup(
			huh.NewSelect[string]().
				Title("Will the new user be standard or dev?").
				Options(
					huh.NewOption("Standard", "standard").Selected(true),
					huh.NewOption("Dev", "admin"),
				).
				Value(&sv.UserRole),
		),
	)

	if err := roleForm.Run(); err != nil {
		return fmt.Errorf("error with form: %v", err)
	}

	confirmMsg := fmt.Sprintf("Full Name: %s\nUsername: %s\nRole: %s\n", sv.FullName, sv.Username, sv.UserRole)
	confirmForm := huh.NewForm(
		huh.NewGroup(
			huh.NewSelect[bool]().
				Title("Is this information correct?").
				Description(confirmMsg).
				Options(
					huh.NewOption("Yes (Send Command)", true).Selected(true),
					huh.NewOption("No (Exit)", false),
				).
				Value(&sv.Confirm),
		),
	)

	if err := confirmForm.Run(); err != nil {
		return fmt.Errorf("error with form: %v", err)
	}

	if !sv.Confirm {
		fmt.Println("Exiting...")
		os.Exit(0)
	}

	return nil
}

func getPlistInfo() (KandjiProfileVars, error) {
	// Get the plist file from the Kandji variables
	plistPath := "/Library/Managed Preferences/io.kandji.globalvariables.plist"

	// Check if the file exists
	_, err := os.Stat(plistPath)
	if err != nil {
		return KandjiProfileVars{}, errors.New("kandji global variables plist does not exist")
	}

	// Open the file
	file, err := os.Open(plistPath)
	if err != nil {
		return KandjiProfileVars{}, fmt.Errorf("error opening file: %v", err)
	}
	defer file.Close()

	// Decode the plist file
	var kVars KandjiProfileVars
	decoder := plist.NewDecoder(file)
	err = decoder.Decode(&kVars)
	if err != nil {
		return KandjiProfileVars{}, fmt.Errorf("error decoding plist: %v", err)
	}

	return kVars, nil
}

func createUser(sv *SetupVars, c *Config) (string, error) {
	// Run the embedded shell script to create the user
	tmpFile, err := os.CreateTemp("", "create_user.sh")
	if err != nil {
		return "", fmt.Errorf("error creating temp file: %v", err)
	}
	defer os.Remove(tmpFile.Name())

	if _, err := tmpFile.WriteString(createUserScript); err != nil {
		return "", fmt.Errorf("error writing to temp file: %v", err)
	}

	if err := tmpFile.Chmod(0700); err != nil {
		return "", fmt.Errorf("error changing temp file permissions: %v", err)
	}

	cmd := exec.Command(tmpFile.Name(), sv.Username, c.TempPassword, sv.UserRole)
	output, err := cmd.CombinedOutput()
	if err != nil {
		outputStr := string(output)
		if strings.Contains(outputStr, "root") {
			return "", fmt.Errorf("shell script error: must be running as root")
		} else if strings.Contains(outputStr, "already exists") {
			return "", fmt.Errorf("shell script error: user already exists - please delete user and try again")
		} else if strings.Contains(outputStr, "failed") {
			return "", fmt.Errorf("shell script error: failed to create user")
		}
		return "", fmt.Errorf("shell script error: %v", err)
	}

	return fmt.Sprintf("User %s created successfully as %s", sv.Username, sv.UserRole), nil
}

func loadConfig(path string) (*Config, error) {
	file, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("error reading config file: %v", err)
	}

	var config Config
	if err := json.Unmarshal(file, &config); err != nil {
		return nil, fmt.Errorf("error unmarshalling config: %v", err)
	}

	return &config, nil
}
