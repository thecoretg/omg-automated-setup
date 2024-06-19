package shared

type SetupVars struct {
	// User options for the setup menu
	SetupType    string
	DeviceID     string
	DeviceName   string
	FullName     string
	Username     string
	TempPassword string
	UserRole     string
	Confirm      bool
	Blueprint    string
	DeleteSpare  bool
}
