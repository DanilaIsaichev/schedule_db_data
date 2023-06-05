package schedule_db_data

type Schedule struct {
	Class   string  `json:"class"`
	Lessons Lessons `json:"lessons"`
}

type Schedules []Schedule
