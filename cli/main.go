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
	message   string // Added to display feedback messages
	inputMode bool   // New field to track if the program is in input mode
	input     string // New field to store user input
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
		message:   "",
		inputMode: false,
		input:     "",
	}
}

func (m model) Init() tea.Cmd {
	return nil
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		if m.inputMode {
			// Handle input mode
			switch msg.String() {
			case "enter":
				m.message = handleOperation(m.operation, m.input)
				m.inputMode = false
				m.input = ""
				m.operation = ""
			case "esc":
				m.inputMode = false
				m.input = ""
			case "backspace": // Handle backspace key
				if len(m.input) > 0 {
					m.input = m.input[:len(m.input)-1]
				}
			default:
				m.input += msg.String()
			}
			return m, nil
		}

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
			case "Add Directory to Index", "Delete Directory from Index", "Search Files and Summaries":
				m.operation = strings.ToLower(strings.ReplaceAll(m.selected, " ", "_"))
				m.inputMode = true
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
	if m.inputMode {
		return fmt.Sprintf("Enter input for %s: %s", m.operation, m.input)
	}

	if m.operation != "" {
		m.message = handleOperation(m.operation, m.input)
		m.operation = "" // Reset operation after handling
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

	if m.message != "" {
		messageStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("10"))
		b.WriteString("\n" + messageStyle.Render(m.message) + "\n")
	}

	return b.String()
}

func handleOperation(operation, input string) string {
	database.Connect()

	switch operation {
	case "add_directory_to_index":
		dirPath := strings.TrimSpace(input)
		indexDir := models.IndexDir{DirectoryLocation: dirPath, IsWhitelisted: true}
		if err := database.DB.Create(&indexDir).Error; err != nil {
			return fmt.Sprintf("Error adding directory: %v", err)
		}
		return "Directory added successfully!"
	case "delete_directory_from_index":
		dirPath := strings.TrimSpace(input)
		if err := database.DB.Where("directory_location = ?", dirPath).Delete(&models.IndexDir{}).Error; err != nil {
			return fmt.Sprintf("Error deleting directory: %v", err)
		}
		return "Directory deleted successfully!"
	case "search_files_and_summaries":
		query := strings.TrimSpace(input)
		var result strings.Builder

		// Search for files
		result.WriteString("Files:\n")
		rows, err := database.DB.Raw("SELECT file_name FROM file_indices WHERE file_name LIKE ?", "%"+query+"%").Rows()
		if err != nil {
			return fmt.Sprintf("Error searching files: %v", err)
		}
		defer rows.Close()
		for rows.Next() {
			var fileName string
			if err := rows.Scan(&fileName); err != nil {
				return fmt.Sprintf("Error reading file result: %v", err)
			}
			result.WriteString(fileName + "\n")
		}

		// Search for summaries
		result.WriteString("\nSummaries:\n")
		rows, err = database.DB.Raw("SELECT summary_keyword FROM file_summaries WHERE summary_keyword LIKE ?", "%"+query+"%").Rows()
		if err != nil {
			return fmt.Sprintf("Error searching summaries: %v", err)
		}
		defer rows.Close()
		for rows.Next() {
			var summaryKeyword string
			if err := rows.Scan(&summaryKeyword); err != nil {
				return fmt.Sprintf("Error reading summary result: %v", err)
			}
			result.WriteString(summaryKeyword + "\n")
		}

		return result.String()
	case "view_whitelisted":
		var dirs []models.IndexDir
		database.DB.Where("is_whitelisted = ?", true).Find(&dirs)
		var dirList []string
		for _, dir := range dirs {
			dirList = append(dirList, dir.DirectoryLocation)
		}
		return fmt.Sprintf("Whitelisted Directories:\n%s", strings.Join(dirList, "\n"))
	case "view_blacklisted":
		var dirs []models.IndexDir
		database.DB.Where("is_whitelisted = ?", false).Find(&dirs)
		var dirList []string
		for _, dir := range dirs {
			dirList = append(dirList, dir.DirectoryLocation)
		}
		return fmt.Sprintf("Blacklisted Directories:\n%s", strings.Join(dirList, "\n"))
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
