package cmd

import (
	"errors"
	"fmt"
	"github.com/Sithell/fox/internal"
	"github.com/go-git/go-git/v5"
	"github.com/spf13/cobra"
	u "net/url"
	"os"
	"regexp"
	"strings"
)

// installCmd represents the "install" command
var installCmd = &cobra.Command{
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

		chromePath, err := internal.PrepareFirefox()
		if err != nil {
			panic(err)
		}

		re, err := regexp.Compile("([A-Za-z0-9_\\-]+).git$")
		if err != nil {
			panic(err)
		}
		repoName := re.FindAllStringSubmatch(url, -1)[0][1]
		fmt.Printf("Installing %s to %s\n", repoName, chromePath)

		_, err = git.PlainClone(chromePath+"/"+repoName, false, &git.CloneOptions{
			URL:      url,
			Progress: os.Stdout,
		})
		if err != nil {
			if err.Error() == "repository already exists" {
				fmt.Println(err)
			} else {
				panic(err)
			}
		}

		userFiles := internal.LocateUserFiles(chromePath + "/" + repoName)
		for userFileName, pathInRepo := range userFiles {
			filename := chromePath + "/" + string(userFileName)
			bytes, _ := os.ReadFile(filename)
			content := addImportToFile(string(bytes), repoName+"/"+pathInRepo)
			err = os.WriteFile(filename, []byte(content), 0644)
			if err != nil {
				panic(err)
			}
		}
	},
}

func addImportToFile(s string, path string) string {
	if strings.Contains(s, "/* Fox mod manager start */\n") &&
		strings.Contains(s, "\n/* Fox mod manager end */") {
		regionStart := strings.Index(s, "/* Fox mod manager start */\n") + 28
		regionEnd := strings.Index(s, "\n/* Fox mod manager end */")
		return s[:regionStart] + addImportToRegion(s[regionStart:regionEnd], path) + s[regionEnd:]
	}
	return fmt.Sprintf("%s\n/* Fox mod manager start */\n@import \"%s\";\n/* Fox mod manager end */\n", s, path)
}

func addImportToRegion(region string, path string) string {
	lines := strings.Split(region, "\n")
	if !internal.Contains(lines, fmt.Sprintf("@import \"%s\";", path)) {
		lines = append(lines, fmt.Sprintf("@import \"%s\";", path))
	}
	if len(lines) == 0 {
		return ""
	}
	return strings.Join(lines, "\n")
}

func init() {
	rootCmd.AddCommand(installCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// installCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// installCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
