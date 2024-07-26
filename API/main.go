package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"time"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/gmail/v1"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
)

// Получить клиента HTTP с OAuth 2.0
func getClient(config *oauth2.Config) *http.Client {
	tokFile := "token.json"
	tok, err := tokenFromFile(tokFile)
	if err != nil {
		tok = getTokenFromWeb(config)
		saveToken(tokFile, tok)
	}
	return config.Client(context.Background(), tok)
}

// Получить токен из файла
func tokenFromFile(file string) (*oauth2.Token, error) {
	f, err := os.Open(file)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	tok := &oauth2.Token{}
	err = json.NewDecoder(f).Decode(tok)
	return tok, err
}

// Получить новый токен с помощью веб-браузера
func getTokenFromWeb(config *oauth2.Config) *oauth2.Token {
	authURL := config.AuthCodeURL("state-token", oauth2.AccessTypeOffline)
	fmt.Printf("Перейдите по следующей ссылке и введите код: \n%v\n", authURL)

	var authCode string
	if _, err := fmt.Scan(&authCode); err != nil {
		log.Fatalf("Не удалось прочитать авторизационный код: %v", err)
	}

	tok, err := config.Exchange(context.TODO(), authCode)
	if err != nil {
		log.Fatalf("Не удалось получить токен: %v", err)
	}
	return tok
}

// Сохранить токен в файл
func saveToken(path string, token *oauth2.Token) {
	fmt.Printf("Сохранение токена в файл %s\n", path)
	f, err := os.Create(path)
	if err != nil {
		log.Fatalf("Не удалось создать файл токена: %v", err)
	}
	defer f.Close()
	json.NewEncoder(f).Encode(token)
}

func main() {
	b, err := ioutil.ReadFile("key.json")
	if err != nil {
		log.Fatalf("Не удалось прочитать файл учетных данных: %v", err)
	}

	config, err := google.ConfigFromJSON(b, gmail.GmailReadonlyScope)
	if err != nil {
		log.Fatalf("Не удалось получить конфигурацию клиента: %v", err)
	}
	client := getClient(config)

	srv, err := gmail.New(client)
	if err != nil {
		log.Fatalf("Не удалось создать службу Gmail: %v", err)
	}

	// Идентификатор вашего Telegram бота
	bot, err := tgbotapi.NewBotAPI("YOUR_TELEGRAM_BOT_TOKEN")
	if err != nil {
		log.Fatalf("Не удалось создать Telegram бот: %v", err)
	}

	// ID вашего чата в Telegram
	chatID := int64(73737)

	// Бесконечный цикл для проверки новых писем
	for {
		checkForNewMessages(srv, bot, chatID)
		time.Sleep(10 * time.Second) // Задержка перед следующей проверкой
	}
}

func checkForNewMessages(srv *gmail.Service, bot *tgbotapi.BotAPI, chatID int64) {
	user := "me"
	r, err := srv.Users.Messages.List(user).LabelIds("INBOX").Q("is:unread").Do()
	if err != nil {
		log.Printf("Не удалось получить сообщения: %v", err)
		return
	}

	if len(r.Messages) > 0 {
		for _, msg := range r.Messages {
			m, err := srv.Users.Messages.Get(user, msg.Id).Do()
			if err != nil {
				log.Printf("Не удалось получить сообщение: %v", err)
				continue
			}

			sendTelegramMessage(bot, chatID, m.Snippet)

			// Пометить сообщение как прочитанное
			mod := &gmail.ModifyMessageRequest{RemoveLabelIds: []string{"UNREAD"}}
			_, err = srv.Users.Messages.Modify(user, msg.Id, mod).Do()
			if err != nil {
				log.Printf("Не удалось пометить сообщение как прочитанное: %v", err)
			}
		}
	}
}

func sendTelegramMessage(bot *tgbotapi.BotAPI, chatID int64, text string) {
	msg := tgbotapi.NewMessage(chatID, text)
	_, err := bot.Send(msg)
	if err != nil {
		log.Printf("Не удалось отправить сообщение в Telegram: %v", err)
	}
}
