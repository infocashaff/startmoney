package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func main() {
	bot, err := tgbotapi.NewBotAPI("7272511680:AAFGlygKMXjVv2shPTw5f7gTrs9WOfAMdzI")
	if err != nil {
		fmt.Println("Failed to create bot:", err)
		return
	}

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates, err := bot.GetUpdates(u)
	if err != nil {
		log.Fatal(err)
	}

	for _, update := range updates {
		if update.Message != nil { // Check if we got a message
			if update.Message.Document != nil {
				fileID := update.Message.Document.FileID
				file, err := bot.GetFile(tgbotapi.FileConfig{FileID: fileID})
				if err != nil {
					fmt.Println("Failed to get file:", err)
					continue
				}

				fileURL := fmt.Sprintf("https://api.telegram.org/file/bot%s/%s", bot.Token, file.FilePath)
				err = downloadFile(fileURL, update.Message.Document.FileName)
				if err != nil {
					fmt.Println("Failed to download file:", err)
				} else {
					fmt.Println("File downloaded successfully:", update.Message.Document.FileName)
				}
			}
		}
	}
}

func downloadFile(url string, fileName string) error {
	// Get the data
	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	// Create the file
	out, err := os.Create(fileName)
	if err != nil {
		return err
	}
	defer out.Close()

	// Write the body to file
	_, err = io.Copy(out, resp.Body)
	return err
}
