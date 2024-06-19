package kandji

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/thecoretg/omg-user-automation/internal/types"
)

type RequestVars struct {
	// For Kandji API request
	Method   string
	Endpoint string
	Payload  string
}

type DeviceDetails struct {
	DeviceID                   string `json:"device_id"`
	DeviceName                 string `json:"device_name"`
	Model                      string `json:"model"`
	SerialNumber               string `json:"serial_number"`
	Platform                   string `json:"platform"`
	OSVersion                  string `json:"os_version"`
	SupplementalBuildVersion   string `json:"supplemental_build_version"`
	SupplementalOSVersionExtra string `json:"supplemental_os_version_extra"`
	User                       User   `json:"user"`
	AssetTag                   string `json:"asset_tag"`
	BlueprintID                string `json:"blueprint_id"`
	MdmEnabled                 bool   `json:"mdm_enabled"`
	AgentInstalled             bool   `json:"agent_installed"`
	IsMissing                  bool   `json:"is_missing"`
	IsRemoved                  bool   `json:"is_removed"`
	AgentVersion               string `json:"agent_version"`
	BlueprintName              string `json:"blueprint_name"`
}

type UpdateResponse struct {
	DeviceID       string `json:"device_id"`
	DeviceName     string `json:"device_name"`
	Model          string `json:"model"`
	SerialNumber   string `json:"serial_number"`
	Platform       string `json:"platform"`
	OSVersion      string `json:"os_version"`
	User           User   `json:"user"`
	AssetTag       string `json:"asset_tag"`
	BlueprintID    string `json:"blueprint_id"`
	MdmEnabled     bool   `json:"mdm_enabled"`
	AgentInstalled bool   `json:"agent_installed"`
	IsMissing      bool   `json:"is_missing"`
	IsRemoved      bool   `json:"is_removed"`
	AgentVersion   string `json:"agent_version"`
	BlueprintName  string `json:"blueprint_name"`
}

type User struct {
	Email      string `json:"email"`
	Name       string `json:"name"`
	ID         int64  `json:"id"`
	IsArchived bool   `json:"is_archived"`
}

type DeleteUserPayload struct {
	DeleteAllUsers bool   `json:"DeleteAllUsers"`
	ForceDeletion  bool   `json:"ForceDeletion"`
	UserName       string `json:"UserName"`
}

func ApiRequest(apiVars RequestVars, c *types.Config) (string, error) {
	// Build and run Kandji API request using RequestVars struct

	apiToken := c.ApiToken
	apiUrl := c.ApiUrl

	url := fmt.Sprintf("https://%s/%s", apiUrl, apiVars.Endpoint)

	payloadStr := apiVars.Payload
	payload := strings.NewReader(payloadStr)

	client := &http.Client{}
	req, err := http.NewRequest(apiVars.Method, url, payload)

	if err != nil {
		return "", fmt.Errorf("error creating request: %v", err)
	}

	req.Header = http.Header{
		"Authorization": []string{fmt.Sprintf("Bearer %s", apiToken)},
		"Content-Type":  []string{"application/json"},
	}

	resp, err := client.Do(req)
	if err != nil {
		return "", fmt.Errorf("error sending request: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return "", fmt.Errorf("error response: %v", resp.Status)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("error reading response: %v", err)
	}

	return string(body), nil
}

func UpdateBlueprint(sv *types.SetupVars, c *types.Config) error {
	reqVars := RequestVars{
		Method:   "PATCH",
		Endpoint: fmt.Sprintf("devices/%s", sv.DeviceID),
		Payload:  fmt.Sprintf(`{"blueprint_id": "%s"}`, sv.Blueprint),
	}

	_, err := ApiRequest(reqVars, c)
	if err != nil {
		return fmt.Errorf("error updating blueprint: %v", err)
	}

	return nil
}

func UpdateComputerName(sv *types.SetupVars, c *types.Config) error {

	reqVars := RequestVars{
		Method:   "POST",
		Endpoint: fmt.Sprintf("devices/%s/action/setname", sv.DeviceID),
		Payload:  fmt.Sprintf(`{"DeviceName": "%s"}`, sv.DeviceName),
	}

	_, err := ApiRequest(reqVars, c)
	if err != nil {
		return fmt.Errorf("error changing computer name: %v", err)
	}

	return nil
}

func DeleteUser(sv *types.SetupVars, c *types.Config, user string) error {
	payloadStruct := DeleteUserPayload{
		DeleteAllUsers: false,
		ForceDeletion:  false,
		UserName:       user,
	}

	payloadBytes, err := json.Marshal(payloadStruct)
	if err != nil {
		return fmt.Errorf("error marshalling payload: %v", err)
	}

	reqVars := RequestVars{
		Method:   "POST",
		Endpoint: fmt.Sprintf("devices/%s/action/deleteuser", sv.DeviceID),
		Payload:  string(payloadBytes),
	}

	_, err = ApiRequest(reqVars, c)
	if err != nil {
		return fmt.Errorf("error deleting user: %v", err)
	}

	return nil
}

func GetComputerDetails(sv *types.SetupVars, c *types.Config) (DeviceDetails, error) {
	// Call the Kandji API to verify the new computer details
	reqVars := RequestVars{
		Method:   "GET",
		Endpoint: fmt.Sprintf("devices/%s", sv.DeviceID),
		Payload:  "",
	}

	resp, err := ApiRequest(reqVars, c)
	if err != nil {
		return DeviceDetails{}, fmt.Errorf("error verifying new computer details: %v", err)
	}

	// Parse the response into a DeviceDetails struct
	var deviceDetails DeviceDetails
	if err := json.Unmarshal([]byte(resp), &deviceDetails); err != nil {
		return DeviceDetails{}, fmt.Errorf("error unmarshalling response: %v", err)
	}

	return deviceDetails, nil

}
