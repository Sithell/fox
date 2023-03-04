package internal

import (
	"errors"
	"github.com/spf13/afero"
	"gopkg.in/ini.v1"
	"os"
	"strings"
)

func PrepareFirefox(fs afero.Fs) (string, error) {
	profileDir, err := locateProfileDir(fs)
	if err != nil {
		return "", err
	}

	err = initChrome(fs, profileDir)
	if err != nil {
		return "", err
	}

	return profileDir + "/chrome", nil
}

func locateProfileDir(fs afero.Fs) (string, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}
	var configDir string
	if os.Getenv("ENV") == "dev" {
		configDir = "tmp"
	} else {
		configDir = homeDir + "/.mozilla/firefox"
	}

	f, err := afero.ReadFile(fs, configDir+"/profiles.ini")
	cfg, err := ini.Load(f)

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
	if _, err := fs.Stat(profileDir); err != nil {
		return "", err
	}
	return profileDir, nil
}

func initChrome(fs afero.Fs, profileDir string) error {
	chromeDir := profileDir + "/chrome"

	if _, err := fs.Stat(chromeDir); os.IsNotExist(err) {
		if err := fs.Mkdir(chromeDir, os.ModePerm); err != nil {
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

func LocateUserFiles(fs afero.Fs, modPath string) (userFilePaths map[UserFile]string) {
	userFilePaths = make(map[UserFile]string)
	for _, userFile := range []UserFile{Chrome, Content} {
		minimalDepth := 999999
		var userFilePath string
		err := afero.Walk(fs, modPath, func(path string, info os.FileInfo, err error) error {
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
