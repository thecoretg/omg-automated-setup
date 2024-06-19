package mac

import (
	"bytes"
	"errors"
	"fmt"
	"os/exec"

	"github.com/thecoretg/omg-user-automation/internal/types"
)

func CreateUser(sv *types.SetupVars) (string, error) {
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

func GetModel() (string, error) {
	// Gets the Mac Model name, such as "MacBook Pro" or "Mac Mini"
	cmdString := `/usr/sbin/system_profiler SPHardwareDataType | awk -F": " '/Model Name/ {print $2}'`
	cmd := exec.Command("bash", "-c", cmdString)

	var out bytes.Buffer
	cmd.Stdout = &out
	err := cmd.Run()
	if err != nil {
		return "", fmt.Errorf("error running mac model retrieval command: %v", err)
	}

	macModel := out.String()
	if macModel == "" {
		return "", errors.New("computer name is empty")
	}

	return macModel, nil

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
