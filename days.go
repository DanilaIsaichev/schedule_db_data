package schedule_db_data

import "encoding/json"

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

func Get_days(buff []byte) (Days, error) {

	data := Days{}
	err := json.Unmarshal(buff, &data)
	if err != nil {
		return Days{}, err
	}

	return data, nil
}
