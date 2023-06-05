package schedule_db_data

type Lesson struct {
	Number  int     `json:"number"`
	Name    string  `json:"name"`
	Room    string  `json:"room"`
	Teacher Teacher `json:"teacher"`
}

type Lessons []Lesson
