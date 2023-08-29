package schedule_db_data

import (
	"database/sql"
	"errors"
	"fmt"
	"time"
)

func Get_teachers() (Teachers, error) {

	db, err := DB_connection(get_db_env())
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

		err := result.Scan(&teacher.Id, &teacher.Name, &teacher.Login)
		if err != nil {
			return Teachers{}, err
		}

		teachers = append(teachers, teacher)
	}

	return teachers, nil
}

func Get_subjects() (Subjects, error) {

	db, err := DB_connection(get_db_env())
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

		err := result.Scan(&subject.Id, &subject.Name, &subject.Description)
		if err != nil {
			return Subjects{}, err
		}

		subjects = append(subjects, subject)
	}

	return subjects, nil
}

func Get_rooms() (Rooms, error) {

	db, err := DB_connection(get_db_env())
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

		err := result.Scan(&room.Id, &room.Name, &room.Wing, &room.Floor)
		if err != nil {
			return Rooms{}, err
		}

		rooms = append(rooms, room)
	}

	return rooms, nil
}

func Get_classes() (Classes, error) {

	db, err := DB_connection(get_db_env())
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

		err := result.Scan(&class.Id, &class.Number, &class.Character)
		if err != nil {
			return Classes{}, err
		}

		classes = append(classes, class)
	}

	return classes, nil
}

func Get_changed_schedule(start time.Time, parallel int) (Days, error) {

	// ищем среди изменений

	start_str := start.Format("02.01.2006")

	// Подключаемся к БД
	db, err := DB_connection(get_db_env())
	if err != nil {
		return Days{}, err
	}
	defer db.Close()

	days := Days{}

	// Запрашиваем расписание по дате
	err = db.QueryRow("SELECT data FROM schedule WHERE is_base = False AND start = " + start_str + " AND parallel = " + fmt.Sprint(parallel) + ";").Scan(&days)
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

				// Получаем класс из строки
				class, err := new(Class).Parse(schedule.Class)
				if err != nil {
					return Days{}, err
				}

				// Ищем класс в данных из БД по номеру и букве
				found_class, err := classes.Find(class.Number, class.Character)
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
	db, err := DB_connection(get_db_env())
	if err != nil {
		return Days{}, err
	}
	defer db.Close()

	days := Days{}

	err = db.QueryRow("SELECT data FROM schedule WHERE is_base = True AND year = " + fmt.Sprint(year) + " AND parallel = " + fmt.Sprint(parallel) + ";").Scan(&days)
	switch {
	case err == sql.ErrNoRows:
		return Days{}, errors.New(fmt.Sprint("no schedules for year ", year))
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

				// Получаем класс из строки
				class, err := new(Class).Parse(schedule.Class)
				if err != nil {
					return Days{}, err
				}

				// Ищем класс в данных из БД по номеру и букве
				found_class, err := classes.Find(class.Number, class.Character)
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
