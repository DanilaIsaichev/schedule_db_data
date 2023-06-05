package schedule_db_data

import (
	"encoding/json"
	"errors"
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
