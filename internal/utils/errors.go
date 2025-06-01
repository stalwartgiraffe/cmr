package utils

import (
	"fmt"

	"github.com/TwiN/go-color"
)

func Redln(a ...any) (int, error) {
	return fmt.Print(color.Colorize(color.Red, fmt.Sprintln(a...)))
}
