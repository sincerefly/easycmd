package cmd

import (
	"testing"

	v "github.com/spf13/viper"
)

func resetConfigForTest(t *testing.T) {
	t.Helper()
	v.Reset()
	cfgFile = ""
}
