package schedule_db_data

import (
	"encoding/json"
)

type Week struct {
	Start    string `json:"start"`
	Year     int    `json:"year"`
	Parallel int    `json:"parallel"`
	Is_Base  bool   `json:"is_base"`
	Data     Days   `json:"data"`
}

func Get_weeks(buff []byte) (Week, error) {

	data := Week{}
	err := json.Unmarshal(buff, &data)
	if err != nil {
		return Week{}, err
	}

	return data, nil
}
