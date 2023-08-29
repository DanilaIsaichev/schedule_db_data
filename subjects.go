package schedule_db_data

import (
	"encoding/json"
	"errors"
	"strings"
)

type Subject struct {
	Name        string `json:"name"`
	Description string `json:"description"`
}

type Subjects []Subject

// Сканер массива предметов
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
				subject.Description = val
			} else {
				return errors.New("couldn't convert 'description' to string")
			}
		} else {
			subject.Description = ""
		}

		subjects_array = append(subjects_array, subject)
	}

	*s = subjects_array

	return nil
}

func UnmarshalSubject(src interface{}) (Subjects, error) {

	// Приведение полученных данных к корректному виду (массив байтов без служебных символов)
	byte_str, err := scan_prepare(src)
	if err != nil {
		return Subjects{}, err
	}

	// Объявляем карту с строчным ключём и значением в виде интерфейса
	var subject_maps []map[string]interface{}

	// Записываем значения в карту
	err = json.Unmarshal(byte_str, &subject_maps)
	if err != nil {
		return Subjects{}, err
	}

	s := Subjects{}

	for _, subject_map := range subject_maps {

		subject := Subject{}

		if subject_map["name"] != nil {
			if val, ok := subject_map["name"].(string); ok {
				subject.Name = val
			} else {
				return Subjects{}, errors.New("couldn't convert subject's name to string")
			}
		} else {
			return Subjects{}, errors.New("no subject's name found")
		}

		if subject_map["description"] != nil {
			if val, ok := subject_map["description"].(string); ok {
				if len(val) > 0 {
					subject.Description = val
				}
			} else {
				return Subjects{}, errors.New("couldn't convert subject's description to string")
			}
		} else {
			subject.Description = ""
		}

		s = append(s, subject)
	}

	return s, nil
}

func (subjects *Subjects) Contain(subject Subject) (res bool) {

	for _, s := range *subjects {
		if s.Name == subject.Name {
			return true
		}
	}

	return false
}

func (subjects *Subjects) Find(name string) (class Subject, err error) {

	for _, subject := range *subjects {
		if name == subject.Name {
			return subject, nil
		}
	}

	return Subject{}, errors.New("no subject with name " + name + " has found")
}
