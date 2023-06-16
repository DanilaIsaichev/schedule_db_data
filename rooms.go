package schedule_db_data

import (
	"encoding/json"
	"errors"
	"strings"
)

type Room struct {
	Name  string `json:"name"`
	Wing  int   `json:"wing"`
	Floor int   `json:"floor"`
}

type Rooms []Room

// Сканер массива кабинетов
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
				room.Wing = int(val)
			} else {
				return errors.New("couldn't convert 'wing' to int")
			}
		} else {
			room.Wing = 0
		}

		if room_map["floor"] != nil {
			if val, ok := room_map["floor"].(float64); ok {
				room.Floor = int(val)
			} else {
				return errors.New("couldn't convert 'floor' to int")
			}
		} else {
			room.Floor = 0
		}

		rooms_array = append(rooms_array, room)
	}

	*r = rooms_array

	return nil
}

func UnmarshalRoom(src interface{}) (Rooms, error) {

	// Приведение полученных данных к корректному виду (массив байтов без служебных символов)
	byte_str, err := scan_prepare(src)
	if err != nil {
		return Rooms{}, err
	}

	// Объявляем карту с строчным ключём и значением в виде интерфейса
	var room_maps []map[string]interface{}

	// Записываем значения в карту
	err = json.Unmarshal(byte_str, &room_maps)
	if err != nil {
		return Rooms{}, err
	}

	r := Rooms{}

	for _, room_map := range room_maps {

		room := Room{}

		if room_map["name"] != nil {
			if val, ok := room_map["name"].(string); ok {
				room.Name = val
			} else {
				return Rooms{}, errors.New("couldn't convert rooms's name to string")
			}
		} else {
			return Rooms{}, errors.New("no rooms's name found")
		}

		if room_map["wing"] != nil {
			if val, ok := room_map["wing"].(float64); ok {
				room.Wing = int(val)
			} else {
				return Rooms{}, errors.New("couldn't convert room's wing to int")
			}
		} else {
			room.Wing = 0
		}

		if room_map["floor"] != nil {
			if val, ok := room_map["floor"].(float64); ok {
				room.Floor = int(val)
			} else {
				return Rooms{}, errors.New("couldn't convert room's floor to string")
			}
		} else {
			room.Floor = 0
		}

		r = append(r, room)
	}

	return r, nil
}
