package mac

import (
	"errors"
	"fmt"
	"os/exec"
	"os/user"
	"strings"

	"github.com/thecoretg/omg-user-automation/internal/shared"
)

func CreateUser(sv *shared.SetupVars) (string, error) {
	// If UserRole is "admin", create the user with the admin role
	// If UserRole is "standard", create the user with the standard role
	cmd := exec.Command("sysadminctl", "-addUser", sv.Username, "-password", sv.Password, "-fullName", sv.FullName)
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

	currentUser, err := user.Current()
	if err != nil {
		return fmt.Errorf("error getting current user: %v", err)
	}

	if currentUser.Username != "root" {
		return errors.New("script must be run as root - please run with sudo, or as root")
	}

	return nil
}

func CheckUserExists(username string) (bool, error) {
	_, err := user.Lookup(username)
	if err != nil {
		return false, nil
	}

	return true, nil
}
