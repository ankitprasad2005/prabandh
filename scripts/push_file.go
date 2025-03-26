package main

import (
	"fmt"
	"os"
	"prabandh/database"
	"prabandh/models"
)

func main() {
	if len(os.Args) < 3 {
		fmt.Println("Usage: push_file <file_name> <file_path>")
		return
	}

	fileName := os.Args[1]
	filePath := os.Args[2]

	database.Connect()

	file := models.File{
		Name: fileName,
		Path: filePath,
	}

	if err := database.DB.Create(&file).Error; err != nil {
		fmt.Printf("Error pushing file data: %v\n", err)
		return
	}

	fmt.Println("File data pushed successfully!")
}
