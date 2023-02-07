package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strings"

	tgbotapi "github.com/Syfaro/telegram-bot-api"
)

type Weather struct {
	TheWeather Fact `json:"fact"`
}
type Fact struct {
	Temp int `json:"temp"`
}

const telegram_token = "5921290713:AAGtK4y74p4BqQCULhxBY3FZsHDCz9VMSvg"

const msc_url string = "https://api.weather.yandex.ru/v2/informers?lat=55.4424&lon=37.3636"
const ufa_url string = "https://api.weather.yandex.ru/v2/informers?lat=54.7431&lon=55.9678"

func (w Weather) GetWeather(url string) int { // метод, отправляющий запрос в янд погоду и получающий ответ
	client := &http.Client{}
	req, _ := http.NewRequest("GET", url, nil)
	req.Header.Set("X-Yandex-API-Key", "d5bbc17e-c642-4176-9df0-9bc77a9cfc0b")
	resp, _ := client.Do(req)
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println(err)
	}
	if err := json.Unmarshal(body, &w); err != nil { // полученный ответ формата json помещаю в структуру Weather
		fmt.Printf("Ошибка декод-я json в структуру [%s]", err.Error())
	}

	return w.TheWeather.Temp // возвращается только температура
}

var numericKeyboard = tgbotapi.NewReplyKeyboard( // кнопки, которые видит пользователь
	tgbotapi.NewKeyboardButtonRow(
		tgbotapi.NewKeyboardButton("Поприветствовать"),
		tgbotapi.NewKeyboardButton("Погода"),
	),
)
var numericKeyboard2 = tgbotapi.NewReplyKeyboard(
	tgbotapi.NewKeyboardButtonRow(
		tgbotapi.NewKeyboardButton("Уфа"),
		tgbotapi.NewKeyboardButton("Москва"),
		tgbotapi.NewKeyboardButton("Назад"),
	),
)

func main() {
	w := Weather{} // Структура, куда парсится полученный запрос формата Json из янд погоды
	// подключаемся к боту с помощью токена
	bot, err := tgbotapi.NewBotAPI(telegram_token)
	if err != nil {
		log.Panic(err)
	}

	bot.Debug = true
	log.Printf("Authorized on account %s", bot.Self.UserName)

	// инициализируем канал, куда будут прилетать обновления от API
	var ucfg tgbotapi.UpdateConfig = tgbotapi.NewUpdate(0)
	ucfg.Timeout = 60
	upd, err := bot.GetUpdatesChan(ucfg)
	if err != nil {
		fmt.Printf("Ошибка созд-я канала [%s]", err.Error())
	}
	// читаем обновления из канала
	for {
		select {
		case update := <-upd:
			//Пользователь, который написал боту
			UserName := update.Message.From.UserName

			// ID чата/диалога.
			// Может быть идентификатором как чата с пользователем
			// (тогда он равен UserID) так и публичного чата/канала
			ChatID := update.Message.Chat.ID

			// Текст сообщения
			Text := update.Message.Text
			log.Printf("[%s] %d %s", UserName, ChatID, Text)
			if update.Message == nil { // ignore non-Message updates
				continue
			}

			var msg = AnswerToUser(Text, w, UserName, ChatID) // метод, который формирует ответ пользователю, в соотв-ии
			// с нажатой кнопкой, результат метода помещается в переменную msg
			if _, err := bot.Send(msg); err != nil { // затем ответ отправляется пользователю
				log.Panic(err)
			}
		}
	}
}

func AnswerToUser(Text string, w Weather, UserName string, ChatID int64) tgbotapi.MessageConfig {
	var result_temp int
	var res string
	var msg = tgbotapi.MessageConfig{}

	res = strings.ToUpper(Text) //преобразую текст, полученный от пользователя, в верхний регистр, для того, чтобы
	// было неважно в каком регистре пользователь напишет слово open или close
	switch res {
	case "OPEN":
		msg = tgbotapi.NewMessage(ChatID, "Весь мой функционал на кнопках")
		msg.ReplyMarkup = numericKeyboard // открывает первые кнопки

	case "CLOSE":
		msg.ReplyMarkup = tgbotapi.NewRemoveKeyboard(true)

	case "УФА":
		result_temp = w.GetWeather(ufa_url)
		reply := fmt.Sprintf("Cейчас температура в Уфе  %d", result_temp)
		msg = tgbotapi.NewMessage(ChatID, reply) // отправляет температуру пользователю

	case "МОСКВА":
		result_temp = w.GetWeather(msc_url)
		reply := fmt.Sprintf("Cейчас температура в Москве  %d", result_temp)
		msg = tgbotapi.NewMessage(ChatID, reply) // отправляет также температуру

	case "ПОПРИВЕТСТВОВАТЬ":
		reply := fmt.Sprintf("Привет %s, я могу подсказать погоду)", UserName)
		msg = tgbotapi.NewMessage(ChatID, reply) // приветствует пользователя

	case "ПОГОДА":
		msg = tgbotapi.NewMessage(ChatID, "Вы выбрали раздел погода")
		msg.ReplyMarkup = numericKeyboard2 // открывает кнопки с городами

	case "НАЗАД":
		msg = tgbotapi.NewMessage(ChatID, "Вы вернулись назад")
		msg.ReplyMarkup = numericKeyboard // Добавил кнопку назад и default

	default:
		msg = tgbotapi.NewMessage(ChatID, "Неизвестная команда")
	}

	return msg
}
