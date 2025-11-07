package telegram

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"
)

// 1. 💡 НОВАЯ СТРУКТУРА (Для "Ответа на Кнопку")
// Она нужна, чтобы "погасить часики"
type AnswerCallbackQueryRequest struct {
	CallbackQueryID string `json:"callback_query_id"`
	// (Мы можем добавить сюда Text, но если оставить пустым,
	// Telegram просто "погасит часики")
}

// одна кнопка
type InlineKeyboardButton struct {
	Text         string `json:"text"`          // Текст, который видит пользователь (например, "УСН 6%")
	CallbackData string `json:"callback_data"` // Секретная команда, которую бот получит при нажатии (например, "tax_6")
}

// ряд кнопок
type InlineKeyboardMarkup struct {
	InlineKeyboard [][]InlineKeyboardButton `json:"inline_keyboard"`
}

// Это структура для "понимания" *ответа* от /getUpdates
type Update struct {
	ID      int      `json:"update_id"` // 👈 ID самого обновления
	Message *Message `json:"message"`   // 👈 Сообщение (команда /start, /price)
	// Запрос от кнопки. Делаем с указателем т.к. может быть nil.
	CallbackQuery *CallbackQuery `json:"callback_query,omitempty"`
}

type CallbackQuery struct {
	ID      string   `json:"id"`
	Data    string   `json:"data"`    // Вот здесь лежит "tax_6"
	Message *Message `json:"message"` // Сообщение, под которым была кнопка
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

// Это структура для "отправки" (запрос на отправку сообщения)
type SendMessageRequest struct {
	ChatID      int64                 `json:"chat_id"`
	Text        string                `json:"text"`
	ReplyMarkup *InlineKeyboardMarkup `json:"reply_markup,omitempty"` // Принимает ссылку, чтобы не отправляться если nil. Надо понять!
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
func (c *Client) SendMessage(chatID int64, text string, markup *InlineKeyboardMarkup) error {
	url := c.botURL + "/sendMessage"

	// 1. "Мысль" (Struct)
	msg := SendMessageRequest{
		ChatID:      chatID,
		Text:        text,
		ReplyMarkup: markup, // Уточнить что делает!
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

// Этот метод вызывает SendMessage, чтобы нам было удобнее, но удобнее что?
func (c *Client) SendMessageWithButtons(chatID int64, text string, markup *InlineKeyboardMarkup) error { // Почему тип данных int64
	return c.SendMessage(chatID, text, markup) //Метод который просто вызывает метод зачем?
}

// 2. 💡 НОВЫЙ МЕТОД (Ответить на нажатие кнопки)
func (c *Client) AnswerCallbackQuery(queryID string) error {
	url := c.botURL + "/answerCallbackQuery"

	// 1. "Мысль" (Struct)
	req := AnswerCallbackQueryRequest{
		CallbackQueryID: queryID,
	}

	// 2. "Говорим" (Marshal)
	jsonData, err := json.Marshal(req)
	if err != nil {
		return err
	}

	// 3. "Звоним" (POST)
	// (Нам не важен ответ, нам важен сам факт "звонка")
	resp, err := c.httpClient.Post(url, "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("статус не 200 OK: %s", resp.Status)
	}

	return nil
}
