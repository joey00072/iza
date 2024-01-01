/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"

	pkg "github.com/joey00072/iza/pkg"
	types "github.com/joey00072/iza/pkg/types"
	"github.com/spf13/cobra"
)

// type PullImageOptions struct {
//     ImageName   string
//     TarballPath string
//     CachePath   string // Optional, path to cache the image
// }
// PullImageOptions holds the options for pulling an image.

// pkg.PullImageOptions on pkg.go.dev

// pullCmd represents the pull command
var pullCmd = &cobra.Command{
	Use:   "pull",
	Short: "Pull Continer Image (eg. ubuntu, alpine, etc.)",
	Long: `Pull images from a registry (only docker.io supported currently):
	Example: iza pull ubuntu
	`,
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) == 0 {
			fmt.Println("Error: Please specify the image name")
			fmt.Println("Example: iza pull alpine")
			return
		}
		imageOpt := types.ImageOptions{
			ImageName:   args[0],
			TarballPath: args[0] + ".tar",
		}
		pkg.PullImage(imageOpt)
	},
}

func init() {
	rootCmd.AddCommand(pullCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// pullCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// pullCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
