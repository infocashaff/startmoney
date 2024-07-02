package main

import (
    "encoding/json"
    "io"
    "net/http"
    "os"
    "path/filepath"
    "fmt"
)

const botToken = "YOUR_BOT_TOKEN"

type GetFileResponse struct {
    Ok     bool `json:"ok"`
    Result struct {
        FilePath string `json:"file_path"`
    } `json:"result"`
}

func getFilePath(fileID string) (string, error) {
    url := fmt.Sprintf("https://api.telegram.org/bot%s/getFile?file_id=%s", botToken, fileID)
    resp, err := http.Get(url)
    if err != nil {
        return "", err
    }
    defer resp.Body.Close()

    var getFileResp GetFileResponse
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

    out, err := os.Create(filepath.Base(filePath))
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

func main() {
    fileID := "FILE_ID_RECEIVED_FROM_USER"
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
