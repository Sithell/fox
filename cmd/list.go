package cmd

import (
	"fmt"
	"github.com/Sithell/fox/internal"
	"github.com/spf13/afero"
	"github.com/spf13/cobra"
)

func NewListCmd(fs afero.Fs) *cobra.Command {
	return &cobra.Command{
		Use:   "list",
		Short: "List installed mods",
		Long:  `List installed mods.`,
		Run: func(cmd *cobra.Command, args []string) {
			chromePath, err := internal.PrepareFirefox(fs)
			if err != nil {
				panic(err)
			}

			foxYml, err := internal.LoadFoxYml(fs, chromePath+"/fox.yml")
			if err != nil {
				panic(err)
			}

			for i, mod := range foxYml.Mods {
				_, err := fmt.Fprintf(cmd.OutOrStdout(), "%d. %s - %s", i+1, mod.Name, mod.Url)
				if err != nil {
					panic(err)
				}
			}
		},
	}
}

func init() {
	rootCmd.AddCommand(NewListCmd(afero.NewOsFs()))
}
