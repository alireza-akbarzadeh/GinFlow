package console

import "fmt"

// ANSI color codes for terminal output
const (
	colorReset  = "\033[0m"
	colorRed    = "\033[31m"
	colorGreen  = "\033[32m"
	colorYellow = "\033[33m"
	colorBlue   = "\033[34m"
	colorPurple = "\033[35m"
	colorCyan   = "\033[36m"
	colorWhite  = "\033[37m"
	colorBold   = "\033[1m"
)

// Console provides colored terminal output for better developer experience
type Console struct{}

// New creates a new Console instance
func New() *Console {
	return &Console{}
}

// Success prints a green success message with emoji
func (c *Console) Success(emoji, message string) {
	fmt.Printf("%s%s %s%s%s\n", colorGreen, emoji, colorBold, message, colorReset)
}

// Info prints a cyan info message with emoji
func (c *Console) Info(emoji, message string) {
	fmt.Printf("%s%s %s%s%s\n", colorCyan, emoji, colorBold, message, colorReset)
}

// Warning prints a yellow warning message with emoji
func (c *Console) Warning(emoji, message string) {
	fmt.Printf("%s%s %s%s%s\n", colorYellow, emoji, colorBold, message, colorReset)
}

// Error prints a red error message with emoji
func (c *Console) Error(emoji, message string) {
	fmt.Printf("%s%s %s%s%s\n", colorRed, emoji, colorBold, message, colorReset)
}

// URL prints a labeled URL with blue highlighting
func (c *Console) URL(emoji, label, url string) {
	fmt.Printf("%s%s %s: %s%s%s\n", colorCyan, emoji, label, colorBlue, url, colorReset)
}

// Line prints a blank line
func (c *Console) Line() {
	fmt.Println()
}

// Divider prints a decorative divider line
func (c *Console) Divider() {
	fmt.Printf("%s%s══════════════════════════════════════════════════%s\n", colorGreen, colorBold, colorReset)
}
