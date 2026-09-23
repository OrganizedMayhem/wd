package main

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"testing"

	"github.com/OrganizedMayhem/wd/cmd"
)

// TestWDHelperProcess runs the CLI in a separate process so exit behavior can
// be asserted without terminating the test runner.
func TestWDHelperProcess(t *testing.T) {
	if os.Getenv("WD_HELPER_PROCESS") != "1" {
		return
	}

	separator := -1
	for i, arg := range os.Args {
		if arg == "--" {
			separator = i
			break
		}
	}
	if separator == -1 {
		os.Exit(2)
	}

	os.Args = append([]string{"wd"}, os.Args[separator+1:]...)
	cmd.Execute()
	os.Exit(0)
}

type invocation struct {
	home string
	dir  string
	env  []string
}

func runWD(t *testing.T, options invocation, args ...string) string {
	t.Helper()

	commandArgs := append([]string{"-test.run=^TestWDHelperProcess$", "--"}, args...)
	command := exec.Command(os.Args[0], commandArgs...)
	command.Dir = options.dir
	command.Env = helperEnv(options.home, options.env...)

	output, err := command.CombinedOutput()
	if err != nil {
		t.Fatalf("wd %s failed: %v\n%s", strings.Join(args, " "), err, output)
	}
	return string(output)
}

func runWDFailure(t *testing.T, options invocation, args ...string) string {
	t.Helper()

	commandArgs := append([]string{"-test.run=^TestWDHelperProcess$", "--"}, args...)
	command := exec.Command(os.Args[0], commandArgs...)
	command.Dir = options.dir
	command.Env = helperEnv(options.home, options.env...)

	output, err := command.CombinedOutput()
	if err == nil {
		t.Fatalf("wd %s succeeded, want failure\n%s", strings.Join(args, " "), output)
	}
	return string(output)
}

// tempDir returns a test directory with symlinks resolved, so paths reported
// by the OS (macOS /var -> /private/var, Windows 8.3 short names) match the
// paths the tests compare against.
func tempDir(t *testing.T) string {
	t.Helper()
	dir, err := filepath.EvalSymlinks(t.TempDir())
	if err != nil {
		t.Fatal(err)
	}
	return dir
}

// helperEnv points both HOME (Unix) and USERPROFILE (Windows) at home, since
// os.UserHomeDir reads a different variable per platform.
func helperEnv(home string, extra ...string) []string {
	overridden := map[string]bool{"HOME": true, "USERPROFILE": true, "WD_HELPER_PROCESS": true}
	for _, value := range extra {
		key, _, _ := strings.Cut(value, "=")
		overridden[key] = true
	}

	env := make([]string, 0, len(os.Environ())+len(extra)+2)
	for _, value := range os.Environ() {
		key, _, _ := strings.Cut(value, "=")
		if !overridden[key] {
			env = append(env, value)
		}
	}
	return append(env, append([]string{"HOME=" + home, "USERPROFILE=" + home, "WD_HELPER_PROCESS=1"}, extra...)...)
}

func writeWarpConfig(t *testing.T, home, content string) {
	t.Helper()
	if err := os.WriteFile(filepath.Join(home, ".warprc"), []byte(content), 0o644); err != nil {
		t.Fatal(err)
	}
}

func TestWarpPointVerbs(t *testing.T) {
	home := tempDir(t)
	target := tempDir(t)
	test := invocation{home: home, dir: target}

	if got, want := runWD(t, test, "addcd", target, "alpha"), fmt.Sprintf("Added warp point 'alpha' to '%s'\n", target); got != want {
		t.Errorf("addcd output = %q, want %q", got, want)
	}
	if got, want := runWD(t, test, "list"), "alpha:"+target+"\n"; got != want {
		t.Errorf("list output = %q, want %q", got, want)
	}
	if got, want := runWD(t, test, "path", "alpha"), target+"\n"; got != want {
		t.Errorf("path output = %q, want %q", got, want)
	}
	if got, want := runWD(t, test, "show", "alpha"), target+"\n"; got != want {
		t.Errorf("show <point> output = %q, want %q", got, want)
	}
	if got, want := runWD(t, test, "show"), "Warp points for current directory: alpha\n"; got != want {
		t.Errorf("show output = %q, want %q", got, want)
	}
	if got, want := runWD(t, test, "alpha"), target+"\n"; got != want {
		t.Errorf("warp output = %q, want %q", got, want)
	}

	replacement := tempDir(t)
	if got, want := runWD(t, test, "addcd", replacement, "alpha"), fmt.Sprintf("Added warp point 'alpha' to '%s'\n", replacement); got != want {
		t.Errorf("duplicate addcd output = %q, want %q", got, want)
	}
	if got, want := runWD(t, test, "alpha"), replacement+"\n"; got != want {
		t.Errorf("updated warp output = %q, want %q", got, want)
	}

	if got, want := runWD(t, test, "add", "current"), fmt.Sprintf("Added warp point 'current' to '%s'\n", target); got != want {
		t.Errorf("add output = %q, want %q", got, want)
	}
	if got, want := runWD(t, test, "rm", "alpha"), "Removed warp point 'alpha'\n"; got != want {
		t.Errorf("rm output = %q, want %q", got, want)
	}
	if got, want := runWD(t, test, "list"), "current:"+target+"\n"; got != want {
		t.Errorf("list after rm = %q, want %q", got, want)
	}
}

func TestAddcdStoresAbsolutePath(t *testing.T) {
	home := tempDir(t)
	parent := tempDir(t)
	target := filepath.Join(parent, "child")
	if err := os.Mkdir(target, 0o755); err != nil {
		t.Fatal(err)
	}

	runWD(t, invocation{home: home, dir: parent}, "addcd", "child")
	if got, want := runWD(t, invocation{home: home}, "path", "child"), target+"\n"; got != want {
		t.Errorf("path for relative addcd = %q, want %q", got, want)
	}
}

func TestShowMatchesSymlinkedDirectory(t *testing.T) {
	home := tempDir(t)
	target := tempDir(t)
	link := filepath.Join(tempDir(t), "link")
	if err := os.Symlink(target, link); err != nil {
		t.Skipf("cannot create symlink: %v", err)
	}
	writeWarpConfig(t, home, "linked:"+link+"\n")

	if got, want := runWD(t, invocation{home: home, dir: target}, "show"), "Warp points for current directory: linked\n"; got != want {
		t.Errorf("show output = %q, want %q", got, want)
	}
}

func TestCleanVerbRemovesMissingWarpPoints(t *testing.T) {
	home := tempDir(t)
	target := tempDir(t)
	missing := filepath.Join(home, "does-not-exist")
	writeWarpConfig(t, home, "live:"+target+"\nstale:"+missing+"\n")

	if got, want := runWD(t, invocation{home: home}, "clean"), "Cleaned warp points.\n"; got != want {
		t.Errorf("clean output = %q, want %q", got, want)
	}

	contents, err := os.ReadFile(filepath.Join(home, ".warprc"))
	if err != nil {
		t.Fatal(err)
	}
	if got, want := string(contents), "live:"+target+"\n"; got != want {
		t.Errorf("cleaned config = %q, want %q", got, want)
	}
}

func TestLSVerbListsWarpPointContents(t *testing.T) {
	home := tempDir(t)
	target := tempDir(t)
	if err := os.WriteFile(filepath.Join(target, "listed.txt"), []byte("contents"), 0o644); err != nil {
		t.Fatal(err)
	}
	writeWarpConfig(t, home, "files:"+target+"\n")

	if got, want := runWD(t, invocation{home: home}, "ls", "files"), "listed.txt\n"; got != want {
		t.Errorf("ls output = %q, want %q", got, want)
	}
}

func TestOpenVerbUsesPlatformOpener(t *testing.T) {
	if runtime.GOOS != "linux" {
		t.Skip("the fake platform opener is specific to the Linux xdg-open implementation")
	}

	home := tempDir(t)
	target := tempDir(t)
	binDir := tempDir(t)
	logFile := filepath.Join(tempDir(t), "open-path")
	openCommand := filepath.Join(binDir, "xdg-open")
	if err := os.WriteFile(openCommand, []byte("#!/bin/sh\nprintf '%s\\n' \"$1\" > \"$WD_OPEN_LOG\"\n"), 0o755); err != nil {
		t.Fatal(err)
	}
	writeWarpConfig(t, home, "project:"+target+"\n")

	runWD(t, invocation{
		home: home,
		env:  []string{"PATH=" + binDir, "WD_OPEN_LOG=" + logFile},
	}, "open", "project")

	contents, err := os.ReadFile(logFile)
	if err != nil {
		t.Fatal(err)
	}
	if got, want := string(contents), target+"\n"; got != want {
		t.Errorf("opener path = %q, want %q", got, want)
	}
}

func TestCommandErrors(t *testing.T) {
	home := tempDir(t)
	test := invocation{home: home, dir: tempDir(t)}

	if got := runWDFailure(t, test, "path", "missing"); !strings.Contains(got, "no warp points yet") {
		t.Errorf("missing config error = %q, want no warp points error", got)
	}
	writeWarpConfig(t, home, "malformed-config\n")
	if got := runWDFailure(t, test, "path", "missing"); !strings.Contains(got, "invalid warp point") {
		t.Errorf("malformed config error = %q, want invalid warp point error", got)
	}
	if got := runWDFailure(t, test, "init", "fish"); !strings.Contains(got, "unsupported shell") {
		t.Errorf("unsupported shell error = %q, want unsupported shell error", got)
	}

	writeWarpConfig(t, home, "")
	for _, name := range []string{"list", "List", "completion", "-x"} {
		if got := runWDFailure(t, test, "add", "--", name); !strings.Contains(got, "is reserved") {
			t.Errorf("add %s error = %q, want reserved name error", name, got)
		}
	}
	if got := runWDFailure(t, test, "path", "missing"); !strings.Contains(got, `warp point "missing" not found`) {
		t.Errorf("unknown point error = %q, want not found error", got)
	}
	if got := runWDFailure(t, test, "rm", "missing"); !strings.Contains(got, `warp point "missing" not found`) {
		t.Errorf("rm unknown point error = %q, want not found error", got)
	}
}

// writeShim installs a `wd` executable in binDir that re-runs the test binary
// as the CLI: a batch file on Windows, a shebang script elsewhere.
func writeShim(t *testing.T, binDir string) {
	t.Helper()
	shim := filepath.Join(binDir, "wd")
	shimContents := fmt.Sprintf("#!/bin/sh\nexec %q -test.run=^TestWDHelperProcess$ -- \"$@\"\n", os.Args[0])
	if runtime.GOOS == "windows" {
		shim += ".cmd"
		shimContents = fmt.Sprintf("@\"%s\" \"-test.run=^TestWDHelperProcess$\" -- %%*\r\n@exit /b %%ERRORLEVEL%%\r\n", os.Args[0])
	}
	if err := os.WriteFile(shim, []byte(shimContents), 0o755); err != nil {
		t.Fatal(err)
	}
}

// shellAdapter captures what differs between the generated shell wrappers so
// the scenarios in runShellWrapperSuite can be written once and run against
// bash, zsh, and PowerShell alike.
type shellAdapter struct {
	name       string   // human-readable label used in skip/error messages
	initArg    string   // argument passed to `wd init`
	binaries   []string // candidate executables to look up, in preference order
	wrapperExt string   // file extension for the generated wrapper script
	pwdExpr    string   // expression that prints the current directory
	exitExpr   string   // expression that prints the last command's exit code as "EXIT:<code>"
	command    func(binary, wrapperPath, script string) *exec.Cmd
}

var bashAdapter = shellAdapter{
	name:       "bash",
	initArg:    "bash",
	binaries:   []string{"bash"},
	wrapperExt: ".sh",
	pwdExpr:    "pwd",
	exitExpr:   `echo "EXIT:$?"`,
	command: func(binary, wrapperPath, script string) *exec.Cmd {
		return exec.Command(binary, "-c", fmt.Sprintf(". %q; %s", wrapperPath, script))
	},
}

var zshAdapter = shellAdapter{
	name:       "zsh",
	initArg:    "zsh",
	binaries:   []string{"zsh"},
	wrapperExt: ".sh",
	pwdExpr:    "pwd",
	exitExpr:   `echo "EXIT:$?"`,
	command: func(binary, wrapperPath, script string) *exec.Cmd {
		return exec.Command(binary, "-c", fmt.Sprintf(". %q; %s", wrapperPath, script))
	},
}

var powershellAdapter = shellAdapter{
	name:       "powershell",
	initArg:    "powershell",
	binaries:   []string{"pwsh", "powershell"},
	wrapperExt: ".ps1",
	pwdExpr:    "Write-Output (Get-Location).Path",
	exitExpr:   `Write-Output "EXIT:$LASTEXITCODE"`,
	command: func(binary, wrapperPath, script string) *exec.Cmd {
		return exec.Command(binary, "-NoProfile", "-NonInteractive", "-Command", fmt.Sprintf(". '%s'; %s", wrapperPath, script))
	},
}

// lastLine returns the final non-empty line of shell output, tolerating both
// LF and CRLF line endings.
func lastLine(output string) string {
	lines := strings.Split(strings.TrimRight(output, "\r\n"), "\n")
	return strings.TrimSpace(lines[len(lines)-1])
}

// runShellWrapperSuite exercises the wrapper `wd init <shell>` generates:
// warping to a point, passthrough commands bypassing the warp, failed warps
// leaving the shell where it started, and paths containing spaces.
func runShellWrapperSuite(t *testing.T, adapter shellAdapter) {
	if runtime.GOOS == "windows" && adapter.name != "powershell" {
		t.Skipf("the %s wrapper is not exercised on windows", adapter.name)
	}

	var binary string
	for _, candidate := range adapter.binaries {
		if _, err := exec.LookPath(candidate); err == nil {
			binary = candidate
			break
		}
	}
	if binary == "" {
		t.Skipf("no %s binary installed (tried %s)", adapter.name, strings.Join(adapter.binaries, ", "))
	}

	home := tempDir(t)
	binDir := tempDir(t)
	writeShim(t, binDir)

	wrapper := filepath.Join(tempDir(t), "wd-wrapper"+adapter.wrapperExt)
	if err := os.WriteFile(wrapper, []byte(runWD(t, invocation{home: home}, "init", adapter.initArg)), 0o600); err != nil {
		t.Fatal(err)
	}

	run := func(t *testing.T, dir, script string) string {
		t.Helper()
		shell := adapter.command(binary, wrapper, script)
		shell.Dir = dir
		shell.Env = helperEnv(home, "PATH="+binDir)
		output, err := shell.CombinedOutput()
		if err != nil {
			t.Fatalf("shell script failed: %v\n%s", err, output)
		}
		return string(output)
	}

	t.Run("WarpsToPoint", func(t *testing.T) {
		target := tempDir(t)
		start := tempDir(t)
		writeWarpConfig(t, home, "alpha:"+target+"\n")

		output := run(t, start, "wd alpha; "+adapter.pwdExpr)
		if got, want := lastLine(output), target; got != want {
			t.Errorf("cwd after warp = %q, want %q", got, want)
		}
	})

	t.Run("PassthroughCommandDoesNotWarp", func(t *testing.T) {
		target := tempDir(t)
		start := tempDir(t)
		writeWarpConfig(t, home, "alpha:"+target+"\n")

		output := run(t, start, "wd list; "+adapter.pwdExpr)
		if !strings.Contains(output, "alpha:"+target) {
			t.Errorf("passthrough output = %q, want it to contain the warp point listing", output)
		}
		if got, want := lastLine(output), start; got != want {
			t.Errorf("cwd after passthrough command = %q, want unchanged %q", got, want)
		}
	})

	t.Run("CompletionCommandDoesNotWarp", func(t *testing.T) {
		start := tempDir(t)

		output := run(t, start, "wd completion bash; "+adapter.exitExpr+"; "+adapter.pwdExpr)
		if !strings.Contains(output, "bash completion") || !strings.Contains(output, "EXIT:0") {
			t.Errorf("completion output = %q, want the generated completion script and EXIT:0", output)
		}
		if got, want := lastLine(output), start; got != want {
			t.Errorf("cwd after completion command = %q, want unchanged %q", got, want)
		}

		output = run(t, start, "wd __complete ad; "+adapter.exitExpr+"; "+adapter.pwdExpr)
		if !strings.Contains(output, "addcd") || !strings.Contains(output, "EXIT:0") {
			t.Errorf("__complete output = %q, want completion candidates and EXIT:0", output)
		}
		if got, want := lastLine(output), start; got != want {
			t.Errorf("cwd after __complete = %q, want unchanged %q", got, want)
		}
	})

	t.Run("UnknownPointFailsWithoutWarping", func(t *testing.T) {
		start := tempDir(t)
		writeWarpConfig(t, home, "alpha:"+tempDir(t)+"\n")

		output := run(t, start, "wd missing; "+adapter.exitExpr+"; "+adapter.pwdExpr)
		if !strings.Contains(output, "EXIT:1") {
			t.Errorf("output = %q, want it to contain EXIT:1", output)
		}
		if got, want := lastLine(output), start; got != want {
			t.Errorf("cwd after failed warp = %q, want unchanged %q", got, want)
		}
	})

	t.Run("WarpsToPointWithSpaces", func(t *testing.T) {
		target := filepath.Join(tempDir(t), "point with spaces")
		if err := os.MkdirAll(target, 0o755); err != nil {
			t.Fatal(err)
		}
		start := tempDir(t)
		writeWarpConfig(t, home, "spacey:"+target+"\n")

		output := run(t, start, "wd spacey; "+adapter.pwdExpr)
		if got, want := lastLine(output), target; got != want {
			t.Errorf("cwd after warp to spaced point = %q, want %q", got, want)
		}
	})
}

func TestBashWrapperSuite(t *testing.T)       { runShellWrapperSuite(t, bashAdapter) }
func TestZshWrapperSuite(t *testing.T)        { runShellWrapperSuite(t, zshAdapter) }
func TestPowerShellWrapperSuite(t *testing.T) { runShellWrapperSuite(t, powershellAdapter) }

func TestInformationalVerbs(t *testing.T) {
	home := tempDir(t)
	test := invocation{home: home}

	if got := runWD(t, test, "init", "bash"); !strings.Contains(got, "wd_cd()") || !strings.Contains(got, "target_path=$(command wd \"$@\")") || !strings.Contains(got, "wd_cd \"$@\"") {
		t.Errorf("init bash output does not contain the shell wrapper that invokes the wd binary:\n%s", got)
	}
	if got := runWD(t, test, "init", "zsh"); !strings.Contains(got, "wd_cd()") || !strings.Contains(got, "target_path=$(command wd \"$@\")") || !strings.Contains(got, "wd_cd \"$@\"") {
		t.Errorf("init zsh output does not contain the shell wrapper that invokes the wd binary:\n%s", got)
	}
	if got := runWD(t, test, "init", "powershell"); !strings.Contains(got, "function wd {") || !strings.Contains(got, "Get-Command wd -CommandType Application") || !strings.Contains(got, "Set-Location -LiteralPath $targetPath") {
		t.Errorf("init powershell output does not contain the shell wrapper that invokes the wd binary:\n%s", got)
	}
	if got, want := runWD(t, test, "init", "powershell"), runWD(t, test, "init", "pwsh"); got != want {
		t.Errorf("init pwsh output = %q, want same as init powershell = %q", want, got)
	}
	if got, want := runWD(t, test, "version"), "wd version 0.10.1\n"; got != want {
		t.Errorf("version output = %q, want %q", got, want)
	}
	if got := runWD(t, test, "help"); !strings.Contains(got, "Available Commands:") || !strings.Contains(got, "addcd") {
		t.Errorf("help output did not list commands:\n%s", got)
	}
	if got := runWD(t, test); !strings.Contains(got, "Available Commands:") || !strings.Contains(got, "addcd") {
		t.Errorf("empty invocation did not show help:\n%s", got)
	}
}
