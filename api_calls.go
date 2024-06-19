package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
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

func kandjiApiRequest(apiVars RequestVars, c *Config) (string, error) {
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

func updateBlueprint(sv *SetupVars, c *Config) (string, error) {
	reqVars := RequestVars{
		Method:   "PATCH",
		Endpoint: fmt.Sprintf("devices/%s", sv.DeviceID),
		Payload:  fmt.Sprintf(`{"blueprint_id": "%s"}`, sv.Blueprint),
	}

	resp, err := kandjiApiRequest(reqVars, c)
	if err != nil {
		return "", fmt.Errorf("error updating blueprint: %v", err)
	}

	// Parse the response into an UpdateResponse struct
	var updateResp UpdateResponse
	if err := json.Unmarshal([]byte(resp), &updateResp); err != nil {
		return "", fmt.Errorf("error unmarshalling device update response: %v", err)
	}

	bpName := updateResp.BlueprintName

	return bpName, nil
}

func verifyNewComputerDetails(sv *SetupVars, c *Config) (string, error) {
	// Call the Kandji API to verify the new computer details
	reqVars := RequestVars{
		Method:   "GET",
		Endpoint: fmt.Sprintf("devices/%s", sv.DeviceID),
		Payload:  "",
	}

	resp, err := kandjiApiRequest(reqVars, c)
	if err != nil {
		return "", fmt.Errorf("error verifying new computer details: %v", err)
	}

	// Parse the response into a DeviceDetails struct
	var deviceDetails DeviceDetails
	if err := json.Unmarshal([]byte(resp), &deviceDetails); err != nil {
		return "", fmt.Errorf("error unmarshalling device detail response: %v", err)
	}

	// Print the new computer details: user's full name, blueprint name, and device name
	deviceSummary := fmt.Sprintf("Assigned User: %s\nBlueprint: %s\nDevice Name: %s\n", deviceDetails.User.Name, deviceDetails.BlueprintName, deviceDetails.DeviceName)
	return deviceSummary, nil
}
