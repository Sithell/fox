package cmd

import (
	"errors"
	"fmt"
	"github.com/Sithell/fox/internal"
	"github.com/go-git/go-git/v5"
	"github.com/spf13/afero"
	"github.com/spf13/cobra"
)

func NewUpdateCmd(fs afero.Fs) *cobra.Command {
	return &cobra.Command{
		Use:   "update",
		Short: "Update installed mod by name",
		Long:  `Update command accepts a single argument: name of the mod folder inside of the chrome directory`,
		Args: func(cmd *cobra.Command, args []string) error {
			if len(args) != 1 {
				return errors.New("update command accepts a single argument: name of the mod folder inside of the chrome directory")
			}
			return nil
		},
		Run: func(cmd *cobra.Command, args []string) {
			var repoName = args[0]

			chromePath, err := internal.PrepareFirefox(fs)
			if err != nil {
				panic(err)
			}

			foxYml, err := internal.LoadFoxYml(fs, chromePath+"/fox.yml")
			if err != nil {
				panic(err)
			}

			var foxMod internal.FoxYmlMod
			for _, mod := range foxYml.Mods {
				if mod.Name == repoName {
					foxMod = mod
					break
				}
			}
			if foxMod.Name == "" {
				fmt.Println(fmt.Errorf("mod %s not found in fox.yml, check if it was installed via fox", repoName))
				return
			}

			// Delete mod
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

			// Install it again
			_, err = fmt.Fprintf(cmd.OutOrStdout(), "Installing %s to %s\n", repoName, chromePath)
			if err != nil {
				panic(err)
			}
			_, err = git.PlainClone(chromePath+"/"+repoName, false, &git.CloneOptions{
				URL:      foxMod.Url,
				Progress: cmd.OutOrStdout(),
			})
			if err != nil {
				if err.Error() == "repository already exists" {
					_, err = fmt.Fprintln(cmd.OutOrStderr(), err)
					if err != nil {
						panic(err)
					}
				} else {
					panic(err)
				}
			}

			userFiles := internal.LocateUserFiles(fs, chromePath+"/"+repoName)
			for userFileName, pathInRepo := range userFiles {
				filename := chromePath + "/" + string(userFileName)
				bytes, _ := afero.ReadFile(fs, filename)
				content := internal.AddImportToFile(string(bytes), repoName+"/"+pathInRepo)
				err = afero.WriteFile(fs, filename, []byte(content), 0644)
				if err != nil {
					panic(err)
				}
			}
		},
	}
}

func init() {
	rootCmd.AddCommand(NewUpdateCmd(afero.NewOsFs()))
}
