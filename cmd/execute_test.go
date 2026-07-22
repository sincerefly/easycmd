package cmd

import (
	"testing"
)

func TestRootCommandUnknownSubcommand(t *testing.T) {
	resetConfigForTest(t)
	t.Cleanup(func() {
		rootCmd.SetArgs(nil)
	})

	rootCmd.SetArgs([]string{"unknown-subcommand"})
	if err := rootCmd.Execute(); err == nil {
		t.Fatal("expected error for unknown subcommand")
	}
}

func TestCheckErr(t *testing.T) {
	t.Run("nil error", func(t *testing.T) {
		checkErr(nil)
	})
}
