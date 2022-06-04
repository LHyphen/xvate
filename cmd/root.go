/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>

*/
package cmd

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"

	"example.com/xvate/app"
	"github.com/spf13/cobra"
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use: "xvate",
	// 	Short: "非对称加解密程序",
	// 	Long: `非对称加解密程序
	// 在程序后面加上文件路径即可加密生成dat文件
	// 在程序后面加上dat文件即可解密还原出原文件`,
	Short: "Asymmetric encryption and decryption program",
	Long: `Asymmetric encryption and decryption program
Add the file after the program to encrypt the generated dat file
Add the dat file after the program to decrypt and restore the original file`,
	Args: func(cmd *cobra.Command, args []string) error {
		if len(args) != 1 {
			// return errors.New("请输入需要加密或者解密的文件")
			return errors.New("please input file")
		}
		return nil
	},
	// Uncomment the following line if your bare application
	// has an action associated with it:
	Run: func(cmd *cobra.Command, args []string) {
		filename := args[0]
		if !filepath.IsAbs(filename) {
			filename, _ = filepath.Abs(filename)
		}
		err := app.Handler(filename)
		if err != nil {
			fmt.Println(err)
		}
	},
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.

	// rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.myapp.yaml)")

	// Cobra also supports local flags, which will only run
	// when this action is called directly.
	// rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
