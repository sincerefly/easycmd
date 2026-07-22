package common

import (
	"fmt"
	"io"
	"os"
)

var (
	Version   = "(untracked)"
	CommitSHA = "(unknown)"
	BuildDate = "(unknown)"

	printWriter io.Writer = os.Stdout
)

func PrintVersion() {
	fmt.Fprintf(printWriter, "easycmd %s (%s %s)\n", Version, CommitSHA, BuildDate)
}

func SetPrintWriter(w io.Writer) {
	printWriter = w
}

func PrintWriter() io.Writer {
	return printWriter
}
