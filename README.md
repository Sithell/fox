# Fox

Simple mod manager for Firefox written in Go.

## Description

Fox allows you to manage Firefox UI modifications via CLI by editing `userChrome.css` and other files. 

## Usage
```shell
# To install a mod
fox install https://github.com/witalihirsch/Mono-firefox-theme.git

# To remove a mod
fox remove Mono-firefox-theme

# To show help message
fox help

```

## Dependencies

Fox can run on every major platform, including:
  * Linux
  * Windows
  * MacOS

Naturally, you need Firefox installed.

Technologies used are:
- [Cobra](https://github.com/spf13/cobra)
- [Afero](https://github.com/spf13/afero)
- [go-git](https://github.com/go-git/go-git)

## Installing

### Building from source
```shell
go build .
```

### Pre-compiled binaries
Coming soon...

## Help

Feel free to open issues and submit pull requests!
