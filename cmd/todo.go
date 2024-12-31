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

func dateConverter(date string) (time.Time, error) {
	formattedDueData := strings.ReplaceAll(date, "/", "-")

	var dueDateParsed time.Time
	var err error

	dueDateParsed, err = time.Parse("2006-01-02", formattedDueData)
	if err != nil {
		dueDateParsed, err = time.Parse("2006-1-2", formattedDueData)
		if err != nil {
			return time.Time{}, err
		}
	}
	return dueDateParsed, nil

}

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
		operation, _ := cmd.Flags().GetString("operation")
		title, err1 := cmd.Flags().GetString("title")
		description, err2 := cmd.Flags().GetString("description")
		dueDate, err3 := cmd.Flags().GetString("due_date")
		id, _ := cmd.Flags().GetInt("id")

		fmt.Println(id)

		if err1 != nil || err2 != nil || err3 != nil {
			fmt.Println("Error getting flags")
			return
		}

		database.Connect()
		models.Migrate(database.DB)

		switch operation {
		case "post":
			if title == "" || dueDate == "" || description == "" {
				fmt.Println("Please provide title, description, due date")
				return
			}

			dueDateParsed, err := dateConverter(dueDate)
			if err != nil {
				fmt.Println("Error parsing due date")
				return
			}
			parsedDueDate := dueDateParsed.Format("2006/01/02")

			if err != nil {
				fmt.Println("Error parsing due date")
				return
			}

			todo := models.Todo{Title: title, Description: description, DueDate: parsedDueDate}
			database.DB.Create(&todo)
			fmt.Println("Todo created successfully!")

		case "get":
			var todos []models.Todo
			if id > 0 {
				var todo models.Todo
				database.DB.First(&todo, id)
				fmt.Printf("ID: %d\nTitle: %s\nDescription: %s\nDue Date: %s\n", todo.ID, todo.Title, todo.Description, todo.DueDate)
			} else {
				database.DB.Find(&todos)
				for _, todo := range todos {
					fmt.Printf("ID: %d\nTitle: %s\nDescription: %s\nDue Date: %s\n\n", todo.ID, todo.Title, todo.Description, todo.DueDate)
				}
			}

		case "put":
			if id == 0 {
				fmt.Println("Please provide todo ID")
				return
			}
			var todo models.Todo
			database.DB.First(&todo, id)

			if title != "" {
				todo.Title = title
			}
			if description != "" {
				todo.Description = description
			}
			if dueDate != "" {
				dueDateParsed, err := dateConverter(dueDate)
				if err != nil {
					fmt.Println("Error parsing due date")
					return
				}
				todo.DueDate = dueDateParsed.Format("2006/01/02")
			}

			database.DB.Save(&todo)
			fmt.Println("Todo updated successfully!")

		case "delete":
			if id == 0 {
				fmt.Println("Please provide todo ID")
				return
			}
			var todo models.Todo
			database.DB.Delete(&todo, id)
			fmt.Println("Todo deleted successfully!")

		default:
			fmt.Println("Invalid operation. Use post, get, put, or delete")
		}
	}}

func init() {
	rootCmd.AddCommand(todoCmd)
	todoCmd.Flags().StringP("operation", "O", "", "Operation(get, post, put, delete)")
	todoCmd.Flags().IntP("id", "I", 0, "Id of each todo. Necessary when put and delete operation")
	todoCmd.Flags().StringP("title", "T", "", "todo title. Necessary when post operation")
	todoCmd.Flags().StringP("description", "D", "", "description. Necessary when post operation")
	todoCmd.Flags().StringP("due_date", "d", "", "due_date. Necessary when post operation")

}
