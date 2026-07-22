package cmd

import (
	"errors"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/sincerefly/easycmd/common"

	"github.com/spf13/cobra"
	v "github.com/spf13/viper"
)

type cobraFunc func(cmd *cobra.Command, args []string)
type pythonFunc func(cmd *cobra.Command, args []string, data Data)

type Data struct {
	Name string
}

func python(fn pythonFunc, envName string) cobraFunc {
	return func(cmd *cobra.Command, args []string) {
		name := os.Getenv(envName)
		data := Data{
			Name: name,
		}
		fn(cmd, args, data)
	}
}

var (
	cfgFile string
)

func init() {
	cobra.OnInitialize(initConfig)

	flags := rootCmd.Flags()
	flags.BoolP("version", "v", false, "output version")
}

var rootCmd = &cobra.Command{
	Use:   "easycmd",
	Short: "A terminal tool template",
	Long:  "Long Terminal Usage desc",
	Run: python(func(cmd *cobra.Command, args []string, data Data) {
		if cmd.Flags().Lookup("version").Changed {
			common.PrintVersion()
			return
		}

		log.Println(cfgFile)
		fmt.Printf("hi,%s\n", data.Name)
	}, "NAME"),
}

func initConfig() {
	if cfgFile == "" {
		home, err := os.UserHomeDir()
		checkErr(err)
		v.AddConfigPath(".")
		v.AddConfigPath(home)
		v.AddConfigPath("/etc/easycmd/")
		v.SetConfigName("config")
		v.SetConfigType("toml")
	} else {
		v.SetConfigFile(cfgFile)
	}

	v.SetEnvPrefix("ec")
	v.AutomaticEnv()
	v.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))

	if err := v.ReadInConfig(); err != nil {
		if _, ok := errors.AsType[v.ConfigParseError](err); ok {
			log.Fatal(err)
		}
		cfgFile = "No config file used"
	} else {
		cfgFile = "Using config file: " + v.ConfigFileUsed()
	}
}
