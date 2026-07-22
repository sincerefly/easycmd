package cmd

import (
	"testing"

	"github.com/spf13/pflag"
	v "github.com/spf13/viper"
)

func resetViper(t *testing.T) {
	t.Helper()
	v.Reset()
}

func TestGetBoolParamB(t *testing.T) {
	flags := pflag.NewFlagSet("test", pflag.ContinueOnError)
	flags.Bool("all", false, "all")

	value, has := getBoolParamB(flags, "all")
	if has || value {
		t.Fatalf("expected default false and has=false, got value=%v has=%v", value, has)
	}

	if err := flags.Set("all", "true"); err != nil {
		t.Fatalf("set flag: %v", err)
	}
	value, has = getBoolParamB(flags, "all")
	if !has || !value {
		t.Fatalf("expected true and has=true, got value=%v has=%v", value, has)
	}
}

func TestGetStringParamB(t *testing.T) {
	resetViper(t)
	flags := pflag.NewFlagSet("test", pflag.ContinueOnError)
	flags.String("server", "", "server")

	value, has := getStringParamB(flags, "server")
	if has || value != "" {
		t.Fatalf("expected empty default and has=false, got value=%q has=%v", value, has)
	}

	if err := flags.Set("server", "http://example.com"); err != nil {
		t.Fatalf("set flag: %v", err)
	}
	value, has = getStringParamB(flags, "server")
	if !has || value != "http://example.com" {
		t.Fatalf("expected changed server, got value=%q has=%v", value, has)
	}

	resetViper(t)
	flags = pflag.NewFlagSet("test", pflag.ContinueOnError)
	flags.String("server", "", "server")
	v.Set("server", "http://from-config")
	value, has = getStringParamB(flags, "server")
	if !has || value != "http://from-config" {
		t.Fatalf("expected viper value, got value=%q has=%v", value, has)
	}
}

func TestGetParam(t *testing.T) {
	resetViper(t)
	flags := pflag.NewFlagSet("test", pflag.ContinueOnError)
	flags.String("server", "default", "server")

	if got := getParam(flags, "server"); got != "default" {
		t.Fatalf("expected default, got %q", got)
	}
}
