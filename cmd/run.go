/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var (
	Server   bool
	WorkDir  string
	Listen   string
	Config   string
)
// runCmd represents the run command
var runCmd = &cobra.Command{
	Use:   "run",
	Short: "boot up the application",
	Long: `It takes four arguments: 1. the mode of the application to run, 2. the dir of the environment to run the application in, 3. if the mode is server, the listen field is required, 4. the path of config file is required.`,
	Run: func(cmd *cobra.Command, args []string) {
		if Server && Listen == "" {
			fmt.Fprintln(os.Stderr, "Error: Listen field is required when the mode is server")
			os.Exit(1)
		} 
	},
}

func init() {
	rootCmd.AddCommand(runCmd)
	runCmd.Flags().StringVarP(&Listen, "listen", "l", "", "the address where the gin binded")
	runCmd.Flags().StringVarP(&Config, "config", "c", "", "Path to the config dir")
	runCmd.Flags().StringVarP(&WorkDir, "workDir", "w", "", "Path to the working dir")
	runCmd.Flags().BoolVarP(&Server, "server", "s", false, "the mode of the application running on")
	runCmd.MarkFlagRequired("config")
	runCmd.MarkFlagRequired("workDir")
	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// runCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// runCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
