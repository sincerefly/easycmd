package common

import (
	"bytes"
	"strings"
	"testing"
)

func TestPrintVersion(t *testing.T) {
	oldVersion := Version
	oldCommit := CommitSHA
	oldBuildDate := BuildDate
	oldWriter := PrintWriter()
	t.Cleanup(func() {
		Version = oldVersion
		CommitSHA = oldCommit
		BuildDate = oldBuildDate
		SetPrintWriter(oldWriter)
	})

	Version = "0.1.0"
	CommitSHA = "abc123"
	BuildDate = "2026-07-22"

	buf := new(bytes.Buffer)
	SetPrintWriter(buf)
	PrintVersion()

	out := buf.String()
	if !strings.Contains(out, "easycmd 0.1.0") {
		t.Fatalf("unexpected output: %q", out)
	}
	if !strings.Contains(out, "abc123") || !strings.Contains(out, "2026-07-22") {
		t.Fatalf("missing commit or build date in output: %q", out)
	}
}
