package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
)

func main() {
	dirPath := os.Args[1]

	// Отримайте список файлів у папці
	files, err := os.ReadDir(dirPath)
	if err != nil {
		log.Fatal(err)
	}

	// Перетворіть список на масив рядків
	var fileNames []string = []string{"", ""}
	for _, file := range files {
		if !file.IsDir() {
			fileNames = append(fileNames, file.Name())
		}
	}
	fmt.Println(len(fileNames))

	// Розбийте файли на групи по 12
	const groupSize = 12
	var groups [][]*string

	for i := 0; i < len(fileNames); i += groupSize {
		var group []*string
		for j := 0; j < groupSize; j++ {
			if i+j < len(fileNames) {
				if fileNames[i+j] == "" {
					group = append(group, nil)
				} else {
					group = append(group, &fileNames[i+j])
				}
			} else {
				group = append(group, nil)
			}
		}
		groups = append(groups, group)
	}

	// Виведіть вміст груп
	for i, group := range groups {
		fmt.Printf("Group %d:\n", i+1)
		for idx, file := range group {
			if *file != nil {
				degree := "90"
				if idx == 1 || idx == 10 || idx == 3 || idx == 8 || idx == 5 || idx == 6 {
					degree = "270"
				}
				run(filepath.Join(os.Args[1], *file), degree, filepath.Join(os.Args[2], *file))
			}
		}
	}
}

func run(filePath, degree, outPath string) {
	cmd := exec.Command("convert", filePath, "-rotate", degree, outPath)

    // Виконайте команду
    err := cmd.Run()
    if err != nil {
        fmt.Printf("Error: %v\n", err)
        return
    }
}
