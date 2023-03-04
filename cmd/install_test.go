package cmd

import (
	"bytes"
	"github.com/Sithell/fox/internal"
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
}
