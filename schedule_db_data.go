package schedule_db_data

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"errors"
	"io/ioutil"
	"strings"
)

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

type Lesson struct {
	Number  int    `json:"number"`
	Name    string `json:"name"`
	Room    string `json:"room"`
	Teacher string `json:"teacher"`
}

type Lessons []Lesson

type Schedule struct {
	Class   string  `json:"class"`
	Lessons Lessons `json:"lessons"`
}

type Schedules []Schedule

type Day struct {
	Date     string    `json:"date"`
	Schedule Schedules `json:"schedule"`
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

	// Приведение полученных данных к корректному виду (массив байтов без служебных символов)
	byte_str, err := scan_prepare(src)
	if err != nil {
		return err
	}

	// Считывание данных из массива байтов
	err = d.scan_days(byte_str)
	if err != nil {
		return err
	}

	return nil
}

func (d *Days) scan_days(src []byte) (err error) {

	// Удаляем экранирование и лишние скобки в начале и конце
	str := strings.ReplaceAll(string(src)[1:len(string(src))-1], "\\", "")
	// Удаляем лишние кавычки
	str = strings.ReplaceAll(str, "\"{", "{")
	str = strings.ReplaceAll(str, "}\"", "}")
	str = "[" + str + "]"

	// Объявляем карту с строчным ключём и значением в виде интерфейса
	var day_maps []map[string]interface{}

	// Записываем значения в карту
	err = json.Unmarshal([]byte(str), &day_maps)
	if err != nil {
		return err
	}

	day_array := Days{}

	// Итерируемся по карте
	for _, day_map := range day_maps {

		day_date := ""

		// Проверяем наличие нужного ключа
		if day_map["date"] != nil {
			if val, ok := day_map["date"].(string); ok {
				day_date = val
			} else {
				return errors.New("couldn't covert 'date' to string")
			}
		} else {
			return errors.New("day has no 'date'")
		}

		sch_array := Schedules{}

		// Проверяем наличие нужного ключа
		if day_map["schedule"] != nil {

			// Приводим нужное значение к типу "массив интерфейсов"
			var lessons []interface{}

			if val, ok := day_map["schedule"].([]interface{}); ok {
				lessons = val
			} else {
				return errors.New("couldn't convert array of 'schedule' to []interface{}")
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
					return errors.New("couldn't convert 'schedule' to map[string]interface{}")
				}

				// Проверяем наличие нужного ключа
				if lesson_map["class"] != nil {

					// Записываем класс
					if val, ok := lesson_map["class"].(string); ok {
						sch.Class = val
					} else {
						return errors.New("couldn't convert 'name' of 'class' to string")
					}

				} else {
					return errors.New("'schedule' has no 'class' key")
				}

				// Массив уроков
				les_array := Lessons{}

				// Проверяем наличие нужного ключа
				if lesson_map["lessons"] != nil {

					// Приводим интерфейс к типу: "массив интерфейсов"
					var lesson_array []interface{}

					if val, ok := lesson_map["lessons"].([]interface{}); ok {
						lesson_array = val
					} else {
						return errors.New("couldn't convert 'lessons' to []interface{}")
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
								return errors.New("couldn't convert 'lesson' to map[string]interface{}")
							}

							// Проверяем наличие нужного ключа
							if lesson_element_map["name"] != nil {

								if val, ok := lesson_element_map["name"].(string); ok {
									les.Name = val
								} else {
									return errors.New("couldn't convert 'name' of 'lesson' to string")
								}

							} else {
								return errors.New("'lesson' has no 'name' field")
							}

							// Проверяем наличие нужного ключа
							if lesson_element_map["number"] != nil {

								// Приводим к int
								if val, ok := lesson_element_map["number"].(float64); ok {
									les.Number = int(val)
								} else {
									return errors.New("couldn't convert 'number' of 'lesson' to int")
								}

							} else {
								return errors.New("'lesson' has no 'number' field")
							}

							// Проверяем наличие нужного ключа
							if lesson_element_map["room"] != nil {

								if val, ok := lesson_element_map["room"].(string); ok {
									les.Room = val
								} else {
									return errors.New("couldn't convert name of 'room' to string")
								}
							} else {
								return errors.New("'lesson' has no 'room' field")
							}

							// Проверяем наличие нужного ключа
							if lesson_element_map["teacher"] != nil {

								if val, ok := lesson_element_map["teacher"].(string); ok {
									les.Teacher = val
								} else {
									return errors.New("couldn't convert name of 'teacher' to string")
								}

							} else {
								return errors.New("'lesson' has no 'teacher' field")
							}

							// Записываем урок в массив уроков
							les_array = append(les_array, les)

						}
					}

				} else {
					return errors.New("'schedule' has no 'lessons' key")
				}

				// Записываем уроки
				sch.Lessons = les_array

				// Записываем расписания
				sch_array = append(sch_array, sch)

			}
		} else {
			return errors.New("'day' has no 'schedule' key")
		}

		// Записываем дни
		day_array = append(day_array, NewDay(day_date, sch_array))
	}

	*d = day_array

	return nil
}

type Week struct {
	Start string `json:"start"`
	Data  Days   `json:"data"`
}

func UnmarshalWeek(src interface{}) (week Week, err error) {

	// Приведение полученных данных к корректному виду (массив байтов без служебных символов)
	byte_str, err := scan_prepare(src)
	if err != nil {
		return Week{}, err
	}

	// Объявляем карту с строчным ключём и значением в виде интерфейса
	var week_map map[string]interface{}

	// Записываем значения в карту
	err = json.Unmarshal(byte_str, &week_map)
	if err != nil {
		return Week{}, err
	}

	start_date := ""

	if week_map["start"] != nil {
		if val, ok := week_map["start"].(string); ok {
			start_date = val
		} else {
			return Week{}, errors.New("couldn't conver week start date to string")
		}
	} else {
		return Week{}, errors.New("no 'start' found")
	}

	days := Days{}

	if week_map["data"] != nil {

		if val, ok := week_map["data"].([]byte); ok {

			err := days.scan_days(val)
			if err != nil {
				return Week{}, err
			}

		} else {
			return Week{}, errors.New("couldn't conver week start date to string")
		}
	} else {
		return Week{}, errors.New("no 'start' found")
	}

	w := Week{}

	w.Start = start_date
	w.Data = days

	return w, nil
}

func scan_prepare(src interface{}) (prepared_bytes []byte, err error) {

	// Массив байтов
	data := []byte{}

	// Приведение к байтам и запись в массив
	if val, ok := src.([]byte); ok {
		data = val
	} else if val, ok := src.([]byte); ok {
		data = []byte(val)
	} else if src == nil {
		return []byte{}, errors.New("couldn't convert db data to []byte")
	}

	// Новый reader для массива
	reader := bytes.NewReader(data)

	// Считываем байты
	bdata, err := ioutil.ReadAll(reader)
	if err != nil {
		return []byte{}, err
	}

	return bdata, nil
}

func DB_connection(hostname string, db_name string, username string, password string, port string) (db_conn *sql.DB, err error) {

	connection_string := "host=" + hostname + " port=" + port + " user=" + username + " password=" + password + " dbname=" + db_name + " sslmode=disable"

	db, err := sql.Open("postgres", connection_string)
	if err != nil {
		return db, err
	}

	return db, nil
}
