/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"strings"
	"time"

	"github.com/spf13/cobra"
	"main.go/database"
	"main.go/models"
)

// todoCmd represents the todo command
var todoCmd = &cobra.Command{
	Use:   "todo",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		title, err1 := cmd.Flags().GetString("title")
		description, err2 := cmd.Flags().GetString("description")
		dueDate, err3 := cmd.Flags().GetString("due_date")

		if err1 != nil || err2 != nil || err3 != nil {
			fmt.Println("Error getting flags")
			return
		}

		if title == "" || dueDate == "" || description == "" {
			fmt.Println("Please provide title, description, due date")
			return
		}

		// replace slash to hyphen in dueDate
		formattedDueData := strings.ReplaceAll(dueDate, "/", "-")
		fmt.Println("formattedDueData: ", formattedDueData)
		// Convert dueDate to time.Time
		dueDateParsed, err := time.Parse("2006-01-02", formattedDueData)
		parsedDueDate := dueDateParsed.Format("2006/01/02")

		if err != nil {
			fmt.Println("Error parsing due date")
			return
		}

		database.Connect()
		models.Migrate(database.DB)

		todo := models.Todo{Title: title, Description: description, DueDate: parsedDueDate}
		database.DB.Create(&todo)

		fmt.Println("Todo created successfully!")

	}}

func init() {
	rootCmd.AddCommand(todoCmd)
	todoCmd.Flags().StringP("title", "T", "", "todo title")
	todoCmd.Flags().StringP("description", "D", "", "description")
	todoCmd.Flags().StringP("due_date", "d", "", "due_date")

}
