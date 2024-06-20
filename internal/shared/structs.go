package shared

type SetupVars struct {
	// User options for the setup menu
	DeviceID    string
	FullName    string
	Username    string
	Password    string
	UserRole    string
	Confirm     bool
	Blueprint   string
	DeleteSpare bool
}
