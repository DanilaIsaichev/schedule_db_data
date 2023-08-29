package schedule_db_data

import (
	"encoding/json"
	"errors"
	"strings"
)

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

func (rooms *Rooms) Contain(room Room) (res bool) {

	for _, r := range *rooms {
		if r.Name == room.Name {
			return true
		}
	}

	return false
}

func (rooms *Rooms) Find(name string) (room Room, err error) {

	for _, room := range *rooms {
		if name == room.Name {
			return room, nil
		}
	}

	return Room{}, errors.New("no room with name " + name + " has found")
}
