package main

import (
	"fmt"
	"os"
	"strings"

	"prabandh/database"
	"prabandh/indexer"
	"prabandh/models"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"gorm.io/gorm"
)

type model struct {
	choices   []string
	cursor    int
	selected  string
	operation string
	message   string
	inputMode bool
	input     string
	verbose   bool
}

func initialModel(verbose bool) model {
	return model{
		choices: []string{
			"Add Directory to Index",
			"Delete Directory from Index",
			"Search Files and Summaries",
			"View Whitelisted Directories",
			"View Blacklisted Directories",
			"Toggle Verbose Mode",
			"Exit",
		},
		cursor:    0,
		selected:  "",
		operation: "",
		message:   "",
		inputMode: false,
		input:     "",
		verbose:   verbose,
	}
}

func (m model) Init() tea.Cmd {
	return nil
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		if m.inputMode {
			switch msg.String() {
			case "enter":
				m.message = m.handleOperation(m.operation, m.input)
				m.inputMode = false
				m.input = ""
				m.operation = ""
			case "esc":
				m.inputMode = false
				m.input = ""
			case "backspace":
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
			case "Toggle Verbose Mode":
				m.verbose = !m.verbose
				if m.verbose {
					m.message = "Verbose mode enabled"
				} else {
					m.message = "Verbose mode disabled"
				}
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
		m.message = m.handleOperation(m.operation, m.input)
		m.operation = ""
	}

	var b strings.Builder
	title := lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("205")).Render("Prabandh CLI")
	b.WriteString(title + "\n\n")

	if m.verbose {
		mode := lipgloss.NewStyle().Foreground(lipgloss.Color("12")).Render("(verbose)")
		b.WriteString(fmt.Sprintf("Mode: %s\n\n", mode))
	}

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

func (m model) handleOperation(operation, input string) string {
	database.Connect()

	switch operation {
	case "add_directory_to_index":
		dirPath := strings.TrimSpace(input)
		if dirPath == "" {
			return "Error: Directory path cannot be empty"
		}

		// Check if directory already exists
		var existingDir models.IndexDir
		if err := database.DB.Where("directory_location = ?", dirPath).First(&existingDir).Error; err == nil {
			return fmt.Sprintf("Directory already indexed: %s", dirPath)
		}

		// Add to IndexDir
		indexDir := models.IndexDir{DirectoryLocation: dirPath, IsWhitelisted: true}
		if err := database.DB.Create(&indexDir).Error; err != nil {
			return fmt.Sprintf("Error adding directory: %v", err)
		}

		// Index files
		indexer := indexer.NewFileIndexer("http://localhost:5051", "gemma:2b", m.verbose)
		indexer.IndexDirectory(dirPath)

		return fmt.Sprintf("Successfully indexed directory: %s", dirPath)

	case "delete_directory_from_index":
		dirPath := strings.TrimSpace(input)
		if dirPath == "" {
			return "Error: Directory path cannot be empty"
		}

		// Delete directory and its associated files
		err := database.DB.Transaction(func(tx *gorm.DB) error {
			// Find the directory record
			var dir models.IndexDir
			if err := tx.Where("directory_location = ?", dirPath).First(&dir).Error; err != nil {
				return err
			}

			// Delete associated files
			if err := tx.Where("file_path LIKE ?", dirPath+"%").Delete(&models.FileIndex{}).Error; err != nil {
				return err
			}

			// Delete the directory
			return tx.Delete(&dir).Error
		})

		if err != nil {
			return fmt.Sprintf("Error deleting directory: %v", err)
		}
		return fmt.Sprintf("Successfully deleted directory: %s", dirPath)

	case "search_files_and_summaries":
		query := strings.TrimSpace(input)
		if query == "" {
			return "Error: Search query cannot be empty"
		}

		var result strings.Builder

		// Search files by name
		var files []models.FileIndex
		database.DB.Where("file_name LIKE ?", "%"+query+"%").Find(&files)
		if len(files) > 0 {
			result.WriteString("Matching Files:\n")
			for _, file := range files {
				result.WriteString(fmt.Sprintf("- %s (%s)\n", file.FileName, file.FilePath))
			}
		} else {
			result.WriteString("No matching files found\n")
		}

		// Search by keywords
		var summaries []models.FileSummary
		database.DB.Where("summary_keyword LIKE ?", "%"+query+"%").Find(&summaries)

		if len(summaries) > 0 {
			result.WriteString("\nMatching Keywords:\n")
			fileKeywords := make(map[string][]string)
			for _, summary := range summaries {
				var file models.FileIndex
				database.DB.First(&file, summary.FileIndexID)
				fileKeywords[file.FileName] = append(fileKeywords[file.FileName], summary.SummaryKeyword)
			}

			for fileName, keywords := range fileKeywords {
				result.WriteString(fmt.Sprintf("- %s: %s\n", fileName, strings.Join(keywords, ", ")))
			}
		} else {
			result.WriteString("\nNo matching keywords found\n")
		}

		return result.String()

	case "view_whitelisted":
		var dirs []models.IndexDir
		database.DB.Where("is_whitelisted = ?", true).Find(&dirs)
		if len(dirs) == 0 {
			return "No whitelisted directories found"
		}

		var result strings.Builder
		result.WriteString("Whitelisted Directories:\n")
		for _, dir := range dirs {
			// Count files in each directory
			var count int64
			database.DB.Model(&models.FileIndex{}).Where("file_path LIKE ?", dir.DirectoryLocation+"%").Count(&count)
			result.WriteString(fmt.Sprintf("- %s (%d files)\n", dir.DirectoryLocation, count))
		}
		return result.String()

	case "view_blacklisted":
		var dirs []models.IndexDir
		database.DB.Where("is_whitelisted = ?", false).Find(&dirs)
		if len(dirs) == 0 {
			return "No blacklisted directories found"
		}

		var result strings.Builder
		result.WriteString("Blacklisted Directories:\n")
		for _, dir := range dirs {
			result.WriteString(fmt.Sprintf("- %s\n", dir.DirectoryLocation))
		}
		return result.String()
	}
	return "Invalid operation"
}

func main() {
	verbose := false // Start with verbose mode off by default
	if len(os.Args) > 1 && os.Args[1] == "-v" {
		verbose = true
	}

	p := tea.NewProgram(initialModel(verbose))
	if err := p.Start(); err != nil {
		fmt.Printf("Error starting CLI: %v\n", err)
		os.Exit(1)
	}
}
