package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
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

var rootCmd = &cobra.Command{

	Use:   "easycmd",
	Short: "A terminal tool template",
	Long:  "Long Terminal Usage desc",

	// Run: func(cmd *cobra.Command, args []string) {
	//   fmt.Println("just run")
	// },

	Run: python(func(cmd *cobra.Command, args []string, data Data) {
		fmt.Printf("hi,%s\n", data.Name)
	}, "NAME"),
}
