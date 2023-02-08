package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"io/ioutil"
	"log"
	"strings"
)

type Lesson struct {
	Number int    `json:"number"`
	Name   string `json:"name"`
	Room   string `json:"room"`
}

type Lessons []Lesson

type Schedule struct {
	Class   string  `json:"class"`
	Lessons Lessons `json:"lessons"`
}

type Schedules []Schedule

type Day struct {
	Date     string    `json:"date"`
	Schedule Schedules `json:"schedules"`
}

func NewDay(d string, s Schedules) Day {

	day_val := Day{}

	day_val.Date = d
	day_val.Schedule = s

	return day_val
}

type Days []Day

// Сканер массива расписаний
func (d *Days) Scan(src interface{}) (err error) {

	// Массив байтов
	data := []byte{}

	// Приведение к байтам и запись в массив
	if val, ok := src.([]byte); ok {
		data = val
	} else if val, ok := src.([]byte); ok {
		data = []byte(val)
	} else if src == nil {
		return errors.New("couldn't convert db data to []byte")
	}

	// Новый reader для массива
	reader := bytes.NewReader(data)

	// Считываем байты
	bdata, err := ioutil.ReadAll(reader)
	if err != nil {
		log.Fatal(err)
	}

	// Удаляем экранирование и лишние скобки в начале и конце
	str := strings.ReplaceAll(string(bdata)[1:len(string(bdata))-1], "\\", "")
	// Удаляем лишние кавычки
	str = strings.ReplaceAll(str, "\"{", "{")
	str = strings.ReplaceAll(str, "}\"", "}")
	str = "[" + str + "]"

	// Объявляем карту с строчным ключём и значением в виде интерфейса
	var day_maps []map[string]interface{}

	// Записываем значения в карту
	err = json.Unmarshal([]byte(str), &day_maps)
	if err != nil {
		log.Fatal(err)
	}

	day_array := Days{}

	// Итерируемся по карте
	for _, day_map := range day_maps {

		day_date := ""

		// Проверяем наличие нужного ключа
		if day_map["Date"] != nil {
			if val, ok := day_map["Date"].(string); ok {
				day_date = val
			}
		}

		sch_array := Schedules{}

		// Проверяем наличие нужного ключа
		if day_map["Schedule"] != nil {

			// Приводим нужное значение к типу "массив интерфейсов"
			var lessons []interface{}

			if val, ok := day_map["Schedule"].([]interface{}); ok {
				lessons = val
			} else {
				return errors.New("couldn't convert Schedules to []interface{}")
			}

			// Итерируемся по массиву
			for _, lesson := range lessons {

				// Структура расписания
				sch := Schedule{}

				// Приводим интерфейсы к типу: "карта со строчным ключём и значением в виде интерфейса"
				var lesson_map map[string]interface{}

				if val, ok := lesson.(map[string]interface{}); ok {
					lesson_map = val
				} else {
					return errors.New("couldn't convert Schedule to map[string]interface{}")
				}

				// Проверяем наличие нужного ключа
				if lesson_map["Class"] != nil {

					// Записываем класс
					if val, ok := lesson_map["Class"].(string); ok {
						sch.Class = val
					} else {
						return errors.New("couldn't convert Class' name to string")
					}

				} else {
					return errors.New("Schedule has no Class key")
				}

				// Массив уроков
				les_array := Lessons{}

				// Проверяем наличие нужного ключа
				if lesson_map["Lessons"] != nil {

					// Приводим интерфейс к типу: "массив интерфейсов"
					var lesson_array []interface{}

					if val, ok := lesson_map["Lessons"].([]interface{}); ok {
						lesson_array = val
					} else {
						return errors.New("couldn't convert Lessons to []interface{}")
					}

					// Если массив не пустой
					if len(lesson_array) != 0 {
						for _, lesson_element := range lesson_array {

							// Структура урока
							les := Lesson{}

							// Приводим интерфейсы к типу: "карта со строчным ключём и значением в виде интерфейса"
							var lesson_element_map map[string]interface{}
							if val, ok := lesson_element.(map[string]interface{}); ok {
								lesson_element_map = val
							} else {
								return errors.New("couldn't convert Lesson to map[string]interface{}")
							}

							// Проверяем наличие нужного ключа
							if lesson_element_map["Name"] != nil {

								if val, ok := lesson_element_map["Name"].(string); ok {
									les.Name = val
								} else {
									return errors.New("couldn't convert Lesson's name to string")
								}

							} else {
								return errors.New("Lesson has no Name field")
							}

							// Проверяем наличие нужного ключа
							if lesson_element_map["Number"] != nil {

								// Приводим к int
								if val, ok := lesson_element_map["Number"].(float64); ok {
									les.Number = int(val)
								} else {
									return errors.New("couldn't convert Lesson's number to int")
								}

							} else {
								return errors.New("Lesson has no Number field")
							}

							// Проверяем наличие нужного ключа
							if lesson_element_map["Room"] != nil {

								if val, ok := lesson_element_map["Room"].(string); ok {
									les.Room = val
								} else {
									return errors.New("couldn't convert Room's name to string")
								}
							} else {
								return errors.New("Room has no Name field")
							}

							// Записываем урок в массив уроков
							les_array = append(les_array, les)

						}
					}

				} else {
					return errors.New("Schedule has no Lessons key")
				}

				// Записываем уроки
				sch.Lessons = les_array

				// Записываем расписания
				sch_array = append(sch_array, sch)

			}
		} else {
			return errors.New("Day has no Schedule key")
		}

		// Записываем дни
		day_array = append(day_array, NewDay(day_date, sch_array))
	}

	*d = day_array

	return nil
}

type Subject struct {
	Id          int     `json:"id"`
	Name        string  `json:"name"`
	Description *string `json:"description"`
}

type Room struct {
	Id    int    `json:"id"`
	Name  string `json:"name"`
	Wing  *int   `json:"wing"`
	Floor *int   `json:"floor"`
}

type Class struct {
	Id        int    `json:"id"`
	Number    int    `json:"number"`
	Character string `json:"сharacter"`
}
