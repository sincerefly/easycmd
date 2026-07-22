package cmd

import (
	"bytes"
	"strings"
	"testing"

	"github.com/sincerefly/easycmd/common"
	"github.com/sincerefly/easycmd/internal/testutil"
	"github.com/spf13/pflag"
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
		resetConfigForTest(t)
		resetRootCommandFlags(t)
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

func resetRootCommandFlags(t *testing.T) {
	t.Helper()
	_ = rootCmd.Flags().Set("version", "false")
	rootCmd.Flags().VisitAll(func(f *pflag.Flag) {
		f.Changed = false
	})
}

func TestRootDefaultRun(t *testing.T) {
	resetConfigForTest(t)
	resetRootCommandFlags(t)
	t.Setenv("NAME", "alice")
	t.Cleanup(func() { rootCmd.SetArgs(nil) })

	out, err := testutil.ExecuteCommand(rootCmd)
	if err != nil {
		t.Fatalf("execute root: %v", err)
	}
	if !strings.Contains(out.Out, "hi,alice") {
		t.Fatalf("unexpected stdout: %q", out.Out)
	}
}

func TestRootDefaultRunWithoutName(t *testing.T) {
	resetConfigForTest(t)
	resetRootCommandFlags(t)
	t.Setenv("NAME", "")
	t.Cleanup(func() { rootCmd.SetArgs(nil) })

	out, err := testutil.ExecuteCommand(rootCmd)
	if err != nil {
		t.Fatalf("execute root: %v", err)
	}
	if !strings.Contains(out.Out, "hi,") {
		t.Fatalf("unexpected stdout: %q", out.Out)
	}
}

func TestVersionSubcommand(t *testing.T) {
	common.Version = "3.0.0"
	common.CommitSHA = "sha"
	common.BuildDate = "date"

	withVersionOutput(t, func(buf *bytes.Buffer) {
		resetConfigForTest(t)
		_, err := testutil.ExecuteCommand(rootCmd, "version")
		if err != nil {
			t.Fatalf("execute version subcommand: %v", err)
		}
		if !strings.Contains(buf.String(), "easycmd 3.0.0") {
			t.Fatalf("unexpected output: %q", buf.String())
		}
	})
}
