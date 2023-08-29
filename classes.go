package schedule_db_data

import (
	"encoding/json"
	"errors"
	"fmt"
	"strconv"
	"strings"
)

type Class struct {
	Id        int    `json:"id"`
	Number    int    `json:"number"`
	Character string `json:"сharacter"`
}

func (c *Class) ToString() (class string) {
	return fmt.Sprint(c.Number, c.Character)
}

func (c *Class) Parse(class_string string) (class Class, err error) {

	c.Number, err = strconv.Atoi(class_string[:len(class_string)-1])
	if err != nil {
		return Class{}, err
	}

	c.Character = class_string[len(class_string)-1:]

	return *c, nil
}

type Classes []Class

// Сканер массива классов
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

func (classes *Classes) Contain(class Class) (res bool) {

	for _, c := range *classes {
		if c.Character == class.Character && c.Number == class.Number {
			return true
		}
	}

	return false
}

func (classes *Classes) Find(name string) (class Class, err error) {

	c, err := new(Class).Parse(name)
	if err != nil {
		return Class{}, err
	}

	for _, class := range *classes {
		if c.Character == class.Character && c.Number == class.Number {
			return class, nil
		}
	}

	return Class{}, errors.New("no class with name " + name + " has found")
}
