# :white_square_button: CHMOD-CLI
[![GitHub release](https://img.shields.io/github/v/release/Mayowa-Ojo/chmod-cli?sort=semver&style=social)](https://GitHub.com/Mayowa-Ojo/chmod-cli/releases/)


Simple cli tool that brings the chmod command in tui format. Generate permissions for files and directories by selecting easy to read config options and copy the result both in numeric and symbolic format.

<p align="center">
   <img width="600" src="docs/cast.svg">
</p>

## Installation
#### Packages

##### Homebrew
```sh
$ brew install Mayowa-Ojo/tap/chmod-cli
```

##### Download
Download one of the pre-compiled binaries from [releases](https://github.com/Mayowa-Ojo/chmod-cli/releases) and add the location to your `PATH`

#### Build from source
[![GitHub go.mod Go version of a Go module](https://img.shields.io/github/go-mod/go-version/Mayowa-Ojo/chmod-cli)](https://github.com/Mayowa-Ojo/chmod-cli)

```sh
$ git clone github.com/Mayowa-Ojo/chmod-cli
$ cd chmod-cli
$ make install
$ chmod-cli
```

## Usage
Run `chmod-cli` in your terminal to start the app.

You can also run `chmod-cli --help` to show an overview of the keybindings

## Navigation
| Key                      | Description                            |
| -----------------------  | -------------------------------------- |
| <kbd> up </kbd>          | Move up in the current section         |
| <kbd> down </kbd>        | Move down in the current section       |
| <kbd> left </kbd>        | Move left in the current section       |
| <kbd> right </kbd>       | Move right in the current section      |
| <kbd> tab/space </kbd>   | Move to the next section               |
| <kbd> shift+tab </kbd>   | Move to the previous section           |
| <kbd> shift+tab </kbd>   | Move to the previous section           |
| <kbd> Enter </kbd>       | Select/toggle current item             |
| <kbd> Ctrl+c </kbd>      | Copy command                           |
| <kbd> Shift+? </kbd>     | toggle help                            |
| <kbd> q </kbd>           | quit                                   |

## Built with
- [Bubbletea](https://github.com/charmbracelet/bubbletea)
- [Lipgloss](https://github.com/charmbracelet/lipgloss)
- [Urfave/cli](https://github.com/urfave/cli/v2)
- [Clipboard](https://github.com/atotto/clipboard)

[![GitHub license](https://img.shields.io/github/license/Naereen/StrapDown.js.svg)](https://github.com/Naereen/StrapDown.js/blob/master/LICENSE)
