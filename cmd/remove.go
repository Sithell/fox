package cmd

import (
	"errors"
	"fmt"
	"github.com/Sithell/fox/internal"
	"github.com/spf13/afero"
	"github.com/spf13/cobra"
	"gopkg.in/yaml.v3"
	"strings"
)

func NewRemoveCmd(fs afero.Fs) *cobra.Command {
	return &cobra.Command{
		Use:   "remove",
		Short: "Remove installed mod by name",
		Long:  `Remove command accepts a single argument: name of the mod folder inside of the chrome directory`,
		Args: func(cmd *cobra.Command, args []string) error {
			if len(args) != 1 {
				return errors.New("Remove command accepts a single argument: name of the mod folder inside of the chrome directory")
			}
			return nil
		},
		Run: func(cmd *cobra.Command, args []string) {
			var repoName = args[0]

			chromePath, err := internal.PrepareFirefox(fs)
			if err != nil {
				panic(err)
			}

			err = fs.RemoveAll(chromePath + "/" + repoName)
			if err != nil {
				panic(err)
			}

			_, err = fmt.Fprintf(cmd.OutOrStdout(), "Deleted %s directory from %s\n", repoName, chromePath)
			if err != nil {
				panic(err)
			}

			for _, userFileName := range []internal.UserFile{internal.Chrome, internal.Content} {
				filename := chromePath + "/" + string(userFileName)
				bytes, _ := afero.ReadFile(fs, filename)
				content := removeImportsFromFile(string(bytes), repoName)
				err = afero.WriteFile(fs, filename, []byte(content), 0644)
				if err != nil {
					panic(err)
				}
				_, err := fmt.Fprintf(cmd.OutOrStdout(), "Removed imports from %s\n", userFileName)
				if err != nil {
					panic(err)
				}
			}

			foxYml, err := internal.LoadFoxYml(fs, chromePath+"/fox.yml")
			if err != nil {
				panic(err)
			}

			for i, mod := range foxYml.Mods {
				if mod.Name == repoName {
					foxYml.Mods = internal.Remove(foxYml.Mods, i)
					break
				}
			}

			rawYaml, err := yaml.Marshal(foxYml)
			err = afero.WriteFile(fs, chromePath+"/fox.yml", rawYaml, 0644)
			if err != nil {
				panic(err)
			}
		},
	}
}

func removeImportsFromFile(s string, dirName string) string {
	lines := strings.Split(s, "\n")
	for i, line := range lines {
		if strings.Contains(line, fmt.Sprintf("@import \"%s", dirName)) {
			internal.Remove(lines, i)
		}
	}

	return strings.Join(lines, "\n")
}

func init() {
	rootCmd.AddCommand(NewRemoveCmd(afero.NewOsFs()))
}
