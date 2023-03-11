package cmd

import (
	"bytes"
	"github.com/Sithell/fox/internal"
	"github.com/spf13/afero"
	"gopkg.in/yaml.v3"
	"io"
	"strings"
	"testing"
)

const firefoxPath = "/home/sithell/.mozilla/firefox"
const profileName = "87syjf8o.default-release"

func TestNewInstallCmd(t *testing.T) {
	fs := internal.InitFs(firefoxPath, profileName)
	cmd := NewInstallCmd(fs)
	b := bytes.NewBufferString("")
	cmd.SetOut(b)
	cmd.SetArgs([]string{"https://github.com/witalihirsch/Mono-firefox-theme.git"})
	err := cmd.Execute()
	if err != nil {
		panic(err)
	}
	out, err := io.ReadAll(b)
	if err != nil {
		panic(err)
	}
	if strings.Split(string(out), "\n")[0] != "Installing Mono-firefox-theme to "+firefoxPath+"/"+profileName+"/chrome" {
		t.Errorf("Invalid output: %s", string(out))
	}
	foxYml, err := afero.ReadFile(fs, firefoxPath+"/"+profileName+"/"+"/chrome/fox.yml")
	if err != nil {
		t.Errorf("Failed to open fox.yml: %s", err)
	}

	expected := internal.FoxYml{[]internal.FoxYmlMod{
		{
			"Mono-firefox-theme",
			"https://github.com/witalihirsch/Mono-firefox-theme.git",
		},
	}}
	content, _ := yaml.Marshal(expected)
	if string(foxYml) != string(content) {
		t.Errorf("Contents of fox.yml do not match: %s", string(foxYml))
	}
}
