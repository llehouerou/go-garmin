package garmin

import (
	"encoding/json"
	"testing"
)

func TestLactateThresholdEntryJSONUnmarshal(t *testing.T) {
	// Speed entry
	speedJSON := `{"userProfilePK":12345678,"version":1769455333539,"calendarDate":"2026-01-26T23:22:13.450","sequence":1769455333539,"speed":0.33611017,"hearRate":null,"heartRateCycling":null}`

	var speedEntry LactateThresholdEntry
	if err := json.Unmarshal([]byte(speedJSON), &speedEntry); err != nil {
		t.Fatalf("Failed to unmarshal speed entry: %v", err)
	}

	if speedEntry.UserProfilePK != 12345678 {
		t.Errorf("UserProfilePK = %d, want 12345678", speedEntry.UserProfilePK)
	}
	if speedEntry.Speed == nil || *speedEntry.Speed < 0.336 || *speedEntry.Speed > 0.337 {
		t.Errorf("Speed = %v, want ~0.33611017", speedEntry.Speed)
	}
	if speedEntry.HeartRate != nil {
		t.Errorf("HeartRate = %v, want nil", speedEntry.HeartRate)
	}

	// HR entry
	hrJSON := `{"userProfilePK":12345678,"version":1769455333539,"calendarDate":"2026-01-26T23:22:13.450","sequence":1769455333539,"speed":null,"hearRate":166,"heartRateCycling":null}`

	var hrEntry LactateThresholdEntry
	if err := json.Unmarshal([]byte(hrJSON), &hrEntry); err != nil {
		t.Fatalf("Failed to unmarshal HR entry: %v", err)
	}

	if hrEntry.Speed != nil {
		t.Errorf("Speed = %v, want nil", hrEntry.Speed)
	}
	if hrEntry.HeartRate == nil || *hrEntry.HeartRate != 166 {
		t.Errorf("HeartRate = %v, want 166", hrEntry.HeartRate)
	}
}

func TestLactateThresholdHelpers(t *testing.T) {
	speed := 0.33611017
	hr := 166

	lt := &LactateThreshold{
		Entries: []LactateThresholdEntry{
			{Speed: &speed, HeartRate: nil},
			{Speed: nil, HeartRate: &hr},
		},
	}

	if s := lt.Speed(); s == nil || *s != speed {
		t.Errorf("Speed() = %v, want %v", s, speed)
	}
	if h := lt.HeartRate(); h == nil || *h != hr {
		t.Errorf("HeartRate() = %v, want %v", h, hr)
	}
}

func TestLactateThresholdRawJSON(t *testing.T) {
	rawJSON := `[{"userProfilePK":123,"speed":0.5}]`
	lt := &LactateThreshold{raw: json.RawMessage(rawJSON)}

	if string(lt.RawJSON()) != rawJSON {
		t.Error("RawJSON should return original JSON")
	}
}

func TestFunctionalThresholdPowerJSONUnmarshal(t *testing.T) {
	// With data
	withDataJSON := `{"userProfilePK":12345678,"version":123,"calendarDate":"2026-01-27","isStale":false,"sequence":456,"sport":"CYCLING","functionalThresholdPower":280,"biometricSourceType":"AUTO"}`

	var ftp FunctionalThresholdPower
	if err := json.Unmarshal([]byte(withDataJSON), &ftp); err != nil {
		t.Fatalf("Failed to unmarshal: %v", err)
	}

	if ftp.UserProfilePK != 12345678 {
		t.Errorf("UserProfilePK = %d, want 12345678", ftp.UserProfilePK)
	}
	if ftp.FunctionalThresholdPower == nil || *ftp.FunctionalThresholdPower != 280 {
		t.Errorf("FunctionalThresholdPower = %v, want 280", ftp.FunctionalThresholdPower)
	}

	// With null values
	nullJSON := `{"userProfilePK":12345678,"version":null,"calendarDate":null,"isStale":null,"sequence":null,"sport":null,"functionalThresholdPower":null,"biometricSourceType":null}`

	var ftpNull FunctionalThresholdPower
	if err := json.Unmarshal([]byte(nullJSON), &ftpNull); err != nil {
		t.Fatalf("Failed to unmarshal null values: %v", err)
	}

	if ftpNull.FunctionalThresholdPower != nil {
		t.Errorf("FunctionalThresholdPower = %v, want nil", ftpNull.FunctionalThresholdPower)
	}
}

func TestFunctionalThresholdPowerRawJSON(t *testing.T) {
	rawJSON := `{"userProfilePK":123}`
	ftp := &FunctionalThresholdPower{raw: json.RawMessage(rawJSON)}

	if string(ftp.RawJSON()) != rawJSON {
		t.Error("RawJSON should return original JSON")
	}
}

func TestPowerToWeightJSONUnmarshal(t *testing.T) {
	rawJSON := `{"userProfilePk":12345678,"calendarDate":"2026-01-27T10:06:56.496","origin":"weight","sport":"RUNNING","functionalThresholdPower":347,"weight":74.2,"powerToWeight":4.67654986522911,"ftpCreateTime":"2025-11-15T10:54:08.91","weightCreateTime":"2026-01-27T10:06:56.496","isStale":false}`

	var ptw PowerToWeight
	if err := json.Unmarshal([]byte(rawJSON), &ptw); err != nil {
		t.Fatalf("Failed to unmarshal: %v", err)
	}

	if ptw.UserProfilePK != 12345678 {
		t.Errorf("UserProfilePK = %d, want 12345678", ptw.UserProfilePK)
	}
	if ptw.FunctionalThresholdPower != 347 {
		t.Errorf("FunctionalThresholdPower = %d, want 347", ptw.FunctionalThresholdPower)
	}
	if ptw.Weight < 74.1 || ptw.Weight > 74.3 {
		t.Errorf("Weight = %f, want ~74.2", ptw.Weight)
	}
	if ptw.PowerToWeightRatio < 4.67 || ptw.PowerToWeightRatio > 4.68 {
		t.Errorf("PowerToWeightRatio = %f, want ~4.676", ptw.PowerToWeightRatio)
	}
	if ptw.Sport != "RUNNING" {
		t.Errorf("Sport = %s, want RUNNING", ptw.Sport)
	}
	if ptw.IsStale {
		t.Error("IsStale = true, want false")
	}
}

func TestPowerToWeightRawJSON(t *testing.T) {
	rawJSON := `{"userProfilePk":123}`
	ptw := &PowerToWeight{raw: json.RawMessage(rawJSON)}

	if string(ptw.RawJSON()) != rawJSON {
		t.Error("RawJSON should return original JSON")
	}
}

func TestBiometricStatJSONUnmarshal(t *testing.T) {
	const statDate = "2026-01-26"
	const seriesRunning = "running"
	rawJSON := `{"from":"2026-01-26","until":"2026-01-26","series":"running","value":0.33611017,"updatedDate":"2026-01-26"}`

	var stat BiometricStat
	if err := json.Unmarshal([]byte(rawJSON), &stat); err != nil {
		t.Fatalf("Failed to unmarshal: %v", err)
	}

	if stat.From != statDate {
		t.Errorf("From = %s, want %s", stat.From, statDate)
	}
	if stat.Until != statDate {
		t.Errorf("Until = %s, want %s", stat.Until, statDate)
	}
	if stat.Series != seriesRunning {
		t.Errorf("Series = %s, want %s", stat.Series, seriesRunning)
	}
	if stat.Value < 0.336 || stat.Value > 0.337 {
		t.Errorf("Value = %f, want ~0.33611017", stat.Value)
	}
}

func TestBiometricStatsRawJSON(t *testing.T) {
	rawJSON := `[{"from":"2026-01-26","until":"2026-01-26","series":"running","value":166.0,"updatedDate":"2026-01-26"}]`
	bs := &BiometricStats{raw: json.RawMessage(rawJSON)}

	if string(bs.RawJSON()) != rawJSON {
		t.Error("RawJSON should return original JSON")
	}
}

func TestHeartRateZoneJSONUnmarshal(t *testing.T) {
	rawJSON := `{"trainingMethod":"HR_RESERVE","restingHeartRateUsed":54,"lactateThresholdHeartRateUsed":166,"zone1Floor":121,"zone2Floor":134,"zone3Floor":148,"zone4Floor":161,"zone5Floor":175,"maxHeartRateUsed":188,"restingHrAutoUpdateUsed":true,"sport":"DEFAULT","changeState":"UNCHANGED"}`

	var zone HeartRateZone
	if err := json.Unmarshal([]byte(rawJSON), &zone); err != nil {
		t.Fatalf("Failed to unmarshal: %v", err)
	}

	if zone.TrainingMethod != "HR_RESERVE" {
		t.Errorf("TrainingMethod = %s, want HR_RESERVE", zone.TrainingMethod)
	}
	if zone.RestingHeartRateUsed != 54 {
		t.Errorf("RestingHeartRateUsed = %d, want 54", zone.RestingHeartRateUsed)
	}
	if zone.LactateThresholdHeartRateUsed != 166 {
		t.Errorf("LactateThresholdHeartRateUsed = %d, want 166", zone.LactateThresholdHeartRateUsed)
	}
	if zone.Zone1Floor != 121 {
		t.Errorf("Zone1Floor = %d, want 121", zone.Zone1Floor)
	}
	if zone.Zone5Floor != 175 {
		t.Errorf("Zone5Floor = %d, want 175", zone.Zone5Floor)
	}
	if zone.MaxHeartRateUsed != 188 {
		t.Errorf("MaxHeartRateUsed = %d, want 188", zone.MaxHeartRateUsed)
	}
	if !zone.RestingHrAutoUpdateUsed {
		t.Error("RestingHrAutoUpdateUsed = false, want true")
	}
	if zone.Sport != "DEFAULT" {
		t.Errorf("Sport = %s, want DEFAULT", zone.Sport)
	}
	if zone.ChangeState != "UNCHANGED" {
		t.Errorf("ChangeState = %s, want UNCHANGED", zone.ChangeState)
	}
}

func TestHeartRateZonesRawJSON(t *testing.T) {
	rawJSON := `[{"trainingMethod":"HR_RESERVE","sport":"DEFAULT"}]`
	hrz := &HeartRateZones{raw: json.RawMessage(rawJSON)}

	if string(hrz.RawJSON()) != rawJSON {
		t.Error("RawJSON should return original JSON")
	}
}
