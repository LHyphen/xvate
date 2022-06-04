/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>

*/
package cmd

import (
	"fmt"

	"example.com/xvate/app"
	"github.com/spf13/cobra"
)

// initCmd represents the init command
var initCmd = &cobra.Command{
	Use: "init",
	// Short: "生成公私钥",
	// Long:  `生成公私钥并存放在当前路径下`,
	Short: "Generate public and private keys",
	Long:  `Generate public and private keys and store them in the ./self path`,
	Run: func(cmd *cobra.Command, args []string) {
		err := app.RsaGenKey(4096)
		if err != nil {
			fmt.Println(err)
		} else {
			fmt.Println("Generate public and private key successfully")
		}
	},
}

func init() {
	rootCmd.AddCommand(initCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// initCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// initCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
