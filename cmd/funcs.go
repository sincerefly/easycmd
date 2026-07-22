package cmd

import (
	"github.com/spf13/pflag"
	v "github.com/spf13/viper"
)

// getBoolParamB returns a bool flag value and whether it was explicitly set on the CLI.
func getBoolParamB(flags *pflag.FlagSet, key string) (bool, bool) {
	value, _ := flags.GetBool(key)
	if flags.Changed(key) {
		return value, true
	}
	return value, false
}

// getStringParamB returns a string parameter and whether it differs from the default.
func getStringParamB(flags *pflag.FlagSet, key string) (string, bool) {
	value, _ := flags.GetString(key)
	if flags.Changed(key) {
		return value, true
	}
	if v.IsSet(key) {
		return v.GetString(key), true
	}
	return value, false
}

func getParam(flags *pflag.FlagSet, key string) string {
	val, _ := getStringParamB(flags, key)
	return val
}
