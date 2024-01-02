/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"strings"

	"github.com/joey00072/iza/pkg"
	types "github.com/joey00072/iza/pkg/types"
	"github.com/spf13/cobra"
)

// runCmd represents the run command
var runCmd = &cobra.Command{
	Use:   "run",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) == 0 {
			fmt.Println("Error: Please specify the image name")
			fmt.Println("Example: iza pull alpine")
			return
		}
		tarPath := strings.ReplaceAll(args[0], "/", "-") + ".tar"
		imageOpt := types.ImageOptions{
			ImageName:   args[0],
			TarballPath: tarPath,
		}
		err := pkg.Run(imageOpt, args)
		if err != nil {
			fmt.Printf("Error Running Container: %v \n", err)
			return
		}

	},
}

func init() {
	rootCmd.AddCommand(runCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// runCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// runCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
