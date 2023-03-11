package internal

import (
	"github.com/spf13/afero"
	"gopkg.in/yaml.v3"
	"reflect"
	"testing"
)

const firefoxPath = "/home/sithell/.mozilla/firefox"
const profileName = "87syjf8o.default-release"
const fullProfilePath = firefoxPath + "/" + profileName
const fullChromePath = firefoxPath + "/" + profileName + "/chrome"

func TestPrepareFirefox(t *testing.T) {
	fs := InitFs(firefoxPath, profileName)
	result, err := PrepareFirefox(fs)
	if err != nil {
		t.Errorf("PrepareFirefox returned error: %s", err)
	}
	if result != fullChromePath {
		t.Errorf("Invalid chrome directory path")
	}
	if _, err = fs.Stat(fullChromePath); err != nil {
		t.Errorf("Failed to locate chrome directory: %s", err)
	}
}

func TestLocateProfileDir(t *testing.T) {
	fs := InitFs(firefoxPath, profileName)
	result, err := locateProfileDir(fs)
	if result != fullProfilePath {
		t.Errorf("locateProfileDir returned %s, %s", result, err)
	}
}

func TestInitChrome(t *testing.T) {
	fs := InitFs(firefoxPath, profileName)
	err := initChrome(fs, fullProfilePath)
	if err != nil {
		t.Errorf("initChrome returned error: %s", err)
	}
	_, err = fs.Stat(fullChromePath)
	if err != nil {
		t.Errorf("%s", err)
	}
}

func TestLocateUserFiles(t *testing.T) {
	fs := afero.NewMemMapFs()
	for _, filename := range []string{
		"/modName/userContent.css",
		"/modName/src/userChrome.css",
	} {
		if _, err := fs.Create(filename); err != nil {
			panic(err)
		}
	}

	result := LocateUserFiles(fs, "/modName")
	if !reflect.DeepEqual(result, map[UserFile]string{
		"userChrome.css":  "src/userChrome.css",
		"userContent.css": "userContent.css",
	}) {
		t.Errorf("LocateUserFiles returned invalid result: %v", result)
	}
}

func TestLoadFoxYml(t *testing.T) {
	expected := FoxYml{[]FoxYmlMod{
		{
			"Mono-firefox-theme",
			"https://github.com/witalihirsch/Mono-firefox-theme.git",
		},
	}}
	content, _ := yaml.Marshal(expected)

	fs := afero.NewMemMapFs()
	file, err := fs.Create("fox.yml")
	if err != nil {
		panic(err)
	}
	if _, err = file.Write(content); err != nil {
		panic(err)
	}

	result, err := LoadFoxYml(fs, "fox.yml")
	if err != nil {
		t.Errorf("LoadFoxYml returned error: %s", err)
	}
	if !reflect.DeepEqual(result, expected) {
		t.Errorf("Expected %v, got %v", expected, result)
	}
}
