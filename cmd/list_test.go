package cmd

import (
	"bytes"
	"github.com/Sithell/fox/internal"
	"github.com/spf13/afero"
	"io"
	"testing"
)

func TestNewListCmd(t *testing.T) {
	// arrange
	fs := internal.InitFs(firefoxPath, profileName)
	chromeDir := firefoxPath + "/" + profileName + "/chrome"
	err := afero.WriteFile(fs, chromeDir+"/fox.yml", []byte(`mods:
    - name: Mono-firefox-theme
      url: https://github.com/witalihirsch/Mono-firefox-theme.git
`), 0644)
	if err != nil {
		panic(err)
	}

	// execute
	cmd := NewListCmd(fs)
	b := bytes.NewBufferString("")
	cmd.SetOut(b)
	cmd.SetArgs([]string{})
	err = cmd.Execute()
	if err != nil {
		panic(err)
	}
	out, err := io.ReadAll(b)
	if err != nil {
		panic(err)
	}

	// assert
	if string(out) != `1. Mono-firefox-theme - https://github.com/witalihirsch/Mono-firefox-theme.git` {
		t.Errorf("Invalid output: %s", string(out))
	}
}
