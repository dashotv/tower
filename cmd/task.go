/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"github.com/spf13/cobra"
)

// taskCmd represents the test command
var taskCmd = &cobra.Command{
	Use:   "task",
	Short: "does nothing, parent of tasks",
	Long:  "does nothing, parent of tasks",
	//Run: func(cmd *cobra.Command, args []string) {
	//	fmt.Println("test called")
	//},
}

func init() {
	rootCmd.AddCommand(taskCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// testCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// testCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
