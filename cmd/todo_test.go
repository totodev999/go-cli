package cmd

import (
	"fmt"
	"testing"
	"time"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"

	"main.go/models"
	_ "main.go/utils"
)

type MockLogger struct {
	LastMessage string
}

func (m *MockLogger) Println(v ...interface{}) {
	m.LastMessage = fmt.Sprint(v...)
}

func setupTestDB(t *testing.T) *gorm.DB {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		t.Fatalf("failed to connect database: %v", err)
	}

	err = db.AutoMigrate(&models.Todo{})
	if err != nil {
		t.Fatalf("failed to migrate database: %v", err)
	}

	return db
}

func TestCreateTodo(t *testing.T) {
	db := setupTestDB(t)
	logger := &MockLogger{}
	handler := NewTodoHandler(db, logger)

	tests := []struct {
		name        string
		title       string
		description string
		dueDate     string
		wantErr     bool
	}{
		{
			name:        "Valid Todo",
			title:       "Test Todo",
			description: "Test Description",
			dueDate:     "2024-01-01",
			wantErr:     false,
		},
		{
			name:        "Invalid Date Format",
			title:       "Test Todo",
			description: "Test Description",
			dueDate:     "invalid-date",
			wantErr:     true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := handler.CreateTodo(tt.title, tt.description, tt.dueDate)
			if (err != nil) != tt.wantErr {
				t.Errorf("CreateTodo() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestGetTodos(t *testing.T) {
	db := setupTestDB(t)
	logger := &MockLogger{}
	handler := NewTodoHandler(db, logger)

	// Create test todo
	testTodo := models.Todo{
		Title:       "Test Todo",
		Description: "Test Description",
		DueDate:     "2024/01/01",
	}
	db.Create(&testTodo)

	tests := []struct {
		name    string
		id      int
		want    int
		wantErr bool
	}{
		{
			name:    "Get All Todos",
			id:      0,
			want:    1,
			wantErr: false,
		},
		{
			name:    "Get Single Todo",
			id:      int(testTodo.ID),
			want:    1,
			wantErr: false,
		},
		{
			name:    "Get Non-existent Todo",
			id:      9999,
			want:    0,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			todos, err := handler.GetTodos(tt.id)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetTodos() error = %v, wantErr %v", err, tt.wantErr)
			}
			if !tt.wantErr && len(todos) != tt.want {
				t.Errorf("GetTodos() got %v todos, want %v", len(todos), tt.want)
			}
		})
	}
}

func TestUpdateTodo(t *testing.T) {
	db := setupTestDB(t)
	logger := &MockLogger{}
	handler := NewTodoHandler(db, logger)

	// Create test todo
	testTodo := models.Todo{
		Title:       "Original Title",
		Description: "Original Description",
		DueDate:     "2024/01/01",
	}
	db.Create(&testTodo)

	tests := []struct {
		name        string
		id          int
		title       string
		description string
		dueDate     string
		wantErr     bool
	}{
		{
			name:        "Update All Fields",
			id:          int(testTodo.ID),
			title:       "Updated Title",
			description: "Updated Description",
			dueDate:     "2024-02-02",
			wantErr:     false,
		},
		{
			name:        "Update Single Field",
			id:          int(testTodo.ID),
			title:       "New Title",
			description: "",
			dueDate:     "",
			wantErr:     false,
		},
		{
			name:        "Non-existent Todo",
			id:          9999,
			title:       "Test",
			description: "Test",
			dueDate:     "2024-01-01",
			wantErr:     true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := handler.UpdateTodo(tt.id, tt.title, tt.description, tt.dueDate)
			if (err != nil) != tt.wantErr {
				t.Errorf("UpdateTodo() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestDeleteTodo(t *testing.T) {
	db := setupTestDB(t)
	logger := &MockLogger{}
	handler := NewTodoHandler(db, logger)

	// Create test todo
	testTodo := models.Todo{
		Title:       "Test Todo",
		Description: "Test Description",
		DueDate:     "2024/01/01",
	}
	db.Create(&testTodo)

	tests := []struct {
		name    string
		id      int
		wantErr bool
	}{
		{
			name:    "Delete Existing Todo",
			id:      int(testTodo.ID),
			wantErr: false,
		},
		{
			name:    "Delete Non-existent Todo",
			id:      9999,
			wantErr: false, // GORM doesn't return error for non-existent ID in Delete
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := handler.DeleteTodo(tt.id)
			if (err != nil) != tt.wantErr {
				t.Errorf("DeleteTodo() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestDateConverter(t *testing.T) {
	tests := []struct {
		name    string
		date    string
		want    time.Time
		wantErr bool
	}{
		{
			name:    "Valid Date Format YYYY-MM-DD",
			date:    "2024-01-01",
			want:    time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC),
			wantErr: false,
		},
		{
			name:    "Valid Date Format YYYY-M-D",
			date:    "2024-1-1",
			want:    time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC),
			wantErr: false,
		},
		{
			name:    "Invalid Date Format",
			date:    "invalid-date",
			want:    time.Time{},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := dateConverter(tt.date)
			if (err != nil) != tt.wantErr {
				t.Errorf("dateConverter() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && !got.Equal(tt.want) {
				t.Errorf("dateConverter() = %v, want %v", got, tt.want)
			}
		})
	}
}
