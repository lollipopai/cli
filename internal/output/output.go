package output

import (
	"encoding/json"
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/fatih/color"
)

var (
	infoPrefix    = color.New(color.FgBlue).Sprint(">")
	successPrefix = color.New(color.FgGreen).Sprint(">")
	warnPrefix    = color.New(color.FgYellow).Sprint("!")
	errorPrefix   = color.New(color.FgRed).Sprint("!")

	cyanColor    = color.New(color.FgCyan)
	greenColor   = color.New(color.FgGreen)
	yellowColor  = color.New(color.FgYellow)
	magentaColor = color.New(color.FgMagenta)
	dimColor     = color.New(color.Faint)
	boldColor    = color.New(color.Bold)
)

func Info(msg string) {
	fmt.Fprintf(os.Stdout, "%s %s\n", infoPrefix, msg)
}

func Success(msg string) {
	fmt.Fprintf(os.Stdout, "%s %s\n", successPrefix, msg)
}

func Warn(msg string) {
	fmt.Fprintf(os.Stderr, "%s %s\n", warnPrefix, msg)
}

func Error(msg string) {
	fmt.Fprintf(os.Stderr, "%s %s\n", errorPrefix, msg)
}

func Fatal(msg string) {
	Error(msg)
	os.Exit(1)
}

func Bold(s string) string {
	return boldColor.Sprint(s)
}

func Dim(s string) string {
	return dimColor.Sprint(s)
}

func IsTTY() bool {
	return !color.NoColor
}

// PrintJSON pretty-prints a value as colored JSON to stdout.
func PrintJSON(v any) {
	data, err := json.MarshalIndent(v, "", "  ")
	if err != nil {
		fmt.Fprintln(os.Stdout, v)
		return
	}
	raw := string(data)

	if !IsTTY() {
		fmt.Fprintln(os.Stdout, raw)
		return
	}

	fmt.Fprintln(os.Stdout, colorizeJSON(raw))
}

func colorizeJSON(raw string) string {
	lines := strings.Split(raw, "\n")
	result := make([]string, 0, len(lines))
	for _, line := range lines {
		stripped := strings.TrimLeft(line, " \t")
		leading := line[:len(line)-len(stripped)]

		if strings.HasPrefix(stripped, "\"") && strings.Contains(stripped, "\": ") {
			keyEnd := strings.Index(stripped, "\": ")
			key := stripped[:keyEnd+1]
			rest := stripped[keyEnd+2:]
			result = append(result, leading+cyanColor.Sprint(key)+":"+colorizeValue(rest))
		} else {
			result = append(result, leading+colorizeValue(stripped))
		}
	}
	return strings.Join(result, "\n")
}

func colorizeValue(s string) string {
	stripped := strings.TrimRight(s, ",")
	trailing := s[len(stripped):]

	switch stripped {
	case "{", "}", "[", "]", "{}", "[]":
		return s
	case "null":
		return " " + dimColor.Sprint("null") + trailing
	case "true", "false":
		return " " + yellowColor.Sprint(stripped) + trailing
	}

	if strings.HasPrefix(stripped, "\"") {
		return " " + greenColor.Sprint(stripped) + trailing
	}

	if _, err := strconv.ParseFloat(stripped, 64); err == nil && stripped != "" {
		return " " + magentaColor.Sprint(stripped) + trailing
	}

	return " " + s
}
