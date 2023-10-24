package schedule_db_data

type Schedule struct {
	Group   string  `json:"group"`
	Lessons Lessons `json:"lessons"`
}

type Schedules []Schedule
