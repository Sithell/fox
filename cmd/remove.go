package cmd

import (
	"errors"
	"fmt"
	"github.com/Sithell/fox/internal"
	"github.com/spf13/cobra"
	"os"
	"strings"
)

// removeCmd represents the remove command
var removeCmd = &cobra.Command{
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

		chromePath, err := internal.PrepareFirefox()
		if err != nil {
			panic(err)
		}

		err = os.RemoveAll(chromePath + "/" + repoName)
		if err != nil {
			panic(err)
		}
		fmt.Printf("Deleted %s directory from %s\n", repoName, chromePath)

		for _, userFileName := range []internal.UserFile{internal.Chrome, internal.Content} {
			filename := chromePath + "/" + string(userFileName)
			bytes, _ := os.ReadFile(filename)
			content := removeImportsFromFile(string(bytes), repoName)
			err = os.WriteFile(filename, []byte(content), 0644)
			if err != nil {
				panic(err)
			}
			fmt.Printf("Removed imports from %s\n", userFileName)
		}
	},
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
	rootCmd.AddCommand(removeCmd)
}
