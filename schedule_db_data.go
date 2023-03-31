package schedule_db_data

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"errors"
	"io"
	"strings"
)

type Subject struct {
	Id          int     `json:"id"`
	Name        string  `json:"name"`
	Description *string `json:"description"`
}

type Subjects []Subject

// Сканер массива учителей
func (s *Subjects) Scan(src interface{}) (err error) {

	// Приведение полученных данных к корректному виду (массив байтов без служебных символов)
	byte_str, err := scan_prepare(src)
	if err != nil {
		return err
	}

	// Считывание данных из массива байтов
	err = s.scan_subjects(byte_str)
	if err != nil {
		return err
	}

	return nil
}

func (s *Subjects) scan_subjects(src []byte) (err error) {

	// Удаляем экранирование и лишние скобки в начале и конце
	str := strings.ReplaceAll(string(src)[1:len(string(src))-1], "\\", "")
	// Удаляем лишние кавычки
	str = strings.ReplaceAll(str, "\"{", "{")
	str = strings.ReplaceAll(str, "}\"", "}")
	str = "[" + str + "]"

	// Объявляем карту с строчным ключём и значением в виде интерфейса
	var subject_maps []map[string]interface{}

	// Записываем значения в карту
	err = json.Unmarshal([]byte(str), &subject_maps)
	if err != nil {
		return err
	}

	subjects_array := Subjects{}

	// Итерируемся по карте
	for _, subject_map := range subject_maps {

		subject := Subject{}

		// Проверяем наличие нужного ключа
		if subject_map["id"] != nil {
			if val, ok := subject_map["id"].(float64); ok {
				subject.Id = int(val)
			} else {
				return errors.New("couldn't convert 'id' to int")
			}
		} else {
			return errors.New("subject has no 'id'")
		}

		if subject_map["name"] != nil {
			if val, ok := subject_map["name"].(string); ok {
				subject.Name = val
			} else {
				return errors.New("couldn't convert 'name' to string")
			}
		} else {
			return errors.New("subject has no 'name'")
		}

		if subject_map["description"] != nil {
			if val, ok := subject_map["description"].(string); ok {
				subject.Description = new(string)
				*subject.Description = val
			} else {
				return errors.New("couldn't convert 'login' to string")
			}
		}

		subjects_array = append(subjects_array, subject)
	}

	*s = subjects_array

	return nil
}

func UnmarshalSubject(src interface{}) (Subject, string, error) {

	// Приведение полученных данных к корректному виду (массив байтов без служебных символов)
	byte_str, err := scan_prepare(src)
	if err != nil {
		return Subject{}, "", err
	}

	// Объявляем карту с строчным ключём и значением в виде интерфейса
	var subject_map map[string]interface{}

	// Записываем значения в карту
	err = json.Unmarshal(byte_str, &subject_map)
	if err != nil {
		return Subject{}, "", err
	}

	s := Subject{}
	a := ""

	if subject_map["action"] != nil {
		if val, ok := subject_map["action"].(string); ok {
			switch {
			case val == "save":
				a = "save"
			case val == "delete":
				a = "delete"
			default:
				return Subject{}, "", errors.New("unknown action")
			}
		} else {
			return Subject{}, "", errors.New("couldn't convert action to string")
		}
	} else {
		return Subject{}, "", errors.New("no action found")
	}

	if subject_map["subject"] != nil {
		if val, ok := subject_map["subject"].(map[string]interface{}); ok {

			subject := val

			if subject["name"] != nil {
				if val, ok := subject["name"].(string); ok {
					s.Name = val
				} else {
					return Subject{}, "", errors.New("couldn't convert subject's name to string")
				}
			} else {
				return Subject{}, "", errors.New("no subject's name found")
			}

			if subject["description"] != nil {
				if val, ok := subject["description"].(string); ok {
					if len(val) > 0 {
						s.Description = new(string)
						*(s.Description) = val
					}
				} else {
					return Subject{}, "", errors.New("couldn't convert subject's description to string")
				}
			}

		} else {
			return Subject{}, "", errors.New("couldn't convert subject's struct to map[string]interface{}")
		}
	} else {
		return Subject{}, "", errors.New("no 'subject' found")
	}

	return s, a, nil
}

type Room struct {
	Id    int    `json:"id"`
	Name  string `json:"name"`
	Wing  *int   `json:"wing"`
	Floor *int   `json:"floor"`
}

type Rooms []Room

// Сканер массива учителей
func (r *Rooms) Scan(src interface{}) (err error) {

	// Приведение полученных данных к корректному виду (массив байтов без служебных символов)
	byte_str, err := scan_prepare(src)
	if err != nil {
		return err
	}

	// Считывание данных из массива байтов
	err = r.scan_rooms(byte_str)
	if err != nil {
		return err
	}

	return nil
}

func (r *Rooms) scan_rooms(src []byte) (err error) {

	// Удаляем экранирование и лишние скобки в начале и конце
	str := strings.ReplaceAll(string(src)[1:len(string(src))-1], "\\", "")
	// Удаляем лишние кавычки
	str = strings.ReplaceAll(str, "\"{", "{")
	str = strings.ReplaceAll(str, "}\"", "}")
	str = "[" + str + "]"

	// Объявляем карту с строчным ключём и значением в виде интерфейса
	var room_maps []map[string]interface{}

	// Записываем значения в карту
	err = json.Unmarshal([]byte(str), &room_maps)
	if err != nil {
		return err
	}

	rooms_array := Rooms{}

	// Итерируемся по карте
	for _, room_map := range room_maps {

		room := Room{}

		// Проверяем наличие нужного ключа
		if room_map["id"] != nil {
			if val, ok := room_map["id"].(float64); ok {
				room.Id = int(val)
			} else {
				return errors.New("couldn't convert 'id' to int")
			}
		} else {
			return errors.New("room has no 'id'")
		}

		if room_map["name"] != nil {
			if val, ok := room_map["name"].(string); ok {
				room.Name = val
			} else {
				return errors.New("couldn't convert 'name' to string")
			}
		} else {
			return errors.New("room has no 'name'")
		}

		if room_map["wing"] != nil {
			if val, ok := room_map["wing"].(float64); ok {
				room.Wing = new(int)
				*room.Wing = int(val)
			} else {
				return errors.New("couldn't convert 'wing' to int")
			}
		}

		if room_map["floor"] != nil {
			if val, ok := room_map["floor"].(float64); ok {
				room.Floor = new(int)
				*room.Floor = int(val)
			} else {
				return errors.New("couldn't convert 'floor' to int")
			}
		}

		rooms_array = append(rooms_array, room)
	}

	*r = rooms_array

	return nil
}

func UnmarshalRoom(src interface{}) (Room, string, error) {

	// Приведение полученных данных к корректному виду (массив байтов без служебных символов)
	byte_str, err := scan_prepare(src)
	if err != nil {
		return Room{}, "", err
	}

	// Объявляем карту с строчным ключём и значением в виде интерфейса
	var room_map map[string]interface{}

	// Записываем значения в карту
	err = json.Unmarshal(byte_str, &room_map)
	if err != nil {
		return Room{}, "", err
	}

	r := Room{}
	a := ""

	if room_map["action"] != nil {
		if val, ok := room_map["action"].(string); ok {
			switch {
			case val == "save":
				a = "save"
			case val == "delete":
				a = "delete"
			default:
				return Room{}, "", errors.New("unknown action")
			}
		} else {
			return Room{}, "", errors.New("couldn't convert action to string")
		}
	} else {
		return Room{}, "", errors.New("no action found")
	}

	if room_map["room"] != nil {
		if val, ok := room_map["room"].(map[string]interface{}); ok {

			room := val

			if room["name"] != nil {
				if val, ok := room["name"].(string); ok {
					r.Name = val
				} else {
					return Room{}, "", errors.New("couldn't convert rooms's name to string")
				}
			} else {
				return Room{}, "", errors.New("no rooms's name found")
			}

			if room["wing"] != nil {
				if val, ok := room["wing"].(float64); ok {
					r.Wing = new(int)
					*(r.Wing) = int(val)
				} else {
					return Room{}, "", errors.New("couldn't convert room's wing to int")
				}
			}

			if room["floor"] != nil {
				if val, ok := room["floor"].(float64); ok {
					r.Floor = new(int)
					*(r.Floor) = int(val)
				} else {
					return Room{}, "", errors.New("couldn't convert room's floor to string")
				}
			}
		} else {
			return Room{}, "", errors.New("couldn't convert room's struct to map[string]interface{}")
		}
	} else {
		return Room{}, "", errors.New("no 'room' found")
	}

	return r, a, nil
}

type Class struct {
	Id        int    `json:"id"`
	Number    int    `json:"number"`
	Character string `json:"сharacter"`
}

type Classes []Class

// Сканер массива учителей
func (c *Classes) Scan(src interface{}) (err error) {

	// Приведение полученных данных к корректному виду (массив байтов без служебных символов)
	byte_str, err := scan_prepare(src)
	if err != nil {
		return err
	}

	// Считывание данных из массива байтов
	err = c.scan_classes(byte_str)
	if err != nil {
		return err
	}

	return nil
}

func (c *Classes) scan_classes(src []byte) (err error) {

	// Удаляем экранирование и лишние скобки в начале и конце
	str := strings.ReplaceAll(string(src)[1:len(string(src))-1], "\\", "")
	// Удаляем лишние кавычки
	str = strings.ReplaceAll(str, "\"{", "{")
	str = strings.ReplaceAll(str, "}\"", "}")
	str = "[" + str + "]"

	// Объявляем карту с строчным ключём и значением в виде интерфейса
	var class_maps []map[string]interface{}

	// Записываем значения в карту
	err = json.Unmarshal([]byte(str), &class_maps)
	if err != nil {
		return err
	}

	classes_array := Classes{}

	// Итерируемся по карте
	for _, class_map := range class_maps {

		class := Class{}

		// Проверяем наличие нужного ключа
		if class_map["id"] != nil {
			if val, ok := class_map["id"].(float64); ok {
				class.Id = int(val)
			} else {
				return errors.New("couldn't convert 'id' to int")
			}
		} else {
			return errors.New("class has no 'id'")
		}

		if class_map["number"] != nil {
			if val, ok := class_map["number"].(float64); ok {
				if int(val) >= 1 && int(val) <= 11 {
					class.Number = int(val)
				} else {
					return errors.New("class number not in [1:11]")
				}
			} else {
				return errors.New("couldn't convert 'number' to int")
			}
		} else {
			return errors.New("class has no 'number'")
		}

		if class_map["character"] != nil {
			if val, ok := class_map["character"].(string); ok {
				if len([]rune(val)) == 1 {
					class.Character = val
				} else {
					return errors.New("wrong class character")
				}
			} else {
				return errors.New("couldn't convert 'character' to string")
			}
		} else {
			return errors.New("class has no 'character'")
		}

		classes_array = append(classes_array, class)
	}

	*c = classes_array

	return nil
}

func UnmarshalClass(src interface{}) (Class, string, error) {

	// Приведение полученных данных к корректному виду (массив байтов без служебных символов)
	byte_str, err := scan_prepare(src)
	if err != nil {
		return Class{}, "", err
	}

	// Объявляем карту с строчным ключём и значением в виде интерфейса
	var class_map map[string]interface{}

	// Записываем значения в карту
	err = json.Unmarshal(byte_str, &class_map)
	if err != nil {
		return Class{}, "", err
	}

	c := Class{}
	a := ""

	if class_map["action"] != nil {
		if val, ok := class_map["action"].(string); ok {
			switch {
			case val == "save":
				a = "save"
			case val == "delete":
				a = "delete"
			default:
				return Class{}, "", errors.New("unknown action")
			}
		} else {
			return Class{}, "", errors.New("couldn't convert action to string")
		}
	} else {
		return Class{}, "", errors.New("no action found")
	}

	if class_map["class"] != nil {
		if val, ok := class_map["class"].(map[string]interface{}); ok {

			class := val

			if class["number"] != nil {
				if val, ok := class["number"].(float64); ok {
					if int(val) >= 1 && int(val) <= 11 {
						c.Number = int(val)
					} else {
						return Class{}, "", errors.New("wrong class number")
					}
				} else {
					return Class{}, "", errors.New("couldn't convert class' number to int")
				}
			} else {
				return Class{}, "", errors.New("no class' number found")
			}

			if class["character"] != nil {
				if val, ok := class["character"].(string); ok {
					if len([]rune(val)) == 1 {
						c.Character = val
					} else {
						return Class{}, "", errors.New("wrong class character")
					}
				} else {
					return Class{}, "", errors.New("couldn't convert class character to string")
				}
			} else {
				return Class{}, "", errors.New("no class character found")
			}

		} else {
			return Class{}, "", errors.New("couldn't convert class' struct to map[string]interface{}")
		}
	} else {
		return Class{}, "", errors.New("no 'class' found")
	}

	return c, a, nil
}

type Lesson struct {
	Number  int     `json:"number"`
	Name    string  `json:"name"`
	Room    string  `json:"room"`
	Teacher Teacher `json:"teacher"`
}

type Teacher struct {
	Id    int    `json:"id"`
	Name  string `json:"name"`
	Login string `json:"login"`
}

type Teachers []Teacher

// Сканер массива учителей
func (t *Teachers) Scan(src interface{}) (err error) {

	// Приведение полученных данных к корректному виду (массив байтов без служебных символов)
	byte_str, err := scan_prepare(src)
	if err != nil {
		return err
	}

	// Считывание данных из массива байтов
	err = t.scan_teachers(byte_str)
	if err != nil {
		return err
	}

	return nil
}

func (t *Teachers) scan_teachers(src []byte) (err error) {

	// Удаляем экранирование и лишние скобки в начале и конце
	str := strings.ReplaceAll(string(src)[1:len(string(src))-1], "\\", "")
	// Удаляем лишние кавычки
	str = strings.ReplaceAll(str, "\"{", "{")
	str = strings.ReplaceAll(str, "}\"", "}")
	str = "[" + str + "]"

	// Объявляем карту с строчным ключём и значением в виде интерфейса
	var teacher_maps []map[string]interface{}

	// Записываем значения в карту
	err = json.Unmarshal([]byte(str), &teacher_maps)
	if err != nil {
		return err
	}

	teachers_array := Teachers{}

	// Итерируемся по карте
	for _, teacher_map := range teacher_maps {

		teacher := Teacher{}

		// Проверяем наличие нужного ключа
		if teacher_map["id"] != nil {
			if val, ok := teacher_map["id"].(float64); ok {
				teacher.Id = int(val)
			} else {
				return errors.New("couldn't convert 'id' to int")
			}
		} else {
			return errors.New("teacher has no 'id'")
		}

		if teacher_map["name"] != nil {
			if val, ok := teacher_map["name"].(string); ok {
				teacher.Name = val
			} else {
				return errors.New("couldn't convert 'name' to string")
			}
		} else {
			return errors.New("teacher has no 'name'")
		}

		if teacher_map["login"] != nil {
			if val, ok := teacher_map["login"].(string); ok {
				teacher.Login = val
			} else {
				return errors.New("couldn't convert 'login' to string")
			}
		} else {
			return errors.New("teacher has no 'login'")
		}

		teachers_array = append(teachers_array, teacher)
	}

	*t = teachers_array

	return nil
}

func UnmarshalTeacher(src interface{}) (Teacher, string, error) {

	// Приведение полученных данных к корректному виду (массив байтов без служебных символов)
	byte_str, err := scan_prepare(src)
	if err != nil {
		return Teacher{}, "", err
	}

	// Объявляем карту с строчным ключём и значением в виде интерфейса
	var teacher_map map[string]interface{}

	// Записываем значения в карту
	err = json.Unmarshal(byte_str, &teacher_map)
	if err != nil {
		return Teacher{}, "", err
	}

	t := Teacher{}
	a := ""

	if teacher_map["action"] != nil {
		if val, ok := teacher_map["action"].(string); ok {
			switch {
			case val == "save":
				a = "save"
			case val == "delete":
				a = "delete"
			default:
				return Teacher{}, "", errors.New("unknown action")
			}
		} else {
			return Teacher{}, "", errors.New("couldn't convert action to string")
		}
	} else {
		return Teacher{}, "", errors.New("no action found")
	}

	if teacher_map["teacher"] != nil {
		if val, ok := teacher_map["teacher"].(map[string]interface{}); ok {

			teacher := val

			if teacher["name"] != nil {
				if val, ok := teacher["name"].(string); ok {
					t.Name = val
				} else {
					return Teacher{}, "", errors.New("couldn't convert teacher's name to string")
				}
			} else {
				return Teacher{}, "", errors.New("no teacher's name found")
			}

			if teacher["login"] != nil {
				if val, ok := teacher["login"].(string); ok {
					t.Login = val
				} else {
					return Teacher{}, "", errors.New("couldn't convert teacher's login to string")
				}
			} else {
				return Teacher{}, "", errors.New("no teacher's login found")
			}

		} else {
			return Teacher{}, "", errors.New("couldn't convertteacher's struct to map[string]interface{}")
		}
	} else {
		return Teacher{}, "", errors.New("no 'teacher' found")
	}

	return t, a, nil
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
				return errors.New("couldn't convert 'date' to string")
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

								if val, ok := lesson_element_map["teacher"].(map[string]interface{}); ok {

									teacher := val

									if val, ok := teacher["name"].(string); ok {
										les.Teacher.Name = val
									} else {
										return errors.New("couldn't convert name of 'teacher' to string")
									}

									if val, ok := teacher["login"].(string); ok {
										les.Teacher.Login = val
									} else {
										return errors.New("couldn't convert login of 'teacher' to string")
									}

								} else {
									return errors.New("couldn't convert teacher object to to map[string]interface{}")
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
	Start    string `json:"start"`
	Year     int    `json:"year"`
	Parallel int    `json:"parallel"`
	Is_Base  bool   `json:"is_base"`
	Data     Days   `json:"data"`
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

	w := Week{}

	if week_map["start"] != nil {
		if val, ok := week_map["start"].(string); ok {
			w.Start = val
		} else {
			return Week{}, errors.New("couldn't convert week start date to string")
		}
	}

	if week_map["year"] != nil {
		if val, ok := week_map["year"].(float64); ok {
			w.Year = int(val)
		} else {
			return Week{}, errors.New("couldn't convert year to int")
		}
	} else {
		return Week{}, errors.New("no 'year' found")
	}

	if week_map["parallel"] != nil {
		if val, ok := week_map["parallel"].(float64); ok {
			w.Parallel = int(val)
		} else {
			return Week{}, errors.New("couldn't convert parallel to int")
		}
	} else {
		return Week{}, errors.New("no 'parallel' found")
	}

	if week_map["is_base"] != nil {
		if val, ok := week_map["is_base"].(bool); ok {
			w.Is_Base = val
		} else {
			return Week{}, errors.New("couldn't convert is_base to bool")
		}
	} else {
		return Week{}, errors.New("no 'is_base' found")
	}

	days := Days{}

	if week_map["data"] != nil {

		if val, ok := week_map["data"].([]byte); ok {

			err := days.scan_days(val)
			if err != nil {
				return Week{}, err
			}

		} else {
			return Week{}, errors.New("couldn't convert week start date to string")
		}
	} else {
		return Week{}, errors.New("no 'start' found")
	}

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
	bdata, err := io.ReadAll(reader)
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
