package timingsmanager

import (
	"encoding/json"
	"os"
	"path/filepath"
	"reflect"
	"time"
)

type Timings struct {
	SeedShop int `json:"seed_shop"`
	GearShop int `json:"gear_shop"`
}

func getTimingsPath() (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}
	dir := filepath.Join(home, "Library", "Application Support", "Taz GAG Macro")
	os.MkdirAll(dir, 0755)
	return filepath.Join(dir, "timings.json"), nil
}

func OnCooldown(objective string, cooldown time.Duration) bool {
	path, err := getTimingsPath()
	if err != nil {
		return false
	}

	data := &Timings{}
	if _, err := os.Stat(path); err == nil {
		fileBytes, err := os.ReadFile(path)
		if err != nil {
			return false
		}
		if err := json.Unmarshal(fileBytes, data); err != nil {
			return false
		}
	}

	v := reflect.ValueOf(data).Elem().FieldByName(objective)
	if !v.IsValid() {
		return false
	}

	lastTime := v.Int()
	return time.Since(time.Unix(lastTime, 0)) < cooldown
}

func UpdateObjectiveTime(objective string) error {
	path, err := getTimingsPath()
	if err != nil {
		return err
	}

	data := &Timings{}
	if _, err := os.Stat(path); err == nil {
		fileBytes, err := os.ReadFile(path)
		if err == nil {
			json.Unmarshal(fileBytes, data)
		}
	}

	v := reflect.ValueOf(data).Elem().FieldByName(objective)
	if !v.IsValid() || !v.CanSet() {
		return nil
	}

	v.SetInt(time.Now().Unix())

	fileBytes, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		return err
	}

	return os.WriteFile(path, fileBytes, 0644)
}
