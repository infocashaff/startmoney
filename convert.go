package main

import (
	"fmt"
	"log"
	"os"
)

func main() {
	dirPath := os.Args[1]

	// Отримайте список файлів у папці
	files, err := os.ReadDir(dirPath)
	if err != nil {
		log.Fatal(err)
	}

	// Перетворіть список на масив рядків
	var fileNames []string = []string{""}
	for _, file := range files {
		if !file.IsDir() {
			fileNames = append(fileNames, file.Name())
		}
	}

	// Розбийте файли на групи по 12
	const groupSize = 12
	var groups [][]*string

	for i := 0; i < len(fileNames); i += groupSize {
		var group []*string
		for j := 0; j < groupSize; j++ {
			if i+j < len(fileNames) {
				group = append(group, &fileNames[i+j])
			}
			
			if fileNames[i+j] == "" {
				group = append(group, nil)
			} else {
				group = append(group, nil)
			}
		}
		groups = append(groups, group)
	}

	// Виведіть вміст груп
	for i, group := range groups {
		fmt.Printf("Group %d:\n", i+1)
		for _, file := range group {
			fmt.Println(file)
		}
	}
}
