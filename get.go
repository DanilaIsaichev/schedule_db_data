package schedule_db_data

import (
	"database/sql"
	"errors"
	"fmt"
	"time"
)

type Request struct {
	Is_base  bool       `json:"is_base"`
	Start    *time.Time `json:"start"`
	Year     int        `json:"year"`
	Parallel int        `json:"parallel"`
}

type Response struct {
	Classes  Classes  `json:"classes"`
	Teachers Teachers `json:"teachers"`
	Rooms    Rooms    `json:"rooms"`
	Subjects Subjects `json:"subjects"`
	Days     Days     `json:"days"`
}

type Generator_response struct {
	Classes  map[string]Class_lessons   `json:"classes"`
	Teachers map[string]Teacher_lessons `json:"teachers"`
}

type Class_lesson struct {
	Time    time.Time `json:"time"`
	Name    string    `json:"name"`
	Room    Room      `json:"room"`
	Teacher Teacher   `json:"teacher"`
}

type Class_lessons []Class_lesson

type Teacher_lesson struct {
	Time  time.Time `json:"time"`
	Name  string    `json:"name"`
	Room  Room      `json:"room"`
	Class Class     `json:"class"`
}

type Teacher_lessons []Teacher_lesson

func Get_teachers() (Teachers, error) {

	db, err := DB_connection(Get_db_env("getter"))
	if err != nil {
		return Teachers{}, err
	}
	defer db.Close()

	result, err := db.Query("SELECT * FROM teacher;")
	if err != nil {
		return Teachers{}, err
	}
	defer result.Close()

	teachers := Teachers{}

	for result.Next() {

		teacher := Teacher{}

		err := result.Scan(&teacher.Login, &teacher.Name)
		if err != nil {
			return Teachers{}, err
		}

		teachers = append(teachers, teacher)
	}

	return teachers, nil
}

func Get_subjects() (Subjects, error) {

	db, err := DB_connection(Get_db_env("getter"))
	if err != nil {
		return Subjects{}, err
	}
	defer db.Close()

	result, err := db.Query("SELECT * FROM subject;")
	if err != nil {
		return Subjects{}, nil
	}
	defer result.Close()

	subjects := Subjects{}

	for result.Next() {

		subject := Subject{}

		err := result.Scan(&subject.Name, &subject.Description)
		if err != nil {
			return Subjects{}, err
		}

		subjects = append(subjects, subject)
	}

	return subjects, nil
}

func Get_rooms() (Rooms, error) {

	db, err := DB_connection(Get_db_env("getter"))
	if err != nil {
		return Rooms{}, err
	}
	defer db.Close()

	result, err := db.Query("SELECT * FROM room;")
	if err != nil {
		return Rooms{}, err
	}
	defer result.Close()

	rooms := Rooms{}

	for result.Next() {

		room := Room{}

		err := result.Scan(&room.Name, &room.Wing, &room.Floor)
		if err != nil {
			return Rooms{}, err
		}

		rooms = append(rooms, room)
	}

	return rooms, nil
}

func Get_classes() (Classes, error) {

	db, err := DB_connection(Get_db_env("getter"))
	if err != nil {
		return Classes{}, err
	}
	defer db.Close()

	result, err := db.Query("SELECT * FROM class;")
	if err != nil {
		return Classes{}, err
	}
	defer result.Close()

	classes := Classes{}

	for result.Next() {

		class := Class{}

		err := result.Scan(&class.Number, &class.Character)
		if err != nil {
			return Classes{}, err
		}

		classes = append(classes, class)
	}

	return classes, nil
}

func Get_changed_schedule(start time.Time, parallel int) (Days, error) {

	// ищем среди изменений

	// Подключаемся к БД
	db, err := DB_connection(Get_db_env("getter"))
	if err != nil {
		return Days{}, err
	}
	defer db.Close()

	days := Days{}

	// Запрашиваем расписание по дате
	err = db.QueryRow("SELECT data FROM schedule WHERE is_base = False AND start = '" + start.Format("2006-01-02") + "' AND parallel = " + fmt.Sprint(parallel) + ";").Scan(&days)
	switch {
	case err == sql.ErrNoRows:

		year := start.Year()

		if start.Month() <= 8 {
			year -= 1
		}

		days, err := Get_base_schedule(year, parallel)
		if err != nil {
			return Days{}, err
		}

		return days, nil

	case err != nil:
		return Days{}, err
	default:

		// Проверка существования классов, учителей, кабинетов и предметов
		classes, err := Get_classes()
		if err != nil {
			return Days{}, err
		}

		teachers, err := Get_teachers()
		if err != nil {
			return Days{}, err
		}

		rooms, err := Get_rooms()
		if err != nil {
			return Days{}, err
		}

		subjects, err := Get_subjects()
		if err != nil {
			return Days{}, err
		}

		for _, day := range days {
			for schedule_id, schedule := range day.Schedule {

				// Ищем класс в данных из БД по номеру и букве
				found_class, err := classes.Find(schedule.Class)
				if err != nil {
					// Отбрасываем расписание для несуществующего класса
					day.Schedule = append(day.Schedule[:schedule_id], day.Schedule[schedule_id+1:]...)
				} else {

					schedule.Class = found_class.ToString()

					for lesson_id, lesson := range schedule.Lessons {

						// Проверяем наличие предмета, кабинета и учителя в данных из БД
						if subjects.Contain(Subject{Name: lesson.Name}) && rooms.Contain(Room{Name: lesson.Room}) && teachers.Contain(lesson.Teacher) {
							// Передаём учителю данные из БД
							lesson.Teacher, err = teachers.Find(lesson.Teacher.Login)
							if err != nil {
								// Отбрасываем урок, если в БД нет данных об учителе
								schedule.Lessons = append(schedule.Lessons[:lesson_id], schedule.Lessons[lesson_id+1:]...)
							}
						} else {
							// Отбрасываем урок, если в БД нет данных о предмете или кабинете, или учителе
							schedule.Lessons = append(schedule.Lessons[:lesson_id], schedule.Lessons[lesson_id+1:]...)
						}
					}
				}
			}
		}

		return days, nil
	}
}

func Get_base_schedule(year int, parallel int) (Days, error) {

	// Подключаемся к БД
	db, err := DB_connection(Get_db_env("getter"))
	if err != nil {
		return Days{}, err
	}
	defer db.Close()

	days := Days{}

	err = db.QueryRow("SELECT data FROM schedule WHERE is_base = True AND year = " + fmt.Sprint(year) + " AND parallel = " + fmt.Sprint(parallel) + ";").Scan(&days)
	if err != nil {
		return Days{}, err
	} else {
		// Проверка существования классов, учителей, кабинетов и предметов
		classes, err := Get_classes()
		if err != nil {
			return Days{}, err
		}

		teachers, err := Get_teachers()
		if err != nil {
			return Days{}, err
		}

		rooms, err := Get_rooms()
		if err != nil {
			return Days{}, err
		}

		subjects, err := Get_subjects()
		if err != nil {
			return Days{}, err
		}

		for _, day := range days {
			for schedule_id, schedule := range day.Schedule {

				// Ищем класс в данных из БД по номеру и букве
				found_class, err := classes.Find(schedule.Class)
				if err != nil {
					// Отбрасываем расписание для несуществующего класса
					day.Schedule = append(day.Schedule[:schedule_id], day.Schedule[schedule_id+1:]...)
				} else {

					schedule.Class = found_class.ToString()

					for lesson_id, lesson := range schedule.Lessons {

						// Проверяем наличие предмета, кабинета и учителя в данных из БД
						if subjects.Contain(Subject{Name: lesson.Name}) && rooms.Contain(Room{Name: lesson.Room}) && teachers.Contain(lesson.Teacher) {
							// Передаём учителю данные из БД
							lesson.Teacher, err = teachers.Find(lesson.Teacher.Login)
							if err != nil {
								// Отбрасываем урок, если в БД нет данных об учителе
								schedule.Lessons = append(schedule.Lessons[:lesson_id], schedule.Lessons[lesson_id+1:]...)
							}
						} else {
							// Отбрасываем урок, если в БД нет данных о предмете или кабинете, или учителе
							schedule.Lessons = append(schedule.Lessons[:lesson_id], schedule.Lessons[lesson_id+1:]...)
						}
					}
				}
			}
		}

		return days, nil

	}
}

// Универсальная функция для получения структуры ответа по структуре запроса
func Get_editor_data(req Request) (Response, error) {

	classes, err := Get_classes()
	if err != nil {
		return Response{}, err
	}

	teachers, err := Get_teachers()
	if err != nil {
		return Response{}, err
	}

	rooms, err := Get_rooms()
	if err != nil {
		return Response{}, err
	}

	subjects, err := Get_subjects()
	if err != nil {
		return Response{}, err
	}

	// Подключаемся к БД
	db, err := DB_connection(Get_db_env("getter"))
	if err != nil {
		return Response{}, err
	}
	defer db.Close()

	days := Days{}

	if req.Is_base {
		err = db.QueryRow("SELECT data FROM schedule WHERE is_base = True AND year = " + fmt.Sprint(req.Year) + " AND parallel = " + fmt.Sprint(req.Parallel) + ";").Scan(&days)
		if err != nil {
			return Response{}, err
		}
	} else {
		err = db.QueryRow("SELECT data FROM schedule WHERE is_base = False AND start = " + req.Start.Format("2006-01-02") + " AND parallel = " + fmt.Sprint(req.Parallel) + ";").Scan(&days)
		if err == sql.ErrNoRows {
			err = db.QueryRow("SELECT data FROM schedule WHERE is_base = True AND year = " + fmt.Sprint(req.Year) + " AND parallel = " + fmt.Sprint(req.Parallel) + ";").Scan(&days)
			if err == sql.ErrNoRows {
				return Response{Teachers: teachers, Classes: classes, Rooms: rooms, Subjects: subjects, Days: Days{}}, err
			} else if err != nil {
				return Response{}, err
			}
		} else if err != nil {
			return Response{}, err
		}
	}

	// Проверка существования классов, учителей, кабинетов и предметов
	for _, day := range days {
		for schedule_id, schedule := range day.Schedule {

			// Ищем класс в данных из БД по номеру и букве
			found_class, err := classes.Find(schedule.Class)
			if err != nil {
				// Отбрасываем расписание для несуществующего класса
				day.Schedule = append(day.Schedule[:schedule_id], day.Schedule[schedule_id+1:]...)
			} else {

				schedule.Class = found_class.ToString()

				for lesson_id, lesson := range schedule.Lessons {

					// Проверяем наличие предмета, кабинета и учителя в данных из БД
					if subjects.Contain(Subject{Name: lesson.Name}) && rooms.Contain(Room{Name: lesson.Room}) && teachers.Contain(lesson.Teacher) {
						// Передаём учителю данные из БД
						lesson.Teacher, err = teachers.Find(lesson.Teacher.Login)
						if err != nil {
							// Отбрасываем урок, если в БД нет данных об учителе
							schedule.Lessons = append(schedule.Lessons[:lesson_id], schedule.Lessons[lesson_id+1:]...)
						}
					} else {
						// Отбрасываем урок, если в БД нет данных о предмете или кабинете, или учителе
						schedule.Lessons = append(schedule.Lessons[:lesson_id], schedule.Lessons[lesson_id+1:]...)
					}
				}
			}
		}
	}

	return Response{Teachers: teachers, Classes: classes, Rooms: rooms, Subjects: subjects, Days: days}, nil
}

// Универсальная функция для получения структуры ответа по структуре запроса
func Get_generator_data(req Request) (Generator_response, error) {

	classes, err := Get_classes()
	if err != nil {
		return Generator_response{}, err
	}

	classes_schedules := make(map[string]Class_lessons)

	for _, class := range classes {
		classes_schedules[class.ToString()] = Class_lessons{}
	}

	teachers, err := Get_teachers()
	if err != nil {
		return Generator_response{}, err
	}

	teachers_schedules := make(map[string]Teacher_lessons)

	for _, teacher := range teachers {
		teachers_schedules[teacher.Login] = Teacher_lessons{}
	}

	rooms, err := Get_rooms()
	if err != nil {
		return Generator_response{}, err
	}

	subjects, err := Get_subjects()
	if err != nil {
		return Generator_response{}, err
	}

	// Подключаемся к БД
	db, err := DB_connection(Get_db_env("getter"))
	if err != nil {
		return Generator_response{}, err
	}
	defer db.Close()

	days := Days{}

	if req.Is_base {
		err = db.QueryRow("SELECT data FROM schedule WHERE is_base = True AND year = " + fmt.Sprint(req.Year) + " AND parallel = " + fmt.Sprint(req.Parallel) + ";").Scan(&days)
		if err != nil {
			return Generator_response{}, err
		}
	} else {
		err = db.QueryRow("SELECT data FROM schedule WHERE is_base = False AND start = " + req.Start.Format("2006-01-02") + " AND parallel = " + fmt.Sprint(req.Parallel) + ";").Scan(&days)
		if err == sql.ErrNoRows {
			err = db.QueryRow("SELECT data FROM schedule WHERE is_base = True AND year = " + fmt.Sprint(req.Year) + " AND parallel = " + fmt.Sprint(req.Parallel) + ";").Scan(&days)
			if err != nil {
				return Generator_response{}, err
			}
		} else if err != nil {
			return Generator_response{}, err
		}
	}

	// Проверка существования классов, учителей, кабинетов и предметов
	for _, day := range days {
		for schedule_id, schedule := range day.Schedule {

			// Ищем класс в данных из БД по номеру и букве
			found_class, err := classes.Find(schedule.Class)
			if err != nil {
				// Отбрасываем расписание для несуществующего класса
				day.Schedule = append(day.Schedule[:schedule_id], day.Schedule[schedule_id+1:]...)
			} else {

				schedule.Class = found_class.ToString()

				for lesson_id, lesson := range schedule.Lessons {

					// Проверяем наличие предмета, кабинета и учителя в данных из БД
					if subjects.Contain(Subject{Name: lesson.Name}) && rooms.Contain(Room{Name: lesson.Room}) && teachers.Contain(lesson.Teacher) {
						// Передаём учителю данные из БД

						time_str := day.Date

						if lesson.Number >= 1 && lesson.Number <= 8 {
							switch {
							case lesson.Number == 1:
								time_str += " 8:30"
							case lesson.Number == 2:
								time_str += " 9:30"
							case lesson.Number == 3:
								time_str += " 10:30"
							case lesson.Number == 4:
								time_str += " 11:30"
							case lesson.Number == 5:
								time_str += " 12:30"
							case lesson.Number == 6:
								time_str += " 13:30"
							case lesson.Number == 7:
								time_str += " 14:30"
							case lesson.Number == 8:
								time_str += " 15:30"
							}
						} else {
							return Generator_response{}, errors.New("wrong lesson number")
						}

						l_time, err := time.Parse("02.01.2006 15:04", time_str)
						if err != nil {
							return Generator_response{}, err
						}

						l_room, err := rooms.Find(lesson.Room)
						if err != nil {
							return Generator_response{}, err
						}

						l_teacher := lesson.Teacher

						l_class, err := classes.Find(schedule.Class)
						if err != nil {
							return Generator_response{}, err
						}

						c_lesson := Class_lesson{Time: l_time, Name: lesson.Name, Room: l_room, Teacher: l_teacher}
						classes_schedules[schedule.Class] = append(classes_schedules[schedule.Class], c_lesson)

						t_lesson := Teacher_lesson{Time: l_time, Name: lesson.Name, Room: l_room, Class: l_class}
						teachers_schedules[l_teacher.Login] = append(teachers_schedules[l_teacher.Login], t_lesson)

					} else {
						// Отбрасываем урок, если в БД нет данных о предмете или кабинете, или учителе
						schedule.Lessons = append(schedule.Lessons[:lesson_id], schedule.Lessons[lesson_id+1:]...)
					}
				}
			}
		}
	}

	return Generator_response{Classes: classes_schedules, Teachers: teachers_schedules}, nil
}
