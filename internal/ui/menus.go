package ui

import (
	"fmt"

	"github.com/charmbracelet/huh"
	"github.com/thecoretg/omg-user-automation/internal/shared"
)

func RunSetupTypeMenu() (string, error) {
	var setupType string
	form := huh.NewForm(
		huh.NewGroup(
			huh.NewSelect[string]().
				Title("Setup Type").
				Options(
					huh.NewOptions("Spare", "User")...).
				Value(&setupType),
		),
	)

	if err := form.WithTheme(huh.ThemeBase()).Run(); err != nil {
		return "", fmt.Errorf("error with setup type form: %v", err)
	}

	return setupType, nil
}

func YesNoMenu(title string) (bool, error) {
	var confirm bool
	form := huh.NewForm(
		huh.NewGroup(
			huh.NewSelect[bool]().
				Title(title).
				Options(
					huh.NewOption("Yes", true).Selected(true),
					huh.NewOption("No", false),
				).
				Value(&confirm),
		),
	)

	if err := form.WithTheme(huh.ThemeBase()).Run(); err != nil {
		return false, fmt.Errorf("error with yes/no form: %v", err)
	}

	return confirm, nil
}

func RunUserMenu(sv *shared.SetupVars) error {
	userForm := huh.NewForm(
		UserRoleMenu(sv),
		DeleteSpareMenu(sv),
	)

	if err := userForm.WithTheme(huh.ThemeBase()).Run(); err != nil {
		return fmt.Errorf("error with user setup form: %v", err)
	}

	confirmForm := UserConfirmMenu(sv)
	if err := confirmForm.WithTheme(huh.ThemeBase()).Run(); err != nil {
		return fmt.Errorf("error with user confirm form: %v", err)
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

func UserConfirmMenu(sv *shared.SetupVars) *huh.Form {
	confirmMsg := fmt.Sprintf("Full Name: %s\nUsername: %s\nRole: %s\n", sv.FullName, sv.Username, sv.UserRole)
	return huh.NewForm(
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
}
