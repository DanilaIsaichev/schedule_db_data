package schedule_db_data

type Lesson struct {
	Number      int         `json:"number"`
	Subject     Subject     `json:"subject"`
	Lesson_data Lesson_data `json:"lesson_data"`
}

type Lessons []Lesson

type Lesson_data []struct {
	Room    Room    `json:"room"`
	Teacher Teacher `json:"teacher"`
}
