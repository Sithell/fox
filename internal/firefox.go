package internal

import (
	"errors"
	"gopkg.in/ini.v1"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
)

func PrepareFirefox() (string, error) {
	profileDir, err := locateProfileDir()
	if err != nil {
		return "", err
	}

	err = initChrome(profileDir)
	if err != nil {
		return "", err
	}

	return profileDir + "/chrome", nil
}

func locateProfileDir() (string, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}
	var configDir = homeDir + "/.mozilla/firefox"
	configDir = "/home/sithell/dev/fox/tmp"
	cfg, err := ini.Load(configDir + "/profiles.ini")

	profiles := cfg.Sections()
	var profileName string
	for _, profile := range profiles {
		if profile.HasKey("Locked") && profile.Key("Locked").String() == "1" {
			profileName = profile.Key("Default").String()
			break
		}
	}
	if profileName == "" {
		return "", errors.New("No active profile found")
	}

	var profileDir = configDir + "/" + profileName
	if _, err := os.Stat(profileDir); err != nil {
		return "", err
	}
	return profileDir, nil
}

func initChrome(profileDir string) error {
	chromeDir := profileDir + "/chrome"

	if _, err := os.Stat(chromeDir); os.IsNotExist(err) {
		if err := os.Mkdir(chromeDir, os.ModePerm); err != nil {
			return err
		}
	} else if err != nil {
		return err
	}
	return nil
}

type UserFile string

const (
	Chrome  = "userChrome.css"
	Content = "userContent.css"
)

func LocateUserFiles(modPath string) (userFilePaths map[UserFile]string) {
	userFilePaths = make(map[UserFile]string)
	for _, userFile := range []UserFile{Chrome, Content} {
		minimalDepth := 999999
		var userFilePath string
		err := filepath.Walk(modPath, func(path string, info fs.FileInfo, err error) error {
			depth := len(strings.Split(path, "/"))
			if strings.Contains(path, string(userFile)) && depth < minimalDepth {
				minimalDepth = depth
				userFilePath = path
			}
			return nil
		})
		if err != nil {
			panic(err)
		}
		if userFilePath != "" {
			userFilePaths[userFile] = userFilePath[len(modPath)+1:]
		}
	}
	return userFilePaths
}
