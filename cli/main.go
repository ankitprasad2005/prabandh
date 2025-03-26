package main

import (
	"fmt"
	"os"
	"strings"

	"prabandh/database"
	"prabandh/models"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type model struct {
	choices   []string
	cursor    int
	selected  string
	operation string
}

func initialModel() model {
	return model{
		choices: []string{
			"Add Directory to Index",
			"Delete Directory from Index",
			"Search Files and Summaries",
			"View Whitelisted Directories",
			"View Blacklisted Directories",
			"Exit",
		},
		cursor:    0,
		selected:  "",
		operation: "",
	}
}

func (m model) Init() tea.Cmd {
	return nil
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "q", "ctrl+c":
			return m, tea.Quit
		case "up", "k":
			if m.cursor > 0 {
				m.cursor--
			}
		case "down", "j":
			if m.cursor < len(m.choices)-1 {
				m.cursor++
			}
		case "enter":
			m.selected = m.choices[m.cursor]
			switch m.selected {
			case "Add Directory to Index":
				m.operation = "add"
			case "Delete Directory from Index":
				m.operation = "delete"
			case "Search Files and Summaries":
				m.operation = "search"
			case "View Whitelisted Directories":
				m.operation = "view_whitelisted"
			case "View Blacklisted Directories":
				m.operation = "view_blacklisted"
			case "Exit":
				return m, tea.Quit
			}
		}
	}
	return m, nil
}

func (m model) View() string {
	if m.operation != "" {
		return handleOperation(m.operation)
	}

	var b strings.Builder
	title := lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("205")).Render("Prabandh CLI")
	b.WriteString(title + "\n\n")

	for i, choice := range m.choices {
		cursor := " "
		if m.cursor == i {
			cursor = ">"
		}
		b.WriteString(fmt.Sprintf("%s %s\n", cursor, choice))
	}
	return b.String()
}

func handleOperation(operation string) string {
	database.Connect()
	switch operation {
	case "add":
		fmt.Println("Enter directory path to add:")
		var dirPath string
		fmt.Scanln(&dirPath)
		indexDir := models.IndexDir{DirectoryLocation: dirPath, IsWhitelisted: true}
		if err := database.DB.Create(&indexDir).Error; err != nil {
			return fmt.Sprintf("Error adding directory: %v", err)
		}
		return "Directory added successfully!"
	case "delete":
		fmt.Println("Enter directory path to delete:")
		var dirPath string
		fmt.Scanln(&dirPath)
		if err := database.DB.Where("directory_location = ?", dirPath).Delete(&models.IndexDir{}).Error; err != nil {
			return fmt.Sprintf("Error deleting directory: %v", err)
		}
		return "Directory deleted successfully!"
	case "search":
		fmt.Println("Enter search query:")
		var query string
		fmt.Scanln(&query)
		var files []models.FileIndex
		var summaries []models.FileSummary
		database.DB.Where("file_name LIKE ?", "%"+query+"%").Find(&files)
		database.DB.Where("summary_keyword LIKE ?", "%"+query+"%").Find(&summaries)
		return fmt.Sprintf("Files: %v\nSummaries: %v", files, summaries)
	case "view_whitelisted":
		var dirs []models.IndexDir
		database.DB.Where("is_whitelisted = ?", true).Find(&dirs)
		return fmt.Sprintf("Whitelisted Directories: %v", dirs)
	case "view_blacklisted":
		var dirs []models.IndexDir
		database.DB.Where("is_whitelisted = ?", false).Find(&dirs)
		return fmt.Sprintf("Blacklisted Directories: %v", dirs)
	}
	return "Invalid operation"
}

func main() {
	p := tea.NewProgram(initialModel())
	if err := p.Start(); err != nil {
		fmt.Printf("Error starting CLI: %v\n", err)
		os.Exit(1)
	}
}
