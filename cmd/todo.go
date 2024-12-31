/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"strings"
	"time"

	"github.com/spf13/cobra"
	"gorm.io/gorm"

	"main.go/database"
	"main.go/models"
	"main.go/utils"
)

type TodoHandler struct {
	db     *gorm.DB
	logger utils.Logger
}

func NewTodoHandler(db *gorm.DB, logger utils.Logger) *TodoHandler {
	return &TodoHandler{
		db:     db,
		logger: logger,
	}
}

func (h *TodoHandler) CreateTodo(title, description, dueDate string) error {
	dueDateParsed, err := dateConverter(dueDate)
	if err != nil {
		return fmt.Errorf("error parsing due date: %w", err)
	}

	todo := models.Todo{
		Title:       title,
		Description: description,
		DueDate:     dueDateParsed.Format("2006/01/02"),
	}
	return h.db.Create(&todo).Error
}

func (h *TodoHandler) GetTodos(id int) ([]models.Todo, error) {
	if id > 0 {
		var todo models.Todo
		if err := h.db.First(&todo, id).Error; err != nil {
			return nil, err
		}
		return []models.Todo{todo}, nil
	}

	var todos []models.Todo
	if err := h.db.Find(&todos).Error; err != nil {
		return nil, err
	}
	return todos, nil
}

func (h *TodoHandler) UpdateTodo(id int, title, description, dueDate string) error {
	var todo models.Todo
	if err := h.db.First(&todo, id).Error; err != nil {
		return err
	}

	if title != "" {
		todo.Title = title
	}
	if description != "" {
		todo.Description = description
	}
	if dueDate != "" {
		dueDateParsed, err := dateConverter(dueDate)
		if err != nil {
			return err
		}
		todo.DueDate = dueDateParsed.Format("2006/01/02")
	}

	return h.db.Save(&todo).Error
}

func (h *TodoHandler) DeleteTodo(id int) error {
	return h.db.Delete(&models.Todo{}, id).Error
}

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

var todoCmd = &cobra.Command{
	Use:   "todo",
	Short: "Todo management command",
	Run: func(cmd *cobra.Command, args []string) {
		database.Connect()
		models.Migrate(database.DB)

		handler := NewTodoHandler(database.DB, &utils.DefaultLogger{})

		operation, _ := cmd.Flags().GetString("operation")
		title, _ := cmd.Flags().GetString("title")
		description, _ := cmd.Flags().GetString("description")
		dueDate, _ := cmd.Flags().GetString("due_date")
		id, _ := cmd.Flags().GetInt("id")

		switch operation {
		case "post":
			if title == "" || dueDate == "" || description == "" {
				handler.logger.Println("Please provide title, description, due date")
				return
			}
			if err := handler.CreateTodo(title, description, dueDate); err != nil {
				handler.logger.Println("Error creating todo:", err)
				return
			}
			handler.logger.Println("Todo created successfully!")

		case "get":
			todos, err := handler.GetTodos(id)
			if err != nil {
				handler.logger.Println("Error getting todos:", err)
				return
			}
			for _, todo := range todos {
				handler.logger.Println("ID:", todo.ID)
				handler.logger.Println("Title:", todo.Title)
				handler.logger.Println("Description:", todo.Description)
				handler.logger.Println("Due Date:", todo.DueDate)
				handler.logger.Println()
			}

		case "put":
			if id == 0 {
				handler.logger.Println("Please provide todo ID")
				return
			}
			if err := handler.UpdateTodo(id, title, description, dueDate); err != nil {
				handler.logger.Println("Error updating todo:", err)
				return
			}
			handler.logger.Println("Todo updated successfully!")

		case "delete":
			if id == 0 {
				handler.logger.Println("Please provide todo ID")
				return
			}
			if err := handler.DeleteTodo(id); err != nil {
				handler.logger.Println("Error deleting todo:", err)
				return
			}
			handler.logger.Println("Todo deleted successfully!")

		default:
			handler.logger.Println("Invalid operation. Use post, get, put, or delete")
		}
	},
}

func init() {
	rootCmd.AddCommand(todoCmd)
	todoCmd.Flags().StringP("operation", "O", "", "Operation(get, post, put, delete)")
	todoCmd.Flags().IntP("id", "I", 0, "Id of each todo. Necessary when put and delete operation")
	todoCmd.Flags().StringP("title", "T", "", "todo title. Necessary when post operation")
	todoCmd.Flags().StringP("description", "D", "", "description. Necessary when post operation")
	todoCmd.Flags().StringP("due_date", "d", "", "due_date. Necessary when post operation")

}
