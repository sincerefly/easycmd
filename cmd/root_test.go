package cmd

import (
	"bytes"
	"strings"
	"testing"

	"github.com/sincerefly/easycmd/common"
)

func withVersionOutput(t *testing.T, fn func(buf *bytes.Buffer)) {
	t.Helper()

	oldWriter := common.PrintWriter()
	buf := new(bytes.Buffer)
	common.SetPrintWriter(buf)
	t.Cleanup(func() { common.SetPrintWriter(oldWriter) })

	fn(buf)
}

func TestVersionCommand(t *testing.T) {
	common.Version = "1.2.3"
	common.CommitSHA = "abc1234"
	common.BuildDate = "2026-07-22"

	withVersionOutput(t, func(buf *bytes.Buffer) {
		versionCmd.Run(versionCmd, []string{})

		out := buf.String()
		if !strings.Contains(out, "easycmd 1.2.3") {
			t.Fatalf("unexpected output: %q", out)
		}
		if !strings.Contains(out, "abc1234") || !strings.Contains(out, "2026-07-22") {
			t.Fatalf("missing commit or build date in output: %q", out)
		}
	})
}

func TestRootVersionFlag(t *testing.T) {
	common.Version = "9.9.9"
	common.CommitSHA = "deadbeef"
	common.BuildDate = "2026-01-01"

	withVersionOutput(t, func(buf *bytes.Buffer) {
		rootCmd.SetArgs([]string{"--version"})
		if err := rootCmd.Execute(); err != nil {
			t.Fatalf("execute root --version: %v", err)
		}

		out := buf.String()
		if !strings.Contains(out, "easycmd 9.9.9") {
			t.Fatalf("unexpected output: %q", out)
		}
	})
}
