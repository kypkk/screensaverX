package main

import (
	"os"
	"os/signal"
	"syscall"
)

// blockExitSignals swallows every signal we're allowed to swallow, so the
// screensaver can't be killed by Ctrl+C, Ctrl+Z, Ctrl+\, `kill PID`, or a
// closed controlling terminal. The only intended way out becomes the
// in-app unlock key (see Update()).
//
// SIGKILL and SIGSTOP intentionally cannot be caught — POSIX guarantees
// those to the kernel as the always-works escape hatches. From another
// shell, `kill -9 <pid>` (SIGKILL) or `kill -STOP <pid>` (SIGSTOP) will
// always end or freeze this process, by design.
func blockExitSignals() {
	sigCh := make(chan os.Signal, 16)
	signal.Notify(sigCh,
		syscall.SIGINT,  // Ctrl+C in cooked mode, or `kill -INT`
		syscall.SIGTERM, // `kill PID` default
		syscall.SIGQUIT, // Ctrl+\ in cooked mode, or `kill -QUIT`
		syscall.SIGTSTP, // Ctrl+Z in cooked mode, or `kill -TSTP`
		syscall.SIGHUP,  // controlling terminal hung up
	)
	go func() {
		for range sigCh {
			// swallow and keep the screensaver alive
		}
	}()
}
