# :white_square_button: CHMOD-CLI

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
##### Snap
```sh
$ sudo snap install chmod-cli
```


#### Build from source
> go version: 1.16
```sh
$ git clone github.com/Mayowa-Ojo/chmod-cli
$ cd chmod-cli
$ make build
```

## Usage
Run `chmod-cli` in your terminal to start the app.

You can also run `chmod-cli --help` to show an overview of the keybindings

## Navigation
| Key                     | Description                             |
|
| ----------------------- | --------------------------------------- |
| <kbd> up <kbd>          | Move up in the current section
| <kbd> down <kbd>        | Move down in the current section
| <kbd> left <kbd>        | Move left in the current section
| <kbd> right <kbd>       | Move right in the current section
| <kbd> tab/space <kbd>   | Move to the next section
| <kbd> shift+tab <kbd>   | Move to the previous section
| <kbd> shift+tab <kbd>   | Move to the previous section
| <kbd> Enter <kbd>       | Select/toggle current item
| <kbd> Ctrl+c <kbd>      | Copy command
| <kbd> Shift+? <kbd>     | toggle help
| <kbd> q <kbd>           | quit
|
