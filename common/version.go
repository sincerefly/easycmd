package common

import (
	"fmt"
)

var (
	Version   = "(untracked)"
	CommitSHA = "(unknown)"
	BuildDate = "(unknown)"
)

func PrintVersion() {
	fmt.Printf("easycmd %s (%s %s)\n", Version, CommitSHA, BuildDate)
}
