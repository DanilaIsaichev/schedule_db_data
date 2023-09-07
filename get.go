package schedule_db_data

import (
	"database/sql"
	"errors"
	"fmt"
	"time"
)

type Request struct {
	Is_base  bool      `json:"is_base"`
	Start    time.Time `json:"start"`
	Year     int       `json:"year"`
	Parallel int       `json:"parallel"`
}

type Response struct {
	Classes  Classes  `json:"classes"`
	Teachers Teachers `json:"teachers"`
	Rooms    Rooms    `json:"rooms"`
	Subjects Subjects `json:"subjects"`
	Days     Days     `json:"days"`
}

/*
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
}*/

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
		if err == sql.ErrNoRows {
			return Response{Teachers: teachers, Classes: classes, Rooms: rooms, Subjects: subjects, Days: Days{}}, err
		} else if err != nil {
			return Response{}, err
		}
	} else {
		err = db.QueryRow("SELECT data FROM schedule WHERE is_base = False AND start = '" + req.Start.Format("2006-01-02") + "' AND parallel = " + fmt.Sprint(req.Parallel) + ";").Scan(&days)
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
	for day_id, day := range days {
		for schedule_id, schedule := range day.Schedule {

			// Ищем класс в данных из БД по номеру и букве
			found_class, err := classes.Find(schedule.Class)
			if err != nil {
				// Отбрасываем расписание для несуществующего класса
				day.Schedule = append(day.Schedule[:schedule_id], day.Schedule[schedule_id+1:]...)
			} else {

				schedule.Class = found_class.ToString()
				for lesson_id, lesson := range schedule.Lessons {

					if len(lesson.Lesson_data) == lesson.Subject.Groups {

						for l_data_id := range lesson.Lesson_data {

							// Проверяем наличие предмета, кабинета и учителя в данных из БД
							if subjects.Contain(lesson.Subject) && rooms.Contain(lesson.Lesson_data[l_data_id].Room) && teachers.Contain(lesson.Lesson_data[l_data_id].Teacher) {

								// Передаём учителю данные из БД
								days[day_id].Schedule[schedule_id].Lessons[lesson_id].Lesson_data[l_data_id].Teacher, err = teachers.Find_by_login(lesson.Lesson_data[l_data_id].Teacher.Login)
								if err != nil {
									// Отбрасываем урок, если в БД нет данных об учителе
									days[day_id].Schedule[schedule_id].Lessons = append(schedule.Lessons[:lesson_id], schedule.Lessons[lesson_id+1:]...)
								}
							} else {
								// Отбрасываем урок, если в БД нет данных о предмете или кабинете, или учителе
								days[day_id].Schedule[schedule_id].Lessons = append(schedule.Lessons[:lesson_id], schedule.Lessons[lesson_id+1:]...)
							}
						}
					} else {
						return Response{}, errors.New("wrong lesson data")
					}
				}
			}
		}
	}

	return Response{Teachers: teachers, Classes: classes, Rooms: rooms, Subjects: subjects, Days: days}, nil
}
