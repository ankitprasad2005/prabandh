package main

import (
	"fmt"
	"os"
	"prabandh/database"
	"prabandh/models"
)

func main() {
	if len(os.Args) < 3 {
		fmt.Println("Usage: push_summary <summary_title> <summary_content>")
		return
	}

	summaryTitle := os.Args[1]
	summaryContent := os.Args[2]

	database.Connect()

	summary := models.Summary{
		Title:   summaryTitle,
		Content: summaryContent,
	}

	if err := database.DB.Create(&summary).Error; err != nil {
		fmt.Printf("Error pushing summary data: %v\n", err)
		return
	}

	fmt.Println("Summary data pushed successfully!")
}
