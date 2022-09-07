/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/dashotv/tower/app"
)

// upcomingCmd represents the upcoming command
var upcomingCmd = &cobra.Command{
	Use:   "upcoming",
	Short: "task: upcoming database call",
	Long:  "task: upcoming database call",
	Run: func(cmd *cobra.Command, args []string) {
		log := app.App().Log
		c, err := app.NewConnector()
		if err != nil {
			log.Errorf("error making connector: %s", err)
			return
		}

		list, err := c.Upcoming()
		if err != nil {
			log.Errorf("error getting upcoming episodes: %s", err)
		}

		fmt.Println("list count = ", len(list))
		//for _, e := range list {
		//	fmt.Printf("%# v\n", pretty.Formatter(e))
		//}
	},
}

func init() {
	taskCmd.AddCommand(upcomingCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// upcomingCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// upcomingCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
