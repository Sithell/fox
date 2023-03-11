package cmd

import (
	"errors"
	"fmt"
	"github.com/Sithell/fox/internal"
	"github.com/go-git/go-git/v5"
	"github.com/spf13/afero"
	"github.com/spf13/cobra"
	"gopkg.in/yaml.v3"
	u "net/url"
	"regexp"
)

func NewInstallCmd(fs afero.Fs) *cobra.Command {
	return &cobra.Command{
		Use:   "install",
		Short: "Install mod by a git link",
		Long:  `Install mod by a git link.`,
		Args: func(cmd *cobra.Command, args []string) error {
			if len(args) < 1 {
				return errors.New("missing argument: url")
			}
			var url = args[0]
			if _, err := u.ParseRequestURI(url); err != nil {
				return fmt.Errorf("not a valid url: %s", url)
			}
			if url[len(url)-4:] != ".git" {
				return fmt.Errorf("not a git repository: %s", url)
			}
			return nil
		},
		Run: func(cmd *cobra.Command, args []string) {
			var url = args[0]

			chromePath, err := internal.PrepareFirefox(fs)
			if err != nil {
				panic(err)
			}

			re, err := regexp.Compile("([A-Za-z0-9_\\-]+).git$")
			if err != nil {
				panic(err)
			}
			repoName := re.FindAllStringSubmatch(url, -1)[0][1]
			_, err = fmt.Fprintf(cmd.OutOrStdout(), "Installing %s to %s\n", repoName, chromePath)
			if err != nil {
				panic(err)
			}

			_, err = git.PlainClone(chromePath+"/"+repoName, false, &git.CloneOptions{
				URL:      url,
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

			foxYml, err := internal.LoadFoxYml(fs, chromePath+"/fox.yml")
			if err != nil {
				panic(err)
			}

			mod := internal.FoxYmlMod{Name: repoName, Url: url}
			if !internal.Contains(foxYml.Mods, mod) {
				foxYml.Mods = append(foxYml.Mods, mod)
			}

			rawYaml, err := yaml.Marshal(foxYml)
			err = afero.WriteFile(fs, chromePath+"/fox.yml", rawYaml, 0644)
			if err != nil {
				panic(err)
			}
		},
	}
}

func init() {
	rootCmd.AddCommand(NewInstallCmd(afero.NewOsFs()))
}
