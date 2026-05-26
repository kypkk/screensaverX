# screensaverX

A terminal screensaver that's intentionally hard to dismiss — built to make you actually take the break.

## Install / Build / Run

Install the latest version from GitHub:

```bash
go install github.com/kypkk/screensaverX@latest
screensaverX
```

Or build from a local clone:

```bash
go build -o screensaverX ./...
./screensaverX
```

Requires Go 1.21+.

## Customizing the exit key

At the top of `main.go`:

```go
// exitKey is the only key combination that ends the screensaver.
const exitKey = "ctrl+q"
```

Change this constant to any string bubbletea recognizes, then rebuild:

| What you want | String |
|---|---|
| Ctrl + letter | `"ctrl+a"`, `"ctrl+q"`, `"ctrl+x"`, … |
| Function keys | `"f1"`, `"f2"`, … `"f12"` |
| Alt + letter | `"alt+a"`, `"alt+x"`, … |
| Escape | `"esc"` |
| Arrow keys | `"up"`, `"down"`, `"left"`, `"right"` |
| Other special keys | `"enter"`, `"tab"`, `"space"`, `"backspace"`, `"home"`, `"end"`, `"pgup"`, `"pgdown"` |

Note: most terminals collapse `Ctrl+<letter>` to a single control character and drop the Shift bit, so `"ctrl+shift+q"` only works in terminals that support the [kitty keyboard protocol](https://sw.kovidgoyal.net/kitty/keyboard-protocol/) (Kitty, WezTerm, Ghostty, etc.). `Cmd+<letter>` on macOS is typically intercepted by the terminal app before reaching the program.

## License

MIT
