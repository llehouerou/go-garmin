// service_device.go
package garmin

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

// Device represents a Garmin device with its capabilities.
type Device struct {
	DeviceID           int64  `json:"deviceId"`
	UnitID             int64  `json:"unitId"`
	DisplayName        string `json:"displayName"`
	ProductDisplayName string `json:"productDisplayName"`
	SerialNumber       string `json:"serialNumber"`
	DeviceStatus       string `json:"deviceStatus"`
	PartNumber         string `json:"partNumber"`
	ProductSku         string `json:"productSku"`
	ApplicationKey     string `json:"applicationKey"`
	ImageURL           string `json:"imageUrl"`
	RegisteredDate     int64  `json:"registeredDate"`

	// Firmware info
	CurrentFirmwareVersion      string `json:"currentFirmwareVersion"`
	CurrentFirmwareVersionMajor int    `json:"currentFirmwareVersionMajor"`
	CurrentFirmwareVersionMinor int    `json:"currentFirmwareVersionMinor"`

	// Device type info
	DeviceTypePK         int64    `json:"deviceTypePk"`
	DeviceTypeSimpleName string   `json:"deviceTypeSimpleName"`
	DeviceCategories     []string `json:"deviceCategories"`

	// Connection capabilities
	WiFi                     bool `json:"wifi"`
	BluetoothClassicDevice   bool `json:"bluetoothClassicDevice"`
	BluetoothLowEnergyDevice bool `json:"bluetoothLowEnergyDevice"`

	// Status indicators
	Primary                         bool `json:"primary"`
	ActiveInd                       int  `json:"activeInd"`
	PrimaryActivityTrackerIndicator bool `json:"primaryActivityTrackerIndicator"`

	// Capability flags (common ones)
	BodyBatteryCapable       bool `json:"bodyBatteryCapable"`
	HRVStatusCapable         bool `json:"hrvStatusCapable"`
	TrainingReadinessCapable bool `json:"trainingReadinessCapable"`
	VO2MaxRunCapable         bool `json:"vo2MaxRunCapable"`
	VO2MaxBikeCapable        bool `json:"vo2MaxBikeCapable"`
	PulseOxAllDayCapable     bool `json:"pulseOxAllDayCapable"`
	PulseOxSleepCapable      bool `json:"pulseOxSleepCapable"`
	SleepScoreCapable        bool `json:"sleepScoreCapable"`
	AllDayStressCapable      bool `json:"allDayStressCapable"`
	RespirationCapable       bool `json:"respirationCapable"`
	HasOpticalHeartRate      bool `json:"hasOpticalHeartRate"`

	raw json.RawMessage
}

// RawJSON returns the original JSON response.
func (d *Device) RawJSON() json.RawMessage {
	return d.raw
}

// DeviceAlarm represents an alarm configuration on a device.
type DeviceAlarm struct {
	AlarmID      int64    `json:"alarmId"`
	AlarmMode    string   `json:"alarmMode"`
	AlarmTime    int      `json:"alarmTime"`
	AlarmDays    []string `json:"alarmDays"`
	AlarmSound   string   `json:"alarmSound"`
	ChangeState  string   `json:"changeState"`
	Backlight    string   `json:"backlight"`
	Enabled      *bool    `json:"enabled,omitempty"`
	AlarmMessage *string  `json:"alarmMessage,omitempty"`
}

// ActivityTracking represents activity tracking settings.
type ActivityTracking struct {
	ActivityTrackingEnabled     *bool `json:"activityTrackingEnabled,omitempty"`
	MoveAlertEnabled            *bool `json:"moveAlertEnabled,omitempty"`
	MoveBarEnabled              *bool `json:"moveBarEnabled,omitempty"`
	PulseOxSleepTrackingEnabled *bool `json:"pulseOxSleepTrackingEnabled,omitempty"`
	PulseOxAcclimationEnabled   *bool `json:"pulseOxAcclimationEnabled,omitempty"`
	HighHrAlertEnabled          *bool `json:"highHrAlertEnabled,omitempty"`
	HighHrAlertThreshold        *int  `json:"highHrAlertThreshold,omitempty"`
	LowHrAlertEnabled           *bool `json:"lowHrAlertEnabled,omitempty"`
	LowHrAlertThreshold         *int  `json:"lowHrAlertThreshold,omitempty"`
}

// SupportedLanguage represents a supported language on the device.
type SupportedLanguage struct {
	ID   int64  `json:"id"`
	Name string `json:"name"`
}

// DeviceSettings represents the settings for a device.
type DeviceSettings struct {
	DeviceID         int64  `json:"deviceId"`
	TimeFormat       string `json:"timeFormat"`
	DateFormat       string `json:"dateFormat"`
	MeasurementUnits string `json:"measurementUnits"`
	AllUnits         string `json:"allUnits"`
	Language         int    `json:"language"`
	StartOfWeek      string `json:"startOfWeek"`

	// Alarms
	Alarms               []DeviceAlarm `json:"alarms"`
	MultipleAlarmEnabled bool          `json:"multipleAlarmEnabled"`

	// Activity tracking
	ActivityTracking *ActivityTracking `json:"activityTracking,omitempty"`

	// Supported configurations
	SupportedLanguages           []SupportedLanguage `json:"supportedLanguages"`
	SupportedAudioPromptDialects []string            `json:"supportedAudioPromptDialects"`

	// Sound and alerts
	AlertTonesEnabled         *bool `json:"alertTonesEnabled,omitempty"`
	SoundVibrationEnabled     *bool `json:"soundVibrationEnabled,omitempty"`
	SoundInAppOnlyEnabled     *bool `json:"soundInAppOnlyEnabled,omitempty"`
	BleConnectionAlertEnabled *bool `json:"bleConnectionAlertEnabled,omitempty"`
	DndEnabled                *bool `json:"dndEnabled,omitempty"`

	// Live features
	LiveEventSharingEnabled     bool     `json:"liveEventSharingEnabled"`
	LiveTrackEnabled            bool     `json:"liveTrackEnabled"`
	LiveEventSharingMsgContents []string `json:"liveEventSharingMsgContents"`
	LiveEventSharingMsgTriggers []string `json:"liveEventSharingMsgTriggers"`

	// Intensity minutes
	IntensityMinutesCalcMethod     string `json:"intensityMinutesCalcMethod"`
	ModerateIntensityMinutesHrZone int    `json:"moderateIntensityMinutesHrZone"`
	VigorousIntensityMinutesHrZone int    `json:"vigorousIntensityMinutesHrZone"`

	// Audio prompts
	AudioPromptDialectType string `json:"audioPromptDialectType"`

	// Sync settings
	AutoSyncStepsBeforeSync   int `json:"autoSyncStepsBeforeSync"`
	AutoSyncMinutesBeforeSync int `json:"autoSyncMinutesBeforeSync"`

	raw json.RawMessage
}

// RawJSON returns the original JSON response.
func (s *DeviceSettings) RawJSON() json.RawMessage {
	return s.raw
}

// DeviceMessage represents a message for a device.
type DeviceMessage struct {
	ID          int64  `json:"id,omitempty"`
	MessageType string `json:"messageType,omitempty"`
	Content     string `json:"content,omitempty"`
}

// DeviceMessages represents the messages response.
type DeviceMessages struct {
	ServiceHost   string          `json:"serviceHost"`
	NumOfMessages int             `json:"numOfMessages"`
	Messages      []DeviceMessage `json:"messages"`

	raw json.RawMessage
}

// RawJSON returns the original JSON response.
func (m *DeviceMessages) RawJSON() json.RawMessage {
	return m.raw
}

// DeviceWeight represents a device's weight/priority info.
type DeviceWeight struct {
	DeviceID               int64  `json:"deviceId"`
	DisplayName            string `json:"displayName"`
	ImageURL               string `json:"imageUrl"`
	Weight                 int    `json:"weight"`
	PrimaryTrainingCapable bool   `json:"primaryTrainingCapable"`
	LhaBackupCapable       bool   `json:"lhaBackupCapable"`
	PrimaryWearableDevice  bool   `json:"primaryWearableDevice"`
}

// DeviceWeightList represents a list of device weights.
type DeviceWeightList struct {
	DeviceWeights              []DeviceWeight `json:"deviceWeights"`
	WearableDeviceCount        int            `json:"wearableDeviceCount,omitempty"`
	PrimaryTrainingDeviceCount int            `json:"primaryTrainingDeviceCount,omitempty"`
}

// PrimaryTrainingDeviceInfo represents the primary training device response.
type PrimaryTrainingDeviceInfo struct {
	PrimaryTrainingDevice struct {
		DeviceID int64 `json:"deviceId"`
	} `json:"PrimaryTrainingDevice"`
	WearableDevices           DeviceWeightList `json:"WearableDevices"`
	TrainingStatusOnlyDevices DeviceWeightList `json:"TrainingStatusOnlyDevices"`
	PrimaryTrainingDevices    DeviceWeightList `json:"PrimaryTrainingDevices"`
	RegisteredDevices         []Device         `json:"RegisteredDevices"`

	raw json.RawMessage
}

// RawJSON returns the original JSON response.
func (p *PrimaryTrainingDeviceInfo) RawJSON() json.RawMessage {
	return p.raw
}

// GetDevices retrieves the list of registered devices.
func (s *DeviceService) GetDevices(ctx context.Context) ([]Device, error) {
	path := "/device-service/deviceregistration/devices"

	resp, err := s.client.doAPI(ctx, http.MethodGet, path, http.NoBody)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusNoContent || resp.StatusCode == http.StatusNotFound {
		return nil, ErrNotFound
	}

	raw, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var devices []Device
	if err := json.Unmarshal(raw, &devices); err != nil {
		return nil, err
	}

	// Store raw JSON in each device
	for i := range devices {
		devices[i].raw = raw
	}

	return devices, nil
}

// GetSettings retrieves the settings for a specific device.
func (s *DeviceService) GetSettings(ctx context.Context, deviceID int64) (*DeviceSettings, error) {
	path := fmt.Sprintf("/device-service/deviceservice/device-info/settings/%d", deviceID)

	resp, err := s.client.doAPI(ctx, http.MethodGet, path, http.NoBody)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusNoContent || resp.StatusCode == http.StatusNotFound {
		return nil, ErrNotFound
	}

	raw, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var settings DeviceSettings
	if err := json.Unmarshal(raw, &settings); err != nil {
		return nil, err
	}
	settings.raw = raw

	return &settings, nil
}

// GetMessages retrieves device messages.
func (s *DeviceService) GetMessages(ctx context.Context) (*DeviceMessages, error) {
	path := "/device-service/devicemessage/messages"

	resp, err := s.client.doAPI(ctx, http.MethodGet, path, http.NoBody)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusNoContent || resp.StatusCode == http.StatusNotFound {
		return nil, ErrNotFound
	}

	raw, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var messages DeviceMessages
	if err := json.Unmarshal(raw, &messages); err != nil {
		return nil, err
	}
	messages.raw = raw

	return &messages, nil
}

// GetPrimaryTrainingDevice retrieves info about the primary training device.
func (s *DeviceService) GetPrimaryTrainingDevice(ctx context.Context) (*PrimaryTrainingDeviceInfo, error) {
	path := "/web-gateway/device-info/primary-training-device"

	resp, err := s.client.doAPI(ctx, http.MethodGet, path, http.NoBody)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusNoContent || resp.StatusCode == http.StatusNotFound {
		return nil, ErrNotFound
	}

	raw, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var info PrimaryTrainingDeviceInfo
	if err := json.Unmarshal(raw, &info); err != nil {
		return nil, err
	}
	info.raw = raw

	return &info, nil
}
