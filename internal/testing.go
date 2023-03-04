package internal

import (
	"github.com/spf13/afero"
	"os"
)

func InitFs(firefoxPath string, profileName string) afero.Fs {
	fs := afero.NewMemMapFs()
	if err := fs.MkdirAll(firefoxPath+"/"+profileName, os.ModePerm); err != nil {
		panic(err)
	}

	file, err := fs.Create(firefoxPath + "/profiles.ini")
	if err != nil {
		panic(err)
	}
	content := `
[Install4F96D1932A9F858E]
Default=` + profileName + `
Locked=1
`
	if _, err = file.WriteString(content); err != nil {
		panic(err)
	}
	if err = file.Close(); err != nil {
		panic(err)
	}
	return fs
}
