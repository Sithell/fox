package cmd

import (
	"bytes"
	"github.com/Sithell/fox/internal"
	"github.com/spf13/afero"
	"io"
	"os"
	"reflect"
	"strings"
	"testing"
)

func TestNewUpdateCmd(t *testing.T) {
	// arrange
	fs := internal.InitFs(firefoxPath, profileName)
	chromeDir := firefoxPath + "/" + profileName + "/chrome"
	if err := fs.Mkdir(chromeDir, os.ModePerm); err != nil {
		panic(err)
	}
	if err := fs.Mkdir(chromeDir+"/Mono-firefox-theme", os.ModePerm); err != nil {
		panic(err)
	}
	err := afero.WriteFile(fs, chromeDir+"/userChrome.css", []byte(`
/* Fox mod manager start */
@import "Mono-firefox-theme/firefox/userChrome.css";
/* Fox mod manager end */
`), 0644)
	if err != nil {
		panic(err)
	}
	err = afero.WriteFile(fs, chromeDir+"/userContent.css", []byte(`
/* Fox mod manager start */
@import "Mono-firefox-theme/firefox/userContent.css";
/* Fox mod manager end */
`), 0644)
	if err != nil {
		panic(err)
	}
	err = afero.WriteFile(fs, chromeDir+"/fox.yml", []byte(`mods:
    - name: Mono-firefox-theme
      url: https://github.com/witalihirsch/Mono-firefox-theme.git
`), 0644)
	if err != nil {
		panic(err)
	}

	// execute
	cmd := NewUpdateCmd(fs)
	b := bytes.NewBufferString("")
	cmd.SetOut(b)
	cmd.SetArgs([]string{"Mono-firefox-theme"})
	err = cmd.Execute()
	if err != nil {
		panic(err)
	}
	out, err := io.ReadAll(b)
	if err != nil {
		panic(err)
	}

	// assert
	if !reflect.DeepEqual(strings.Split(string(out), "\n")[:4], []string{
		"Deleted Mono-firefox-theme directory from " + chromeDir,
		"Removed imports from userChrome.css",
		"Removed imports from userContent.css",
		"Installing Mono-firefox-theme to " + chromeDir,
	}) {
		t.Errorf("Invalid output: %s", string(out))
	}
}
