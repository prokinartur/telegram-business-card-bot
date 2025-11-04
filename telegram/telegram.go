package telegram

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"
)

// 💡 НЮАНС 1: Вложенные структуры
// Ответ от Telegram - это сложный JSON. Нам нужны
// структуры Go, чтобы "понять" (Unmarshal) его.
//
// Update (Обновление)
// └── Message (Сообщение)
//
//	├── Chat (Чат, откуда пришло)
//	│   └── ID (ID чата)
//	└── Text (Текст сообщения)
type Update struct {
	ID      int     `json:"update_id"` // 👈 ID самого обновления
	Message Message `json:"message"`
}

type Message struct {
	Text string `json:"text"`
	Chat Chat   `json:"chat"`
}

type Chat struct {
	ID int64 `json:"id"` // 👈 ID чата (пользователя)
}

// Это структура для "понимания" *ответа* от /getUpdates
type GetUpdatesResponse struct {
	Ok     bool     `json:"ok"`
	Result []Update `json:"result"` // 👈 Массив (срез) обновлений
}

// Это структура для "отправки" (мы её знаем)
type SendMessageRequest struct {
	ChatID int64  `json:"chat_id"`
	Text   string `json:"text"`
}

// Client - это наш "пульт"
type Client struct {
	botURL     string      // (https://api.telegram.org/bot<TOKEN>)
	httpClient http.Client // (Наш "телефон")
}

// NewClient - "конструктор" для нашего пульта
func NewClient(token string) *Client {
	return &Client{
		botURL:     "https://api.telegram.org/bot" + token,
		httpClient: http.Client{},
	}
}

// ---
// 💡 МЕТОД 1: "Слушать" (GetUpdates)
// ---
func (c *Client) GetUpdates(offset int) ([]Update, error) {
	// 💡 НЮАНС 2: Offset
	// Мы "звоним" и говорим: "Дай мне обновления, НАЧИНАЯ С (offset)"
	// Это нужно, чтобы не получать одни и те же сообщения 100 раз
	url := c.botURL + "/getUpdates?offset=" + strconv.Itoa(offset)

	// 1. "Звоним" (GET-запрос)
	resp, err := c.httpClient.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	// 2. Читаем ответ
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	// 3. "Понимаем" (Unmarshal)
	var response GetUpdatesResponse
	if err := json.Unmarshal(body, &response); err != nil {
		return nil, err
	}

	// Возвращаем только *список* обновлений
	return response.Result, nil
}

// ---
// 💡 МЕТОД 2: "Говорить" (SendMessage)
// (Это наш старый код, но "завернутый" в метод)
// ---
func (c *Client) SendMessage(chatID int64, text string) error {
	url := c.botURL + "/sendMessage"

	// 1. "Мысль" (Struct)
	msg := SendMessageRequest{
		ChatID: chatID,
		Text:   text,
	}

	// 2. "Говорим" (Marshal)
	jsonData, err := json.Marshal(msg)
	if err != nil {
		return err
	}

	// 3. "Звоним" (POST)
	resp, err := c.httpClient.Post(url, "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	// 4. Проверяем статус
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("статус не 200 OK: %s", resp.Status)
	}

	return nil
}
