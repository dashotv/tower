/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"bytes"
	"encoding/json"
	"log"
	"net/http"

	"github.com/spf13/cobra"
)

// logCmd represents the log command
var logCmd = &cobra.Command{
	Use:   "log",
	Short: "test sending logs to ui",
	Long:  "test sending logs to ui",
	Run: func(cmd *cobra.Command, args []string) {
		if err := sendLog("debug", "cli", "debug message"); err != nil {
			log.Fatal(err)
		}
		if err := sendLog("info", "cli", "test message"); err != nil {
			log.Fatal(err)
		}
		if err := sendLog("warn", "cli", "warn message"); err != nil {
			log.Fatal(err)
		}
		if err := sendLog("error", "cli", "error message"); err != nil {
			log.Fatal(err)
		}
	},
}

func sendLog(level, facility, message string) error {
	json_data, err := json.Marshal(map[string]string{"level": level, "facility": facility, "message": message})
	if err != nil {
		return err
	}

	resp, err := http.Post("http://localhost:9000/messages", "application/json", bytes.NewBuffer(json_data))
	if err != nil {
		return err
	}

	return resp.Body.Close()
}

func init() {
	taskCmd.AddCommand(logCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// logCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// logCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
