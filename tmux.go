package main

import (
	"fmt"
	"os"
	"os/exec"
	"strings"
)

// tmuxReentryEnv is set on the re-exec'd child so we don't try to detach a
// second time (in case $TMUX leaks into the freed outer shell's env).
const tmuxReentryEnv = "SCREENSAVERX_TMUX_REENTRY"

// handleTmux checks whether we're running inside a tmux session and, if so,
// performs a detach / re-exec / re-attach handoff so the screensaver can take
// over the real outer terminal instead of running inside a tmux pane.
//
// Returns:
//
//	true  → caller should exit immediately; the real work continues in a
//	        new copy of this binary running in the freed outer shell.
//	false → not in tmux (or sentinel already set, meaning we ARE the
//	        re-exec'd child); caller should run the screensaver inline.
func handleTmux() bool {
	if os.Getenv("TMUX") == "" || os.Getenv(tmuxReentryEnv) != "" {
		return false
	}
	if err := runViaTmuxDetach(); err != nil {
		fmt.Fprintf(os.Stderr, "tmux detach failed (%v); running inside tmux pane instead.\n", err)
		return false
	}
	return true
}

// runViaTmuxDetach asks tmux to detach the current client and exec a shell
// command in its place: re-run this binary, then `tmux attach -t <session>`.
func runViaTmuxDetach() error {
	sessionBytes, err := exec.Command("tmux", "display-message", "-p", "#S").Output()
	if err != nil {
		return fmt.Errorf("get session name: %w", err)
	}
	session := strings.TrimSpace(string(sessionBytes))
	if session == "" {
		return fmt.Errorf("empty tmux session name")
	}

	bin, err := os.Executable()
	if err != nil {
		return fmt.Errorf("locate own binary: %w", err)
	}

	cmd := fmt.Sprintf("%s=1 %s && tmux attach -t %s",
		tmuxReentryEnv, shellQuote(bin), shellQuote(session))

	return exec.Command("tmux", "detach-client", "-E", cmd).Run()
}

// shellQuote wraps s in single quotes, escaping any embedded single quotes so
// it's safe to splice into a /bin/sh command.
func shellQuote(s string) string {
	return "'" + strings.ReplaceAll(s, "'", `'\''`) + "'"
}
