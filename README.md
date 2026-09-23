# wd

`wd` (warp directory) is a tool that allows you to jump to custom directories in your terminal. It's a Go rewrite of the original [wd](https://github.com/mfaerevaag/wd) zsh plugin, designed for easier installation and cross-shell compatibility.

## Why wd?

The original `wd` is a zsh plugin. This Go implementation means:

- **Easy Installation**: Just a single binary, no need for complex shell plugin managers.
- **Cross-Shell Support**: Works with `bash`, `zsh`, and PowerShell.
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
eval "$(wd init bash)"  # Use 'zsh' if you are using Zsh
```

For PowerShell, add the following to your profile (`$PROFILE`):

```powershell
Invoke-Expression (@(wd init powershell) -join "`n")
```

The wrapper passes `wd`'s own commands (`list`, `completion`, and so on) straight through instead of treating them as warp points, and that list is built from the installed `wd` when `wd init` runs. Loading it with `eval`/`Invoke-Expression` as above rebuilds it in every new shell. If you saved the output of `wd init` to a file instead, regenerate that file after upgrading `wd` so new commands keep working.

### Tab completion

To complete commands and warp point names (`wd pr<Tab>`, `wd rm <Tab>`), also load the completion script after the wrapper:

```bash
source <(wd completion bash)  # or: source <(wd completion zsh)
```

```powershell
wd completion powershell | Out-String | Invoke-Expression
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

### Pick a warp point with fzf

If [fzf](https://github.com/junegunn/fzf) is installed, running `wd` with no arguments opens a fuzzy finder over your warp points and warps to the one you pick. Without fzf, it prints the help as before.

```bash
wd
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

- `add [point]`: Adds the current working directory to your warp points, replacing an existing point with the same name.
- `list`: Print all stored warp points.
- `rm <point>`: Removes the given warp point.
- `show [point]`: Print path to given warp point, or show points for current directory.
- `open <point>`: Open the warp point in the default file explorer.
- `ls <point>`: List files in the target warp point (without warping).
- `path <point>`: Show the path of a warp point.
- `clean`: Remove warp points that no longer exist.

## Credits

This is a fork of [mfaerevaag/wd](https://github.com/mfaerevaag/wd) written in Go for easier installation.
