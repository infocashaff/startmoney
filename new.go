package main

import (
	"fmt"
	"image"
	_ "image/png"
	"io/fs"
	"log"
	"os"
	"path/filepath"

	"github.com/jung-kurt/gofpdf"
)

func main() {
	filesDir := os.Args[1]
	// Ім'я файлу PDF, який ми створимо
	pdfFileName := "output.pdf"

	// Ініціалізуємо новий PDF документ
	pdf := gofpdf.New("P", "mm", "A4", "") // "P" - портретний режим, "mm" - міліметри, "A4" - формат сторінки

	filenames, err := os.ReadDir(filesDir)
	if err != nil {
		log.Fatal(err)
	}
	filenames = append([]fs.DirEntry{nil, nil}, filenames...)

	var groups [][]fs.DirEntry
	for i := 0; i <= len(filenames); i += 12 {
		groups = append(groups, filenames[i:i+12])
	}
	pageWidth, pageHeigth := pdf.GetPageSize()

	for idx, group := range groups {
		fmt.Println(idx)
		cycle(pdf, group[0], group[11], pageWidth, pageHeigth)
		cycle(pdf, group[1], group[10], pageWidth, pageHeigth)
		cycle(pdf, group[2], group[9], pageWidth, pageHeigth)
		cycle(pdf, group[3], group[8], pageWidth, pageHeigth)
		cycle(pdf, group[4], group[7], pageWidth, pageHeigth)
		cycle(pdf, group[5], group[6], pageWidth, pageHeigth)

	}

	err = pdf.OutputFileAndClose(pdfFileName)
	if err != nil {
		log.Fatalf("Помилка при збереженні файлу PDF: %v", err)
	}

	log.Printf("PDF файл успішно створено: %s\n", pdfFileName)
}

func cycle(pdf *gofpdf.Fpdf, elem1, elem2 fs.DirEntry, pageWidth, pageHeigth float64) {
	pdf.AddPage()

	if elem1 != nil {
		width, height := imgSize(filepath.Join(os.Args[1], elem1.Name()))

		imgH1 := (float64(height) * pageWidth) / float64(width)
		pdf.Image(filepath.Join(os.Args[1], elem1.Name()), 0, pageHeigth-imgH1-15, pageWidth, 0, false, "", 0, "")
	}
	if elem2 != nil {
		pdf.Image(filepath.Join(os.Args[1], elem2.Name()), 0, 15, pageWidth, 0, false, "", 0, "")
	}
}

func imgSize(path string) (int, int) {
	file, err := os.Open(path)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	img, _, err := image.Decode(file)
	if err != nil {
		log.Fatal(err)
	}
	bounds := img.Bounds()
	return bounds.Max.X, bounds.Max.Y
}
