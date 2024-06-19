package mac

import (
	"errors"
	"fmt"
	"os/exec"
	"strings"

	"github.com/thecoretg/omg-user-automation/internal/shared"
)

func CreateUser(sv *shared.SetupVars) (string, error) {
	// If UserRole is "admin", create the user with the admin role
	// If UserRole is "standard", create the user with the standard role
	cmd := exec.Command("sysadminctl", "-addUser", sv.Username, "-password", sv.TempPassword, "-fullName", sv.FullName)
	if sv.UserRole == "admin" {
		cmd.Args = append(cmd.Args, "-admin")
	}

	err := cmd.Run()
	if err != nil {
		return "", fmt.Errorf("error creating user: %v", err)
	}

	return fmt.Sprintf("User %s created as %s", sv.Username, sv.UserRole), nil
}

func CreateShortname(fullName string) string {
	loweCaseName := strings.ToLower(fullName)
	username := strings.ReplaceAll(loweCaseName, " ", "")
	return username
}

func CheckRoot() error {
	// Check if the script is running as root
	cmd := exec.Command("whoami")
	output, err := cmd.Output()
	if err != nil {
		return fmt.Errorf("error checking root: %v", err)
	}

	if string(output) != "root\n" {
		return errors.New("script must be run as root - please run with sudo, or as root")
	}

	return nil
}
