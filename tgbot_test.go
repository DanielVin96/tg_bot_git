package main

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestAnswerToUserOpen(t *testing.T) { // тестирование кнопок
	name := "DanielVin"

	// Временный "Костыль" использования Map, пока не нашел способ, как проводить несколько тестов одновременно
	textList := map[string]string{ // map, где ключ - это текст пользователя, а значение - это возвращаемый текст
		"Open":             "Весь мой функционал на кнопках",
		"Поприветствовать": fmt.Sprintf("Привет %s, я могу подсказать погоду)", name),
		"Погода":           "Вы выбрали раздел погода",
		"Назад":            "Вы вернулись назад",
		"default":          "Неизвестная команда",
	}
	for key, value := range textList { // иду циклом по map и вызываю тест на проверку каждой кнопки
		assert.Equal(t, value, AnswerToUser(key, Weather{}, name, 25).Text) //В агрумент метода Equal в качестве value попадает
		// ожидаемый текст, а в аргумент метода AnswerToUser в качестве key попадает текст пользователя
	}

}
func TestAnswerToUserUfa(t *testing.T) { // тестирование кнопки, которая возвращает текст с погодой в Уфе
	w := Weather{}
	res_weather := w.GetWeather(ufa_url)
	expectedUfa := fmt.Sprintf("Cейчас температура в Уфе  %d", res_weather) // ожидаемый текст
	assert.Equal(t, expectedUfa, AnswerToUser("Уфа", Weather{}, "DanielVin", 25).Text)
}
func TestAnswerToUserMsc(t *testing.T) { // тестирование кнопки, которая возвращает текст с погодой в Москве
	w := Weather{}
	res_weather := w.GetWeather(msc_url)
	expectedMsc := fmt.Sprintf("Cейчас температура в Москве  %d", res_weather) // ожидаемый текст
	assert.Equal(t, expectedMsc, AnswerToUser("Москва", Weather{}, "DanielVin", 25).Text)
}
