package ui

import (
	"fmt"

	"github.com/charmbracelet/huh"
	"github.com/thecoretg/omg-user-automation/internal/types"
)

func RunMenu(sv *types.SetupVars) error {
	form := huh.NewForm(
		SetupTypeMenu(sv),
	)

	if err := form.WithTheme(huh.ThemeBase()).Run(); err != nil {
		return fmt.Errorf("error with form: %v", err)
	}

	switch sv.SetupType {
	case "Spare":
		// Run the spare setup menu
		return fmt.Errorf("spare setup not implemented yet")

	case "User":
		// Run the user setup menu
		if err := RunUserMenu(sv); err != nil {
			return fmt.Errorf("error running user menu: %v", err)
		}
		return nil

	default:
		return fmt.Errorf("invalid setup type: %s", sv.SetupType)

	}
}

func RunUserMenu(sv *types.SetupVars) error {
	form := huh.NewForm(
		InitialUserMenu(sv),
		UserConfirmMenu(sv),
	)

	if err := form.WithTheme(huh.ThemeBase()).Run(); err != nil {
		return fmt.Errorf("error with form: %v", err)
	}

	return nil
}
func SetupTypeMenu(sv *types.SetupVars) *huh.Group {
	return huh.NewGroup(
		huh.NewSelect[string]().
			Title("Setup Type").
			Options(
				huh.NewOptions("Spare", "User")...).
			Value(&sv.SetupType),
	)
}

func InitialUserMenu(sv *types.SetupVars) *huh.Group {
	// User inputs role type and if they want to delete the spare user (if it exists)
	return huh.NewGroup(
		huh.NewSelect[string]().
			Title("Will the new user be standard or dev?").
			Options(
				huh.NewOption("Standard", "standard").Selected(true),
				huh.NewOption("Dev", "admin"),
			).
			Value(&sv.UserRole),
		huh.NewSelect[bool]().
			Title("Delete spare user, if it exists?").
			Options(
				huh.NewOption("Yes", true).Selected(true),
				huh.NewOption("No", false),
			).
			Value(&sv.DeleteSpare),
	)
}

func UserConfirmMenu(sv *types.SetupVars) *huh.Group {
	confirmMsg := fmt.Sprintf("Computer Name: %s\nFull Name: %s\nUsername: %s\nRole: %s\n", sv.DeviceName, sv.FullName, sv.Username, sv.UserRole)
	return huh.NewGroup(
		huh.NewSelect[bool]().
			Title("Is this information correct?").
			Description(confirmMsg).
			Options(
				huh.NewOption("Yes (Send Command)", true).Selected(true),
				huh.NewOption("No (Exit)", false),
			).
			Value(&sv.Confirm),
	)
}
