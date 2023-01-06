package cmd

import (
	"easycmd/version"
	"fmt"
	"log"
	"strings"

	"github.com/mitchellh/go-homedir"
	"github.com/spf13/cobra"
	v "github.com/spf13/viper"
	"gitlab.com/avarf/getenvs"
)

type cobraFunc func(cmd *cobra.Command, args []string)
type pythonFunc func(cmd *cobra.Command, args []string, data Data)

type Data struct {
	Name string
}

func python(fn pythonFunc, envName string) cobraFunc {
	return func(cmd *cobra.Command, args []string) {
		name := getenvs.GetEnvString(envName, "dong")
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
	flags.BoolP("version", "v", false, "version")
}

var rootCmd = &cobra.Command{

	Use:   "easycmd",
	Short: "A terminal tool template",
	Long:  "Long Terminal Usage desc",

	// Run: func(cmd *cobra.Command, args []string) {
	//   fmt.Println("just run")
	// },

	Run: python(func(cmd *cobra.Command, args []string, data Data) {
		log.Println(cfgFile)
		fmt.Printf("hi,%s\n", data.Name)

		if cmd.Flags().Lookup("version").Changed {
			fmt.Println("Easycmd v" + version.Version + "/" + version.CommitSHA)
			return
		}
	}, "NAME"),
}

func initConfig() {
	if cfgFile == "" {
		home, err := homedir.Dir()

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
		if _, ok := err.(v.ConfigParseError); ok {
			panic(err)
		}
		cfgFile = "No config file used"
	} else {
		cfgFile = "Using config file: " + v.ConfigFileUsed()
	}
}
