/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"go-hj-hospital/HospitalQueue/app"
	"go-hj-hospital/config"
	"go-hj-hospital/util"
	"os"

	"github.com/spf13/cobra"
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "go-hj-hospital",
	Short: "A brief description of your application",
	Long: `A longer description that spans multiple lines and likely contains
	examples and usage of using your application.`,
	Run: start,
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func start(cmd *cobra.Command, args []string) {
	config.Init()
	util.CreateDB()
	util.MigrateDB(util.Master())
	app.IrisInit()
	app.IrisStart()
}
