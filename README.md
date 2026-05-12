# wd (go re-write)

`wd` (warp directory) is a tool that allows you to jump to custom directories in your terminal. It's a Go rewrite of the original [wd](https://github.com/mfaerevaag/wd) zsh plugin, designed for easier installation and cross-shell compatibility.

## Why wd-go?

The original `wd` is a zsh plugin. This `wd` is written in Go, which means:
- **Easy Installation**: Just a single binary, no need for complex shell plugin managers.
- **Cross-Shell Support**: Works with `bash` and `zsh`.
- **Fast**: Built with Go for speed and efficiency.

## Installation

### Using Go

```bash
go install github.com/OrganizedMayhem/wd@latest
```

## Setup

Since a child process cannot change the parent shell's directory, `wd` requires a small shell wrapper. 

Add the following to your `.bashrc`, `.zshrc`, or equivalent:

```bash
eval "$(wd-go init bash)"  # Use 'zsh' if you are using Zsh
```

## Usage

### Add a warp point
```bash
# Add current directory as 'myproject'
wd add myproject

# Add current directory using the directory name
wd add
```

### Warp to a point
```bash
wd myproject
```

### List all warp points
```bash
wd list
```

### Remove a warp point
```bash
wd rm myproject
```

### Show current warp point
```bash
wd show
```

## Commands

- `add [point]`: Adds the current working directory to your warp points.
- `list`: Print all stored warp points.
- `rm <point>`: Removes the given warp point.
- `show [point]`: Print path to given warp point, or show points for current directory.
- `open <point>`: Open the warp point in the default file explorer.
- `ls <point>`: List files in the target warp point (without warping).
- `path <point>`: Show the path of a warp point.
- `clean`: Remove warp points that no longer exist.

## Credits

This is a fork of [mfaerevaag/wd](https://github.com/mfaerevaag/wd) written in Go for easier installation.
