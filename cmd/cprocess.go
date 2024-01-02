/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"

	pkg "github.com/joey00072/iza/pkg"
	"github.com/spf13/cobra"
)

// cprocessCmd represents the cprocess command
var cprocessCmd = &cobra.Command{
	Use:   "cprocess",
	Short: "runtime process fork (not for use cammand)",
	Long: `This is used to Setup the container runtime process.
	eg setting up hostname, enviourment variables etc.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Running cprocess...")
		pkg.RunContainer(args)
	},
}

func init() {
	rootCmd.AddCommand(cprocessCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// cprocessCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// cprocessCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
