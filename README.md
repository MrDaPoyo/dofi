# Dofi!!!!
Dofi is an open source fantasy console written in Go.
Dofi's main goal is to be a fun and easy to use console for learning programming and basic game development.

## Features
Some of Dofi's features are:
- Balena language
- Code editor
- CLI
- Cute icons (made by a friend)
- Cute cursor!!! (made by me)

The Balena programming language is a simple, optimized to use the least amount of characters and very slow interpreted language that I built from scratch, for the sole purpose of being used in Dofi.

Balena requires a mouse as of right now, and is not compatible with touchscreens.

The code editor just works as a text editor, and is not a full-featured code editor. It's based on the same text engine as the CLI, which I wrote from scratch and was a pain to implement (I spent several hours of my life working on it).

## Keybindings
`ESC` - Toggle between the code editor and the CLI

That's.... pretty much it.

## Compiling to WASM

On a Linux machine, you can do the following to compile and build to WASM:

```
env GOOS=js GOARCH=wasm go build -o dofi.wasm github.com/mrdapoyo/dofi
```

On a Windows machine, you can do the following to compile and build to WASM:

```
$Env:GOOS = 'js'
$Env:GOARCH = 'wasm'
go build -o dofi.wasm github.com/mrdapoyo/dofi
Remove-Item Env:GOOS
Remove-Item Env:GOARCH
```

## Compiling to a binary
To compile to a binary, it's as simple as running the following command:
`go build -o dofi main.go`