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

func UnmarshalTeacher(src interface{}) (Teachers, error) {

	// Приведение полученных данных к корректному виду (массив байтов без служебных символов)
	byte_str, err := scan_prepare(src)
	if err != nil {
		return Teachers{}, err
	}

	// Объявляем массив карт со строчным ключём и значением в виде интерфейса
	var teacher_maps []map[string]interface{}

	// Записываем значения в массив карт
	err = json.Unmarshal(byte_str, &teacher_maps)
	if err != nil {
		return Teachers{}, err
	}

	t := Teachers{}

	for _, teacher_map := range teacher_maps {

		teacher := Teacher{}

		if teacher_map["login"] != nil {
			if val, ok := teacher_map["login"].(string); ok {
				teacher.Login = val
			} else {
				return Teachers{}, errors.New("couldn't convert teacher's login to string")
			}
		} else {
			return Teachers{}, errors.New("no teacher's login found")
		}

		if teacher_map["name"] != nil {
			if val, ok := teacher_map["name"].(string); ok {
				teacher.Name = val
			} else {
				return Teachers{}, errors.New("couldn't convert teacher's name to string")
			}
		} else {
			return Teachers{}, errors.New("no teacher's name found")
		}

		t = append(t, teacher)
	}

	return t, nil
}
