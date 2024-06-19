package ui

import (
	"fmt"

	"github.com/charmbracelet/huh"
	"github.com/thecoretg/omg-user-automation/internal/shared"
)

func RunSetupTypeMenu(sv *shared.SetupVars) error {
	form := huh.NewForm(
		huh.NewGroup(
			huh.NewSelect[string]().
				Title("Setup Type").
				Options(
					huh.NewOptions("Spare", "User")...).
				Value(&sv.SetupType),
		),
	)

	if err := form.WithTheme(huh.ThemeBase()).WithAccessible(true).Run(); err != nil {
		return fmt.Errorf("error with setup type form: %v", err)
	}

	return nil
}

func RunUserMenu(sv *shared.SetupVars) error {
	form := huh.NewForm(
		UserRoleMenu(sv),
		DeleteSpareMenu(sv),
		UserConfirmMenu(sv),
	)

	if err := form.WithTheme(huh.ThemeBase()).WithAccessible(true).Run(); err != nil {
		return fmt.Errorf("error with user setup form: %v", err)
	}

	return nil
}

func UserRoleMenu(sv *shared.SetupVars) *huh.Group {
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

func DeleteSpareMenu(sv *shared.SetupVars) *huh.Group {
	// User inputs role type and if they want to delete the spare user (if it exists)
	return huh.NewGroup(
		huh.NewSelect[bool]().
			Title("Delete spare user, if it exists?").
			Options(
				huh.NewOption("Yes", true).Selected(true),
				huh.NewOption("No", false),
			).
			Value(&sv.DeleteSpare),
	)
}

func UserConfirmMenu(sv *shared.SetupVars) *huh.Group {
	confirmMsg := fmt.Sprintf("Full Name: %s\nUsername: %s\nRole: %s\n", sv.FullName, sv.Username, sv.UserRole)
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
