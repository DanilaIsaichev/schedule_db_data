package schedule_db_data

import (
	"encoding/json"
	"errors"
	"strings"
)

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

func (teachers *Teachers) Contain(teacher Teacher) (res bool) {

	for _, t := range *teachers {
		if t.Login == teacher.Login {
			return true
		}
	}

	return false
}

func (teachers *Teachers) Find(login string) (teacher Teacher, err error) {

	for _, teacher := range *teachers {
		if teacher.Login == login {
			return teacher, nil
		}
	}

	return Teacher{}, errors.New("no teacher with login " + login + " has found")
}
