// service_device_test.go
package garmin

import (
	"encoding/json"
	"testing"
)

const (
	testAnonymousName  = "anonymous"
	testMetricSystem   = "metric"
	testPrivateTypeKey = "private"
)

func TestDeviceJSONUnmarshal(t *testing.T) {
	rawJSON := `{
		"deviceId": 12345678,
		"unitId": 12345678,
		"displayName": "anonymous",
		"productDisplayName": "Forerunner 965",
		"serialNumber": "ABC123456",
		"deviceStatus": "active",
		"partNumber": "006-B4315-00",
		"productSku": "010-02809-10",
		"applicationKey": "forerunner965",
		"imageUrl": "https://res.garmin.com/en/products/010-02809-10/v/c1_01_md.png",
		"registeredDate": 1737031539000,
		"currentFirmwareVersion": "27.09",
		"currentFirmwareVersionMajor": 27,
		"currentFirmwareVersionMinor": 9,
		"deviceTypePk": 37086,
		"deviceTypeSimpleName": "Garmin Forerunner 965",
		"deviceCategories": ["FITNESS", "WELLNESS", "GOLF", "OUTDOOR"],
		"wifi": true,
		"bluetoothClassicDevice": false,
		"bluetoothLowEnergyDevice": true,
		"primary": true,
		"activeInd": 1,
		"primaryActivityTrackerIndicator": true,
		"bodyBatteryCapable": true,
		"hrvStatusCapable": true,
		"trainingReadinessCapable": true,
		"vo2MaxRunCapable": true,
		"vo2MaxBikeCapable": true,
		"pulseOxAllDayCapable": true,
		"pulseOxSleepCapable": true,
		"sleepScoreCapable": true,
		"allDayStressCapable": true,
		"respirationCapable": true,
		"hasOpticalHeartRate": true
	}`

	var device Device
	if err := json.Unmarshal([]byte(rawJSON), &device); err != nil {
		t.Fatalf("Failed to unmarshal: %v", err)
	}

	if device.DeviceID != 12345678 {
		t.Errorf("DeviceID = %d, want 12345678", device.DeviceID)
	}
	if device.UnitID != 12345678 {
		t.Errorf("UnitID = %d, want 12345678", device.UnitID)
	}
	if device.DisplayName != testAnonymousName {
		t.Errorf("DisplayName = %s, want %s", device.DisplayName, testAnonymousName)
	}
	if device.ProductDisplayName != "Forerunner 965" {
		t.Errorf("ProductDisplayName = %s, want Forerunner 965", device.ProductDisplayName)
	}
	if device.SerialNumber != "ABC123456" {
		t.Errorf("SerialNumber = %s, want ABC123456", device.SerialNumber)
	}
	if device.DeviceStatus != "active" {
		t.Errorf("DeviceStatus = %s, want active", device.DeviceStatus)
	}
	if !device.WiFi {
		t.Error("expected WiFi to be true")
	}
	if device.BluetoothClassicDevice {
		t.Error("expected BluetoothClassicDevice to be false")
	}
	if !device.BluetoothLowEnergyDevice {
		t.Error("expected BluetoothLowEnergyDevice to be true")
	}
	if !device.Primary {
		t.Error("expected Primary to be true")
	}
	if device.ActiveInd != 1 {
		t.Errorf("ActiveInd = %d, want 1", device.ActiveInd)
	}
	if !device.BodyBatteryCapable {
		t.Error("expected BodyBatteryCapable to be true")
	}
	if !device.HRVStatusCapable {
		t.Error("expected HRVStatusCapable to be true")
	}
	if !device.TrainingReadinessCapable {
		t.Error("expected TrainingReadinessCapable to be true")
	}
	if len(device.DeviceCategories) != 4 {
		t.Errorf("DeviceCategories length = %d, want 4", len(device.DeviceCategories))
	}
	if device.CurrentFirmwareVersion != "27.09" {
		t.Errorf("CurrentFirmwareVersion = %s, want 27.09", device.CurrentFirmwareVersion)
	}
}

func TestDeviceRawJSON(t *testing.T) {
	rawJSON := `{"deviceId":12345678,"displayName":"anonymous"}`

	var device Device
	if err := json.Unmarshal([]byte(rawJSON), &device); err != nil {
		t.Fatal(err)
	}
	device.raw = json.RawMessage(rawJSON)

	if string(device.RawJSON()) != rawJSON {
		t.Error("RawJSON should return original JSON")
	}
}

func TestDeviceSettingsJSONUnmarshal(t *testing.T) {
	rawJSON := `{
		"deviceId": 12345678,
		"timeFormat": "time_twenty_four_hr",
		"dateFormat": "date_day_month",
		"measurementUnits": "metric",
		"allUnits": "metric",
		"language": 1,
		"startOfWeek": "MONDAY",
		"alarms": [
			{
				"alarmId": 1745847607,
				"alarmMode": "OFF",
				"alarmTime": 1320,
				"alarmDays": ["ONCE"],
				"alarmSound": "TONE_AND_VIBRATION",
				"changeState": "UNCHANGED",
				"backlight": "ON"
			}
		],
		"multipleAlarmEnabled": true,
		"supportedLanguages": [
			{"id": 12345678, "name": "en"},
			{"id": 12345678, "name": "fr"}
		],
		"supportedAudioPromptDialects": ["EN_US", "FR_FR"],
		"alertTonesEnabled": false,
		"soundVibrationEnabled": true,
		"soundInAppOnlyEnabled": false,
		"bleConnectionAlertEnabled": true,
		"dndEnabled": true,
		"liveEventSharingEnabled": false,
		"liveTrackEnabled": false,
		"intensityMinutesCalcMethod": "AUTO",
		"moderateIntensityMinutesHrZone": 3,
		"vigorousIntensityMinutesHrZone": 4,
		"audioPromptDialectType": "FR_FR",
		"autoSyncStepsBeforeSync": 2000,
		"autoSyncMinutesBeforeSync": 240
	}`

	var settings DeviceSettings
	if err := json.Unmarshal([]byte(rawJSON), &settings); err != nil {
		t.Fatalf("Failed to unmarshal: %v", err)
	}

	if settings.DeviceID != 12345678 {
		t.Errorf("DeviceID = %d, want 12345678", settings.DeviceID)
	}
	if settings.TimeFormat != "time_twenty_four_hr" {
		t.Errorf("TimeFormat = %s, want time_twenty_four_hr", settings.TimeFormat)
	}
	if settings.DateFormat != "date_day_month" {
		t.Errorf("DateFormat = %s, want date_day_month", settings.DateFormat)
	}
	if settings.MeasurementUnits != testMetricSystem {
		t.Errorf("MeasurementUnits = %s, want %s", settings.MeasurementUnits, testMetricSystem)
	}
	if settings.StartOfWeek != "MONDAY" {
		t.Errorf("StartOfWeek = %s, want MONDAY", settings.StartOfWeek)
	}
	if len(settings.Alarms) != 1 {
		t.Fatalf("Alarms length = %d, want 1", len(settings.Alarms))
	}
	if settings.Alarms[0].AlarmMode != "OFF" {
		t.Errorf("Alarms[0].AlarmMode = %s, want OFF", settings.Alarms[0].AlarmMode)
	}
	if !settings.MultipleAlarmEnabled {
		t.Error("expected MultipleAlarmEnabled to be true")
	}
	if len(settings.SupportedLanguages) != 2 {
		t.Errorf("SupportedLanguages length = %d, want 2", len(settings.SupportedLanguages))
	}
	if settings.SupportedLanguages[0].Name != "en" {
		t.Errorf("SupportedLanguages[0].Name = %s, want en", settings.SupportedLanguages[0].Name)
	}
	if settings.IntensityMinutesCalcMethod != "AUTO" {
		t.Errorf("IntensityMinutesCalcMethod = %s, want AUTO", settings.IntensityMinutesCalcMethod)
	}
	if settings.AudioPromptDialectType != "FR_FR" {
		t.Errorf("AudioPromptDialectType = %s, want FR_FR", settings.AudioPromptDialectType)
	}
}

func TestDeviceSettingsRawJSON(t *testing.T) {
	rawJSON := `{"deviceId":12345678,"timeFormat":"time_twenty_four_hr"}`

	var settings DeviceSettings
	if err := json.Unmarshal([]byte(rawJSON), &settings); err != nil {
		t.Fatal(err)
	}
	settings.raw = json.RawMessage(rawJSON)

	if string(settings.RawJSON()) != rawJSON {
		t.Error("RawJSON should return original JSON")
	}
}

func TestDeviceMessagesJSONUnmarshal(t *testing.T) {
	rawJSON := `{
		"serviceHost": "http://connectapi.garmin.com/",
		"numOfMessages": 0,
		"messages": []
	}`

	var messages DeviceMessages
	if err := json.Unmarshal([]byte(rawJSON), &messages); err != nil {
		t.Fatalf("Failed to unmarshal: %v", err)
	}

	if messages.ServiceHost != "http://connectapi.garmin.com/" {
		t.Errorf("ServiceHost = %s, want http://connectapi.garmin.com/", messages.ServiceHost)
	}
	if messages.NumOfMessages != 0 {
		t.Errorf("NumOfMessages = %d, want 0", messages.NumOfMessages)
	}
	if len(messages.Messages) != 0 {
		t.Errorf("Messages length = %d, want 0", len(messages.Messages))
	}
}

func TestDeviceMessagesRawJSON(t *testing.T) {
	rawJSON := `{"serviceHost":"http://example.com/","numOfMessages":0,"messages":[]}`

	var messages DeviceMessages
	if err := json.Unmarshal([]byte(rawJSON), &messages); err != nil {
		t.Fatal(err)
	}
	messages.raw = json.RawMessage(rawJSON)

	if string(messages.RawJSON()) != rawJSON {
		t.Error("RawJSON should return original JSON")
	}
}

func TestPrimaryTrainingDeviceInfoJSONUnmarshal(t *testing.T) {
	rawJSON := `{
		"PrimaryTrainingDevice": {"deviceId": 12345678},
		"WearableDevices": {
			"deviceWeights": [
				{
					"deviceId": 12345678,
					"displayName": "anonymous",
					"imageUrl": "https://res.garmin.com/en/products/010-02809-10/v/c1_01_md.png",
					"weight": 3,
					"primaryTrainingCapable": true,
					"lhaBackupCapable": true,
					"primaryWearableDevice": true
				}
			],
			"wearableDeviceCount": 1
		},
		"TrainingStatusOnlyDevices": {"deviceWeights": []},
		"PrimaryTrainingDevices": {
			"deviceWeights": [
				{
					"deviceId": 12345678,
					"displayName": "anonymous",
					"imageUrl": "https://res.garmin.com/en/products/010-02809-10/v/c1_01_md.png",
					"weight": 3,
					"primaryTrainingCapable": true,
					"lhaBackupCapable": true,
					"primaryWearableDevice": true
				}
			],
			"primaryTrainingDeviceCount": 1
		},
		"RegisteredDevices": []
	}`

	var info PrimaryTrainingDeviceInfo
	if err := json.Unmarshal([]byte(rawJSON), &info); err != nil {
		t.Fatalf("Failed to unmarshal: %v", err)
	}

	if info.PrimaryTrainingDevice.DeviceID != 12345678 {
		t.Errorf("PrimaryTrainingDevice.DeviceID = %d, want 12345678", info.PrimaryTrainingDevice.DeviceID)
	}
	if len(info.WearableDevices.DeviceWeights) != 1 {
		t.Fatalf("WearableDevices.DeviceWeights length = %d, want 1", len(info.WearableDevices.DeviceWeights))
	}
	if info.WearableDevices.DeviceWeights[0].DisplayName != testAnonymousName {
		t.Errorf("WearableDevices.DeviceWeights[0].DisplayName = %s, want %s", info.WearableDevices.DeviceWeights[0].DisplayName, testAnonymousName)
	}
	if !info.WearableDevices.DeviceWeights[0].PrimaryTrainingCapable {
		t.Error("expected PrimaryTrainingCapable to be true")
	}
	if info.WearableDevices.WearableDeviceCount != 1 {
		t.Errorf("WearableDevices.WearableDeviceCount = %d, want 1", info.WearableDevices.WearableDeviceCount)
	}
	if info.PrimaryTrainingDevices.PrimaryTrainingDeviceCount != 1 {
		t.Errorf("PrimaryTrainingDevices.PrimaryTrainingDeviceCount = %d, want 1", info.PrimaryTrainingDevices.PrimaryTrainingDeviceCount)
	}
}

func TestPrimaryTrainingDeviceInfoRawJSON(t *testing.T) {
	rawJSON := `{"PrimaryTrainingDevice":{"deviceId":12345678}}`

	var info PrimaryTrainingDeviceInfo
	if err := json.Unmarshal([]byte(rawJSON), &info); err != nil {
		t.Fatal(err)
	}
	info.raw = json.RawMessage(rawJSON)

	if string(info.RawJSON()) != rawJSON {
		t.Error("RawJSON should return original JSON")
	}
}
