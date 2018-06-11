package api

import (
	"fmt"
	"os"
)

func Fatal(msg string, err error) {
	fmt.Fprintln(os.Stderr, fmt.Sprintf("%s: %s", msg, err))

	os.Exit(1)
}
