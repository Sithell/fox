package cmd

import (
	"bytes"
	"github.com/Sithell/fox/internal"
	"github.com/spf13/afero"
	"gopkg.in/yaml.v3"
	"io"
	"os"
	"testing"
)

func TestNewRemoveCmd(t *testing.T) {
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
	cmd := NewRemoveCmd(fs)
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
	if string(out) != `Deleted Mono-firefox-theme directory from `+chromeDir+`
Removed imports from userChrome.css
Removed imports from userContent.css
` {
		t.Errorf("Invalid output: %s", string(out))
	}
	if _, err := fs.Stat(chromeDir + "/Mono-firefox-theme"); !os.IsNotExist(err) {
		t.Errorf("Mono-firefox-theme directory not removed: %s", err)
	}

	foxYml, err := afero.ReadFile(fs, chromeDir+"/fox.yml")
	if err != nil {
		t.Errorf("Failed to open fox.yml: %s", err)
	}
	expected := internal.FoxYml{Mods: []internal.FoxYmlMod{}}
	content, _ := yaml.Marshal(expected)
	if string(foxYml) != string(content) {
		t.Errorf("Contents of fox.yml do not match: %s", string(foxYml))
	}
}
