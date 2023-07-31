package models

import (
	"encoding/json"
	"fmt"
	"os"
)

type Config struct {
	DefaultFolder   string   `json:"default_folder"`
	DefaultPrefix   string   `json:"default_prefix"`
	TrackedPrefixes []Prefix `json:"tracked_prefixes"`
}

func ConfigFilePath() string {
	return os.Getenv("HOME") + "/.config/wineprefixer.json"
}

func CreateConfigFileIfNotExist() error {
	// Verify file status
	_, err := os.Stat(ConfigFilePath())
	if os.IsNotExist(err) {
		// Create empty config file
		config := Config{
			DefaultFolder:   os.Getenv("HOME"),
			DefaultPrefix:   "",
			TrackedPrefixes: []Prefix{},
		}

		// create file
		file, err := os.Create(ConfigFilePath())
		if err != nil {
			return fmt.Errorf("error while creating the file: %w", err)
		}
		defer file.Close()

		// convert to json
		configJSON, err := json.MarshalIndent(config, "", "	")
		if err != nil {
			return fmt.Errorf("error while converting to JSON: %w", err)
		}

		// Write File
		_, err = file.Write(configJSON)
		if err != nil {
			return fmt.Errorf("error while writing the config file: %w", err)
		}
	} else if err != nil {
		return fmt.Errorf("error checking for existence of configuration file: %w", err)
	}
	return nil
}

func ReadConfigFromFile() (Config, error) {
	// reading file
	fileBytes, err := os.ReadFile(ConfigFilePath())
	if err != nil {
		return Config{}, fmt.Errorf("error reading configuration file: %w", err)
	}

	// convert from json
	var config Config
	err = json.Unmarshal(fileBytes, &config)
	if err != nil {
		return Config{}, fmt.Errorf("error while converting JSON: %w", err)
	}

	return config, nil
}

// WriteConfigFile
func (config *Config) SaveConfigToFile() error {
	// convert to json
	configJSON, err := json.MarshalIndent(config, "", "	")
	if err != nil {
		return fmt.Errorf("error while converting to JSON: %w", err)
	}

	// Write File
	err = os.WriteFile(ConfigFilePath(), configJSON, 0644)
	if err != nil {
		return fmt.Errorf("error while saving the config file: %w", err)
	}

	return nil
}

// TrackPrefixIfUnique
func (config *Config) TrackPrefix(newPrefix Prefix) error {
	isUniqueID, isUniquePath := true, true
	for _, p := range config.TrackedPrefixes {
		isUniqueID = isUniqueID && p.ID != newPrefix.ID
		isUniquePath = isUniquePath && p.Path != newPrefix.Path
	}

	if !isUniqueID {
		return fmt.Errorf("the id: %s is already in use", newPrefix.ID)
	}
	if !isUniquePath {
		return fmt.Errorf("the path: %s is already in use", newPrefix.Path)
	}

	config.TrackedPrefixes = append(config.TrackedPrefixes, newPrefix)

	return nil
}

// SelectPrefix
func (config *Config) GetPrefix(id string) (Prefix, error) {
	var prefix Prefix
	isValidID := false
	for _, p := range config.TrackedPrefixes {
		if p.ID == id {
			isValidID = true
			prefix = p
			break
		}
	}

	if !isValidID {
		return prefix, fmt.Errorf("id does not correspond to a tracked prefix")
	}

	return prefix, nil
}

// RemoveTrackedPrefix
func (config *Config) RemovePrefix(id string) (Prefix, error) {
	index := -1
	var delPrefix Prefix
	for i, p := range config.TrackedPrefixes {
		if p.ID == id {
			delPrefix = p
			index = i
			break
		}
	}

	if index == -1 {
		return Prefix{}, fmt.Errorf("id %s does not correspond to any tracked prefix", id)
	}

	config.TrackedPrefixes = append(config.TrackedPrefixes[:index], config.TrackedPrefixes[index+1:]...)

	return delPrefix, nil
}
