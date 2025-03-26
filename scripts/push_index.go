package main

import (
	"fmt"
	"os"
	"prabandh/database"
	"prabandh/models"
)

func main() {
	if len(os.Args) < 3 {
		fmt.Println("Usage: push_index <index_key> <index_value>")
		return
	}

	indexKey := os.Args[1]
	indexValue := os.Args[2]

	database.Connect()

	index := models.Index{
		Key:   indexKey,
		Value: indexValue,
	}

	if err := database.DB.Create(&index).Error; err != nil {
		fmt.Printf("Error pushing index data: %v\n", err)
		return
	}

	fmt.Println("Index data pushed successfully!")
}
