package schedule_db_data

import (
	"encoding/json"
	"errors"
)

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
