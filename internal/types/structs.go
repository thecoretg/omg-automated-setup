package types

type Config struct {
	ApiToken          string `json:"KANDJI_API_TOKEN"`
	ApiUrl            string `json:"KANDJI_API_URL"`
	StandardBlueprint string `json:"KANDJI_STANDARD_BLUEPRINT"`
	DevBlueprint      string `json:"KANDJI_DEV_BLUEPRINT"`
	TempPassword      string `json:"TEMP_PASSWORD"`
}

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
