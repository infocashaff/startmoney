package main

import (
    "encoding/json"
    "fmt"
    "net/http"
    "os"
)

const botToken = "7272511680:AAFGlygKMXjVv2shPTw5f7gTrs9WOfAMdzI"

type Update struct {
    UpdateID int `json:"update_id"`
    Message  struct {
        MessageID int `json:"message_id"`
        From      struct {
            ID        int    `json:"id"`
            IsBot     bool   `json:"is_bot"`
            FirstName string `json:"first_name"`
            Username  string `json:"username"`
            LanguageCode string `json:"language_code"`
        } `json:"from"`
        Chat struct {
            ID        int    `json:"id"`
            FirstName string `json:"first_name"`
            Username  string `json:"username"`
            Type      string `json:"type"`
        } `json:"chat"`
        Date int `json:"date"`
        Document struct {
            FileID   string `json:"file_id"`
            FileName string `json:"file_name"`
            MimeType string `json:"mime_type"`
            FileSize int    `json:"file_size"`
        } `json:"document"`
    } `json:"message"`
}

type GetUpdatesResponse struct {
    Ok     bool     `json:"ok"`
    Result []Update `json:"result"`
}

func getUpdates() ([]Update, error) {
    url := fmt.Sprintf("https://api.telegram.org/bot%s/getUpdates", botToken)
    resp, err := http.Get(url)
    if err != nil {
        return nil, err
    }
    defer resp.Body.Close()

    var updatesResp GetUpdatesResponse
    if err := json.NewDecoder(resp.Body).Decode(&updatesResp); err != nil {
        return nil, err
    }

    if !updatesResp.Ok {
        return nil, fmt.Errorf("telegram API returned not ok")
    }

    return updatesResp.Result, nil
}

func main() {
    updates, err := getUpdates()
    if err != nil {
        fmt.Printf("Error getting updates: %v\n", err)
        return
    }

    for _, update := range updates {
        if update.Message.Document.FileID != "" {
            fileID := update.Message.Document.FileID
            fileName := update.Message.Document.FileName
            fmt.Printf("Received file with ID: %s and name: %s\n", fileID, fileName)

            // Виклик функцій для завантаження файлу
            filePath, err := getFilePath(fileID)
            if err != nil {
                fmt.Printf("Error getting file path: %v\n", err)
                return
            }

            if err := downloadFile(filePath); err != nil {
                fmt.Printf("Error downloading file: %v\n", err)
                return
            }

            fmt.Println("File downloaded successfully")
        }
    }
}

func getFilePath(fileID string) (string, error) {
    url := fmt.Sprintf("https://api.telegram.org/bot%s/getFile?file_id=%s", botToken, fileID)
    resp, err := http.Get(url)
    if err != nil {
        return "", err
    }
    defer resp.Body.Close()

    var getFileResp struct {
        Ok     bool `json:"ok"`
        Result struct {
            FilePath string `json:"file_path"`
        } `json:"result"`
    }
    if err := json.NewDecoder(resp.Body).Decode(&getFileResp); err != nil {
        return "", err
    }

    if !getFileResp.Ok {
        return "", fmt.Errorf("telegram API returned not ok")
    }

    return getFileResp.Result.FilePath, nil
}

func downloadFile(filePath string) error {
    url := fmt.Sprintf("https://api.telegram.org/file/bot%s/%s", botToken, filePath)
    resp, err := http.Get(url)
    if err != nil {
        return err
    }
    defer resp.Body.Close()

    out, err := os.Create(filePath)
    if err != nil {
        return err
    }
    defer out.Close()

    _, err = io.Copy(out, resp.Body)
    if err != nil {
        return err
    }

    return nil
}
